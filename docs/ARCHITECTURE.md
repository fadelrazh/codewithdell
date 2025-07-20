# CodeWithDell Architecture

## Overview

CodeWithDell is a modern blog and showcase platform built with a microservices architecture using Next.js for the frontend and Golang for the backend.

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (Next.js)     │◄──►│   (Golang)      │◄──►│  (PostgreSQL)   │
│                 │    │                 │    │                 │
│ - React 18      │    │ - Gin Framework │    │ - ACID Compliant│
│ - TypeScript    │    │ - GORM ORM      │    │ - Full-text     │
│ - Tailwind CSS  │    │ - JWT Auth      │    │   Search        │
│ - Framer Motion │    │ - Redis Cache   │    │ - Trigram Index │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx Proxy   │    │   Redis Cache   │    │   File Storage  │
│                 │    │                 │    │                 │
│ - Load Balancer │    │ - Session Store │    │ - Local/S3      │
│ - Rate Limiting │    │ - Cache Layer   │    │ - Image Upload  │
│ - SSL Termination│   │ - Pub/Sub       │    │ - Static Files  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Backend Architecture (Clean Architecture)

### Layers

1. **Presentation Layer** (Handlers)

   - HTTP request/response handling
   - Input validation
   - Authentication middleware
   - Rate limiting

2. **Business Logic Layer** (Services)

   - Core business rules
   - Data transformation
   - External service integration
   - Caching logic

3. **Data Access Layer** (Repositories)

   - Database operations
   - Query optimization
   - Data mapping
   - Transaction management

4. **Infrastructure Layer**
   - Database connections
   - External APIs
   - File storage
   - Email services

### Directory Structure

```
backend/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── database/         # Database connection & migrations
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repositories/     # Data access layer
│   ├── routes/           # Route definitions
│   ├── services/         # Business logic
│   ├── utils/            # Utility functions
│   └── logger/           # Logging system
├── pkg/                  # Public packages
└── main.go              # Application entry point
```

## Frontend Architecture

### Technology Stack

- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Data Fetching**: React Query (TanStack Query)
- **Forms**: React Hook Form + Zod
- **Authentication**: NextAuth.js
- **Animations**: Framer Motion

### Directory Structure

```
frontend/
├── src/
│   ├── app/              # App Router pages
│   ├── components/       # Reusable components
│   │   ├── ui/          # Base UI components
│   │   ├── forms/       # Form components
│   │   └── layout/      # Layout components
│   ├── hooks/           # Custom React hooks
│   ├── lib/             # Utility libraries
│   ├── stores/          # Zustand stores
│   ├── types/           # TypeScript types
│   └── utils/           # Utility functions
├── public/              # Static assets
└── prisma/              # Database schema (if using Prisma)
```

## Database Design

### Core Tables

1. **users** - User accounts and profiles
2. **posts** - Blog posts
3. **projects** - Showcase projects
4. **categories** - Content categories
5. **tags** - Content tags
6. **technologies** - Technology stack
7. **comments** - User comments
8. **likes** - User likes
9. **bookmarks** - User bookmarks
10. **screenshots** - Project screenshots

### Relationships

- Users can have many Posts and Projects
- Posts and Projects can have many Categories and Tags
- Projects can have many Technologies
- Posts and Projects can have many Comments, Likes, and Bookmarks
- Comments can have parent-child relationships

### Indexes

- Primary keys on all tables
- Unique indexes on slugs, emails, usernames
- Full-text search indexes on content
- Trigram indexes for fuzzy search
- Composite indexes for common queries

## Security Architecture

### Authentication

- JWT-based authentication
- Refresh token rotation
- Password hashing with bcrypt
- Rate limiting on auth endpoints
- Session management with Redis

### Authorization

- Role-based access control (RBAC)
- Resource-level permissions
- Admin, Editor, and User roles
- Middleware-based authorization

### Data Protection

- Input validation and sanitization
- SQL injection prevention
- XSS protection
- CSRF protection
- Content Security Policy (CSP)

## Performance Architecture

### Caching Strategy

- Redis for session storage
- Redis for API response caching
- Browser caching for static assets
- CDN for global content delivery

### Database Optimization

- Connection pooling
- Query optimization
- Index strategy
- Read replicas (for production)

### Frontend Optimization

- Code splitting
- Image optimization
- Lazy loading
- Service worker for offline support

## Monitoring & Observability

### Logging

- Structured logging with zerolog
- Request/response logging
- Error tracking
- Performance metrics

### Metrics

- Prometheus metrics
- Custom business metrics
- Health checks
- Performance monitoring

### Tracing

- Request tracing
- Database query tracing
- External service tracing

## Deployment Architecture

### Development

- Docker Compose for local development
- Hot reloading for both frontend and backend
- Local database and Redis instances
- Development tools (Adminer, Redis Commander)

### Production

- Containerized deployment
- Load balancer (Nginx)
- Database clustering
- Redis clustering
- CDN integration
- SSL/TLS termination

## API Design

### RESTful Endpoints

- `/api/v1/auth/*` - Authentication
- `/api/v1/posts/*` - Blog posts
- `/api/v1/projects/*` - Showcase projects
- `/api/v1/users/*` - User management
- `/api/v1/admin/*` - Admin operations

### Response Format

```json
{
  "success": true,
  "data": {},
  "message": "Success",
  "meta": {
    "pagination": {},
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### Error Handling

- Standardized error responses
- HTTP status codes
- Error codes for client handling
- Detailed logging for debugging

## Scalability Considerations

### Horizontal Scaling

- Stateless backend services
- Database read replicas
- Redis clustering
- Load balancer distribution

### Vertical Scaling

- Resource optimization
- Database query optimization
- Caching strategies
- CDN usage

### Future Considerations

- Microservices decomposition
- Event-driven architecture
- GraphQL API
- Real-time features (WebSocket)
- Mobile app support
