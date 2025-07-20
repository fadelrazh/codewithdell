# CodeWithDell API Documentation

## Base URL

```
Development: http://localhost:8080
Production: https://api.codewithdell.com
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "data": {},
  "message": "Success",
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "pages": 10
    },
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

## Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "field": "email",
      "message": "Email is required"
    }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Validation Error
- `429` - Too Many Requests
- `500` - Internal Server Error

## Authentication Endpoints

### Register User

```http
POST /api/v1/auth/register
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "uuid": "uuid-string",
      "email": "user@example.com",
      "username": "username",
      "firstName": "John",
      "lastName": "Doe",
      "role": "user",
      "status": "active",
      "verified": false,
      "createdAt": "2024-01-01T00:00:00Z"
    },
    "tokens": {
      "accessToken": "jwt-token",
      "refreshToken": "refresh-token"
    }
  },
  "message": "User registered successfully"
}
```

### Login User

```http
POST /api/v1/auth/login
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "uuid": "uuid-string",
      "email": "user@example.com",
      "username": "username",
      "firstName": "John",
      "lastName": "Doe",
      "role": "user",
      "status": "active",
      "verified": true,
      "lastLogin": "2024-01-01T00:00:00Z"
    },
    "tokens": {
      "accessToken": "jwt-token",
      "refreshToken": "refresh-token"
    }
  },
  "message": "Login successful"
}
```

### Refresh Token

```http
POST /api/v1/auth/refresh
```

**Request Body:**

```json
{
  "refreshToken": "refresh-token"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "accessToken": "new-jwt-token",
    "refreshToken": "new-refresh-token"
  },
  "message": "Token refreshed successfully"
}
```

## Blog Posts Endpoints

### Get All Posts

```http
GET /api/v1/posts
```

**Query Parameters:**

- `page` (number) - Page number (default: 1)
- `limit` (number) - Items per page (default: 10)
- `category` (string) - Filter by category slug
- `tag` (string) - Filter by tag slug
- `search` (string) - Search in title and content
- `status` (string) - Filter by status (published, draft)

**Response:**

```json
{
  "success": true,
  "data": {
    "posts": [
      {
        "id": 1,
        "uuid": "uuid-string",
        "title": "Getting Started with Go",
        "slug": "getting-started-with-go",
        "excerpt": "Learn the basics of Go programming...",
        "content": "Full post content...",
        "featuredImage": "https://example.com/image.jpg",
        "status": "published",
        "publishedAt": "2024-01-01T00:00:00Z",
        "viewCount": 150,
        "likeCount": 25,
        "commentCount": 10,
        "author": {
          "id": 1,
          "username": "author",
          "firstName": "John",
          "lastName": "Doe",
          "avatar": "https://example.com/avatar.jpg"
        },
        "categories": [
          {
            "id": 1,
            "name": "Programming",
            "slug": "programming",
            "color": "#3B82F6"
          }
        ],
        "tags": [
          {
            "id": 1,
            "name": "Go",
            "slug": "go",
            "color": "#00ADD8"
          }
        ],
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ]
  },
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "pages": 10
    }
  }
}
```

### Get Post by Slug

```http
GET /api/v1/posts/{slug}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "post": {
      "id": 1,
      "uuid": "uuid-string",
      "title": "Getting Started with Go",
      "slug": "getting-started-with-go",
      "excerpt": "Learn the basics of Go programming...",
      "content": "Full post content...",
      "featuredImage": "https://example.com/image.jpg",
      "status": "published",
      "publishedAt": "2024-01-01T00:00:00Z",
      "viewCount": 150,
      "likeCount": 25,
      "commentCount": 10,
      "author": {
        "id": 1,
        "username": "author",
        "firstName": "John",
        "lastName": "Doe",
        "avatar": "https://example.com/avatar.jpg"
      },
      "categories": [...],
      "tags": [...],
      "comments": [...],
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

## Projects Endpoints

### Get All Projects

```http
GET /api/v1/projects
```

**Query Parameters:**

- `page` (number) - Page number (default: 1)
- `limit` (number) - Items per page (default: 10)
- `category` (string) - Filter by category slug
- `technology` (string) - Filter by technology slug
- `difficulty` (string) - Filter by difficulty (beginner, intermediate, advanced, expert)
- `search` (string) - Search in title and description

**Response:**

```json
{
  "success": true,
  "data": {
    "projects": [
      {
        "id": 1,
        "uuid": "uuid-string",
        "title": "E-commerce Platform",
        "slug": "ecommerce-platform",
        "description": "A full-stack e-commerce platform...",
        "content": "Detailed project description...",
        "featuredImage": "https://example.com/image.jpg",
        "status": "published",
        "publishedAt": "2024-01-01T00:00:00Z",
        "viewCount": 200,
        "likeCount": 35,
        "commentCount": 15,
        "liveUrl": "https://demo.example.com",
        "sourceUrl": "https://github.com/user/project",
        "demoUrl": "https://demo.example.com",
        "difficulty": "intermediate",
        "duration": "2-3 weeks",
        "teamSize": 1,
        "author": {
          "id": 1,
          "username": "author",
          "firstName": "John",
          "lastName": "Doe",
          "avatar": "https://example.com/avatar.jpg"
        },
        "categories": [...],
        "tags": [...],
        "technologies": [
          {
            "id": 1,
            "name": "React",
            "slug": "react",
            "icon": "react-icon.svg",
            "color": "#61DAFB"
          }
        ],
        "screenshots": [
          {
            "id": 1,
            "title": "Homepage",
            "imageUrl": "https://example.com/screenshot1.jpg",
            "altText": "Project homepage screenshot",
            "order": 1
          }
        ],
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ]
  },
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "pages": 5
    }
  }
}
```

### Get Project by Slug

```http
GET /api/v1/projects/{slug}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "project": {
      "id": 1,
      "uuid": "uuid-string",
      "title": "E-commerce Platform",
      "slug": "ecommerce-platform",
      "description": "A full-stack e-commerce platform...",
      "content": "Detailed project description...",
      "featuredImage": "https://example.com/image.jpg",
      "status": "published",
      "publishedAt": "2024-01-01T00:00:00Z",
      "viewCount": 200,
      "likeCount": 35,
      "commentCount": 15,
      "liveUrl": "https://demo.example.com",
      "sourceUrl": "https://github.com/user/project",
      "demoUrl": "https://demo.example.com",
      "difficulty": "intermediate",
      "duration": "2-3 weeks",
      "teamSize": 1,
      "author": {...},
      "categories": [...],
      "tags": [...],
      "technologies": [...],
      "screenshots": [...],
      "comments": [...],
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

## Comments Endpoints

### Get Comments for Post/Project

```http
GET /api/v1/posts/{slug}/comments
GET /api/v1/projects/{slug}/comments
```

**Query Parameters:**

- `page` (number) - Page number (default: 1)
- `limit` (number) - Items per page (default: 10)

**Response:**

```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": 1,
        "uuid": "uuid-string",
        "content": "Great article! Very helpful.",
        "status": "approved",
        "user": {
          "id": 1,
          "username": "commenter",
          "firstName": "Jane",
          "lastName": "Smith",
          "avatar": "https://example.com/avatar.jpg"
        },
        "parent": null,
        "children": [
          {
            "id": 2,
            "uuid": "uuid-string",
            "content": "I agree!",
            "status": "approved",
            "user": {
              "id": 2,
              "username": "replier",
              "firstName": "Bob",
              "lastName": "Johnson",
              "avatar": "https://example.com/avatar2.jpg"
            },
            "parent": {
              "id": 1,
              "uuid": "uuid-string"
            },
            "children": [],
            "createdAt": "2024-01-01T00:00:00Z",
            "updatedAt": "2024-01-01T00:00:00Z"
          }
        ],
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    ]
  },
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "pages": 3
    }
  }
}
```

### Create Comment (Authenticated)

```http
POST /api/v1/comments
```

**Request Body:**

```json
{
  "content": "Great article! Very helpful.",
  "postId": 1,
  "parentId": null
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "comment": {
      "id": 1,
      "uuid": "uuid-string",
      "content": "Great article! Very helpful.",
      "status": "pending",
      "user": {
        "id": 1,
        "username": "commenter",
        "firstName": "Jane",
        "lastName": "Smith",
        "avatar": "https://example.com/avatar.jpg"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  },
  "message": "Comment created successfully"
}
```

## User Profile Endpoints

### Get User Profile (Authenticated)

```http
GET /api/v1/profile
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "uuid": "uuid-string",
      "email": "user@example.com",
      "username": "username",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "https://example.com/avatar.jpg",
      "bio": "Software developer passionate about...",
      "role": "user",
      "status": "active",
      "verified": true,
      "lastLogin": "2024-01-01T00:00:00Z",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Update User Profile (Authenticated)

```http
PUT /api/v1/profile
```

**Request Body:**

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "bio": "Updated bio...",
  "username": "newusername"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "uuid": "uuid-string",
      "email": "user@example.com",
      "username": "newusername",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "https://example.com/avatar.jpg",
      "bio": "Updated bio...",
      "role": "user",
      "status": "active",
      "verified": true,
      "lastLogin": "2024-01-01T00:00:00Z",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  },
  "message": "Profile updated successfully"
}
```

## Search Endpoint

### Search Content

```http
GET /api/v1/search
```

**Query Parameters:**

- `q` (string) - Search query (required)
- `type` (string) - Content type (posts, projects, all)
- `page` (number) - Page number (default: 1)
- `limit` (number) - Items per page (default: 10)

**Response:**

```json
{
  "success": true,
  "data": {
    "results": {
      "posts": [...],
      "projects": [...]
    },
    "total": 150,
    "query": "golang"
  },
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 150,
      "pages": 15
    }
  }
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **General endpoints**: 100 requests per minute
- **Authentication endpoints**: 10 requests per minute
- **Search endpoints**: 30 requests per minute

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Error Codes

| Code                   | Description               |
| ---------------------- | ------------------------- |
| `VALIDATION_ERROR`     | Request validation failed |
| `AUTHENTICATION_ERROR` | Authentication failed     |
| `AUTHORIZATION_ERROR`  | Insufficient permissions  |
| `NOT_FOUND`            | Resource not found        |
| `DUPLICATE_ENTRY`      | Resource already exists   |
| `RATE_LIMIT_EXCEEDED`  | Rate limit exceeded       |
| `INTERNAL_ERROR`       | Internal server error     |

## WebSocket Events (Future)

Real-time features will be added in future versions:

- Live comments
- Real-time notifications
- Live collaboration
- Chat functionality

## SDKs and Libraries

Official SDKs will be provided for:

- JavaScript/TypeScript
- Python
- Go
- PHP
- Ruby

## Support

For API support and questions:

- Email: api@codewithdell.com
- Documentation: https://docs.codewithdell.com
- GitHub Issues: https://github.com/codewithdell/api/issues
