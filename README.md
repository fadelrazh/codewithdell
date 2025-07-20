# CodeWithDell - Modern Blog Platform

A full-stack blog platform built with **Next.js 15** (Frontend) and **Go** (Backend), featuring modern authentication, real-time content management, and a beautiful responsive UI.

## 🚀 Features

### Backend (Go)

- **RESTful API** with Gin framework
- **JWT Authentication** with refresh tokens
- **Role-based Authorization** (User/Admin)
- **PostgreSQL Database** with GORM ORM
- **CRUD Operations** for posts, users, and profiles
- **Database Migrations** and seeding
- **Middleware** for authentication and logging
- **Health Checks** and monitoring
- **Docker** support

### Frontend (Next.js 15)

- **Modern UI** with Tailwind CSS
- **Responsive Design** for all devices
- **Dark Mode** support
- **TypeScript** for type safety
- **Client-side Authentication** with localStorage
- **Real-time Updates** with React hooks
- **Form Validation** and error handling
- **SEO Optimized** with metadata

## 🛠️ Tech Stack

### Backend

- **Go 1.21+** - Programming language
- **Gin** - Web framework
- **GORM** - ORM for database
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Docker** - Containerization

### Frontend

- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **React Hook Form** - Form management
- **Zustand** - State management
- **Framer Motion** - Animations

## 📦 Installation

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Docker (optional)

### Quick Start

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/codewithdell.git
cd codewithdell
```

2. **Set up environment variables**

```bash
cp env.example .env
# Edit .env with your database credentials
```

3. **Start PostgreSQL with Docker**

```bash
docker-compose up -d postgres
```

4. **Run Backend**

```bash
cd backend
go mod download
go run main.go
```

5. **Run Frontend**

```bash
cd frontend
npm install
npm run dev
```

6. **Access the application**

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger

## 🗄️ Database Setup

The application uses PostgreSQL with the following main tables:

- **users** - User accounts and profiles
- **posts** - Blog posts and articles
- **tags** - Post categorization
- **post_tags** - Many-to-many relationship

### Sample Data

The application includes sample data for testing:

- Admin user: `admin@codewithdell.com` / `password`
- Sample posts with different statuses

## 🔐 Authentication

### User Registration

```bash
POST /api/v1/auth/register
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "password123",
  "username": "johndoe"
}
```

### User Login

```bash
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}
```

### Protected Routes

- `/api/v1/profile` - User profile (requires auth)
- `/api/v1/admin/*` - Admin operations (requires admin role)

## 📝 API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /api/v1/test` - Test endpoint
- `GET /api/v1/posts` - Get all posts
- `GET /api/v1/posts/:slug` - Get post by slug
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token

### Protected Endpoints

- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

### Admin Endpoints

- `POST /api/v1/admin/posts` - Create post
- `PUT /api/v1/admin/posts/:id` - Update post
- `DELETE /api/v1/admin/posts/:id` - Delete post

## 🎨 Frontend Pages

### Public Pages

- **Home** (`/`) - Landing page with features and CTA
- **Posts** (`/posts`) - Blog post listing
- **Post Detail** (`/posts/[slug]`) - Individual post view
- **Login** (`/login`) - User authentication
- **Register** (`/register`) - User registration

### Protected Pages

- **Profile** (`/profile`) - User profile management

## 🚀 Deployment

### Backend Deployment

```bash
# Build Docker image
docker build -t codewithdell-backend ./backend

# Run container
docker run -p 8080:8080 --env-file .env codewithdell-backend
```

### Frontend Deployment

```bash
# Build for production
cd frontend
npm run build

# Start production server
npm start
```

### Environment Variables

#### Backend (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=codewithdell
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret
```

#### Frontend (.env.local)

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm test
```

## 📁 Project Structure

```
codewithdell/
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── routes/
│   │   └── server/
│   ├── main.go
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── package.json
│   └── next.config.js
├── docker-compose.yml
└── README.md
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Go web framework
- [Next.js](https://nextjs.org/) - React framework
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [GORM](https://gorm.io/) - Go ORM library

## 📞 Support

For support, email support@codewithdell.com or create an issue in this repository.

---

**Built with ❤️ by the CodeWithDell Team**
