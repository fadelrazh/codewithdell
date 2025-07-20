#!/bin/bash

# Development script for CodeWithDell
# This script helps set up and run the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    local missing_deps=()
    
    if ! command_exists docker; then
        missing_deps+=("docker")
    fi
    
    if ! command_exists docker-compose; then
        missing_deps+=("docker-compose")
    fi
    
    if ! command_exists node; then
        missing_deps+=("node")
    fi
    
    if ! command_exists npm; then
        missing_deps+=("npm")
    fi
    
    if ! command_exists go; then
        missing_deps+=("go")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        print_status "Please install the missing dependencies and try again."
        exit 1
    fi
    
    print_success "All prerequisites are installed!"
}

# Function to setup environment
setup_environment() {
    print_status "Setting up environment..."
    
    # Create .env files if they don't exist
    if [ ! -f ".env" ]; then
        cat > .env << EOF
# Application
ENVIRONMENT=development
PORT=8080
CORS_ORIGIN=http://localhost:3000
LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5433
DB_NAME=codewithdell
DB_USER=codewithdell
DB_PASSWORD=codewithdell123
DB_SSLMODE=disable
DB_MAX_OPEN=25
DB_MAX_IDLE=5
DB_TIMEOUT=5s

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# Email
EMAIL_PROVIDER=sendgrid
EMAIL_API_KEY=
EMAIL_FROM=noreply@codewithdell.com

# Storage
STORAGE_PROVIDER=local
STORAGE_BUCKET=codewithdell
STORAGE_REGION=us-east-1
STORAGE_ACCESS_KEY=
STORAGE_SECRET_KEY=
EOF
        print_success "Created .env file"
    fi
    
    if [ ! -f "frontend/.env.local" ]; then
        cat > frontend/.env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret-change-in-production
EOF
        print_success "Created frontend/.env.local file"
    fi
}

# Function to install dependencies
install_dependencies() {
    print_status "Installing dependencies..."
    
    # Backend dependencies
    if [ -f "backend/go.mod" ]; then
        print_status "Installing Go dependencies..."
        cd backend
        go mod download
        cd ..
        print_success "Go dependencies installed"
    fi
    
    # Frontend dependencies
    if [ -f "frontend/package.json" ]; then
        print_status "Installing Node.js dependencies..."
        cd frontend
        npm install
        cd ..
        print_success "Node.js dependencies installed"
    fi
}

# Function to start services
start_services() {
    print_status "Starting services with Docker Compose..."
    
    # Start database and Redis
    docker-compose up -d postgres redis
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check if services are running
    if ! docker-compose ps | grep -q "postgres.*Up"; then
        print_error "PostgreSQL failed to start"
        exit 1
    fi
    
    if ! docker-compose ps | grep -q "redis.*Up"; then
        print_error "Redis failed to start"
        exit 1
    fi
    
    print_success "Database and Redis are running!"
}

# Function to run migrations
run_migrations() {
    print_status "Running database migrations..."
    
    cd backend
    go run main.go migrate
    cd ..
    
    print_success "Database migrations completed!"
}

# Function to start development servers
start_dev_servers() {
    print_status "Starting development servers..."
    
    # Start backend in background
    print_status "Starting backend server..."
    cd backend
    go run main.go &
    BACKEND_PID=$!
    cd ..
    
    # Start frontend in background
    print_status "Starting frontend server..."
    cd frontend
    npm run dev &
    FRONTEND_PID=$!
    cd ..
    
    print_success "Development servers started!"
    print_status "Backend PID: $BACKEND_PID"
    print_status "Frontend PID: $FRONTEND_PID"
    print_status "Backend: http://localhost:8080"
    print_status "Frontend: http://localhost:3000"
    print_status "API Docs: http://localhost:8080/swagger/"
    print_status "pgAdmin4 (Database GUI): http://localhost:8081"
    print_status "Adminer (Alternative): http://localhost:8082"
    print_status "Redis GUI: http://localhost:8083"
    
    # Wait for user to stop
    echo ""
    print_warning "Press Ctrl+C to stop all servers"
    
    # Function to cleanup on exit
    cleanup() {
        print_status "Stopping servers..."
        kill $BACKEND_PID 2>/dev/null || true
        kill $FRONTEND_PID 2>/dev/null || true
        print_success "Servers stopped!"
        exit 0
    }
    
    # Set trap to cleanup on exit
    trap cleanup SIGINT SIGTERM
    
    # Wait for background processes
    wait
}

# Function to stop all services
stop_services() {
    print_status "Stopping all services..."
    
    # Stop Docker services
    docker-compose down
    
    # Kill any running processes
    pkill -f "go run main.go" 2>/dev/null || true
    pkill -f "npm run dev" 2>/dev/null || true
    
    print_success "All services stopped!"
}

# Function to show status
show_status() {
    print_status "Checking service status..."
    
    echo ""
    echo "Docker Services:"
    docker-compose ps
    
    echo ""
    echo "Running Processes:"
    ps aux | grep -E "(go run main.go|npm run dev)" | grep -v grep || echo "No development servers running"
}

# Function to show logs
show_logs() {
    local service=${1:-"all"}
    
    case $service in
        "backend"|"frontend"|"postgres"|"redis")
            print_status "Showing logs for $service..."
            docker-compose logs -f $service
            ;;
        "all")
            print_status "Showing all logs..."
            docker-compose logs -f
            ;;
        *)
            print_error "Unknown service: $service"
            print_status "Available services: backend, frontend, postgres, redis, all"
            exit 1
            ;;
    esac
}

# Function to reset database
reset_database() {
    print_warning "This will delete all data in the database!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Resetting database..."
        docker-compose down -v
        docker-compose up -d postgres redis
        sleep 10
        run_migrations
        print_success "Database reset completed!"
    else
        print_status "Database reset cancelled"
    fi
}

# Function to show help
show_help() {
    echo "CodeWithDell Development Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  setup     - Set up the development environment"
    echo "  start     - Start all development services"
    echo "  stop      - Stop all services"
    echo "  status    - Show service status"
    echo "  logs      - Show logs (backend|frontend|postgres|redis|all)"
    echo "  reset-db  - Reset the database"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 setup"
    echo "  $0 start"
    echo "  $0 logs backend"
    echo "  $0 reset-db"
}

# Main script logic
case "${1:-help}" in
    "setup")
        check_prerequisites
        setup_environment
        install_dependencies
        start_services
        run_migrations
        print_success "Setup completed! Run '$0 start' to start development servers."
        ;;
    "start")
        start_services
        start_dev_servers
        ;;
    "stop")
        stop_services
        ;;
    "status")
        show_status
        ;;
    "logs")
        show_logs "$2"
        ;;
    "reset-db")
        reset_database
        ;;
    "help"|*)
        show_help
        ;;
esac 