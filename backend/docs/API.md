# CodeWithDell API Documentation

## Overview

The CodeWithDell API is a RESTful API for managing blog posts, projects, user interactions, and content management. This API provides comprehensive functionality for a YouTube project showcase platform.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Public endpoints**: 100 requests per minute per IP
- **Authenticated endpoints**: 200 requests per minute per user
- **Admin endpoints**: 500 requests per minute per user

Rate limit headers are included in responses:

- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the rate limit resets (Unix timestamp)

## Endpoints

### Authentication

#### Register User

```http
POST /auth/register
```

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "Password123!",
  "username": "johndoe"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "username": "johndoe",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login User

```http
POST /auth/login
```

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "Password123!"
}
```

**Response:** Same as register response

#### Refresh Token

```http
POST /auth/refresh
```

**Headers:**

```
Authorization: Bearer <refresh-token>
```

**Response:**

```json
{
  "token": "new-access-token",
  "refresh_token": "new-refresh-token"
}
```

### Posts

#### Get All Posts

```http
GET /posts?page=1&limit=10&category=web-development&tags=javascript,react
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 50)
- `category` (optional): Filter by category slug
- `tags` (optional): Comma-separated tag slugs
- `author` (optional): Filter by author username
- `status` (optional): Filter by status (published, draft, archived)

**Response:**

```json
{
  "posts": [
    {
      "id": 1,
      "title": "Building a Modern Web App",
      "slug": "building-modern-web-app",
      "excerpt": "Learn how to build a modern web application...",
      "content": "Full post content...",
      "status": "published",
      "view_count": 150,
      "like_count": 25,
      "comment_count": 8,
      "author": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "username": "johndoe"
      },
      "tags": [
        {
          "id": 1,
          "name": "JavaScript",
          "slug": "javascript"
        }
      ],
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 10,
  "pages": 5
}
```

#### Get Post by Slug

```http
GET /posts/building-modern-web-app
```

**Response:**

```json
{
  "id": 1,
  "title": "Building a Modern Web App",
  "slug": "building-modern-web-app",
  "content": "Full post content...",
  "excerpt": "Learn how to build a modern web application...",
  "featured_image": "/uploads/images/featured.jpg",
  "status": "published",
  "published_at": "2024-01-01T00:00:00Z",
  "view_count": 150,
  "like_count": 25,
  "comment_count": 8,
  "author": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe"
  },
  "categories": [
    {
      "id": 1,
      "name": "Web Development",
      "slug": "web-development"
    }
  ],
  "tags": [
    {
      "id": 1,
      "name": "JavaScript",
      "slug": "javascript"
    }
  ],
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Comments

#### Get Comments

```http
GET /comments?post_id=1&project_id=2
```

**Query Parameters:**

- `post_id` (optional): Get comments for specific post
- `project_id` (optional): Get comments for specific project

**Response:**

```json
{
  "comments": [
    {
      "id": 1,
      "content": "Great article! Very helpful.",
      "status": "approved",
      "created_at": "2024-01-01T00:00:00Z",
      "user": {
        "id": 2,
        "first_name": "Jane",
        "last_name": "Smith",
        "username": "janesmith"
      },
      "children": []
    }
  ],
  "total": 1
}
```

#### Create Comment (Authenticated)

```http
POST /comments
```

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "content": "Great article! Very helpful.",
  "post_id": 1,
  "parent_id": null
}
```

**Response:**

```json
{
  "message": "Comment created successfully and pending approval",
  "comment": {
    "id": 1,
    "content": "Great article! Very helpful.",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
    "user": {
      "id": 2,
      "first_name": "Jane",
      "last_name": "Smith",
      "username": "janesmith"
    }
  }
}
```

### User Interactions

#### Like Post (Authenticated)

```http
POST /interactions/posts/1/like
```

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
{
  "message": "Post liked successfully"
}
```

#### Unlike Post (Authenticated)

```http
DELETE /interactions/posts/1/like
```

**Headers:**

```
Authorization: Bearer <token>
```

#### Bookmark Post (Authenticated)

```http
POST /interactions/posts/1/bookmark
```

**Headers:**

```
Authorization: Bearer <token>
```

#### Remove Bookmark (Authenticated)

```http
DELETE /interactions/posts/1/bookmark
```

**Headers:**

```
Authorization: Bearer <token>
```

#### Check User Interaction (Authenticated)

```http
GET /interactions/posts/1/check
```

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
{
  "is_liked": true,
  "is_bookmarked": false
}
```

#### Get User Likes (Authenticated)

```http
GET /interactions/likes
```

**Headers:**

```
Authorization: Bearer <token>
```

#### Get User Bookmarks (Authenticated)

```http
GET /interactions/bookmarks
```

**Headers:**

```
Authorization: Bearer <token>
```

### Categories

#### Get All Categories

```http
GET /categories
```

**Response:**

```json
{
  "categories": [
    {
      "id": 1,
      "name": "Web Development",
      "slug": "web-development",
      "description": "Web development tutorials and guides",
      "color": "#3B82F6",
      "icon": "code"
    }
  ],
  "total": 1
}
```

#### Get Category by Slug

```http
GET /categories/web-development
```

#### Get Posts by Category

```http
GET /categories/web-development/posts
```

#### Get Projects by Category

```http
GET /categories/web-development/projects
```

### Tags

#### Get All Tags

```http
GET /tags
```

#### Get Popular Tags

```http
GET /tags/popular
```

#### Get Tag by Slug

```http
GET /tags/javascript
```

#### Get Posts by Tag

```http
GET /tags/javascript/posts
```

#### Get Projects by Tag

```http
GET /tags/javascript/projects
```

### Search

#### Search Content

```http
GET /search?q=javascript&type=posts&category=web-development&tags=react&sort_by=relevance&page=1&limit=10
```

**Query Parameters:**

- `q` (optional): Search query
- `type` (optional): Search type (posts, projects, all)
- `category` (optional): Filter by category
- `tags` (optional): Comma-separated tag slugs
- `author` (optional): Filter by author
- `status` (optional): Filter by status
- `sort_by` (optional): Sort by (relevance, date, views, likes)
- `sort_order` (optional): Sort order (asc, desc)
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response:**

```json
{
  "query": "javascript",
  "type": "posts",
  "results": [
    {
      "id": 1,
      "title": "JavaScript Fundamentals",
      "slug": "javascript-fundamentals",
      "excerpt": "Learn JavaScript basics...",
      "author": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "username": "johndoe"
      },
      "tags": [
        {
          "id": 1,
          "name": "JavaScript",
          "slug": "javascript"
        }
      ],
      "view_count": 150,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 10,
  "pages": 3
}
```

#### Get Search Suggestions

```http
GET /search/suggestions?q=jav
```

**Response:**

```json
{
  "suggestions": ["JavaScript Fundamentals", "Java Programming", "JavaScript"]
}
```

#### Get Search Statistics

```http
GET /search/stats
```

**Response:**

```json
{
  "total_posts": 150,
  "total_projects": 75,
  "total_tags": 50,
  "total_categories": 10,
  "popular_tags": [
    {
      "id": 1,
      "name": "JavaScript",
      "slug": "javascript",
      "usage_count": 45
    }
  ]
}
```

### Analytics

#### Get Analytics Overview

```http
GET /analytics
```

**Response:**

```json
{
  "overview": {
    "total_posts": 150,
    "total_projects": 75,
    "total_users": 500,
    "total_comments": 1200,
    "total_likes": 3500,
    "total_bookmarks": 800
  },
  "trends": {
    "recent_posts": 15,
    "recent_projects": 8,
    "new_users": 25,
    "recent_comments": 45,
    "post_growth_rate": 12.5
  },
  "popular": {
    "popular_posts": [...],
    "most_liked_posts": [...],
    "popular_projects": [...],
    "popular_tags": [...]
  },
  "user_stats": {
    "active_users": 150,
    "top_contributors": [...],
    "most_engaged_users": [...]
  },
  "engagement": {
    "avg_views_per_post": 125.5,
    "avg_likes_per_post": 23.3,
    "avg_comments_per_post": 8.0,
    "engagement_rate": 25.0,
    "total_views": 18825,
    "total_likes": 3500,
    "total_comments": 1200
  }
}
```

#### Get Post Statistics

```http
GET /analytics/posts/1
```

#### Get User Statistics

```http
GET /analytics/users/1
```

### File Upload (Authenticated)

#### Upload Image

```http
POST /upload/image
```

**Headers:**

```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data:**

- `image`: Image file (max 5MB, formats: JPEG, PNG, GIF, WebP)

**Response:**

```json
{
  "message": "Image uploaded successfully",
  "data": {
    "url": "/uploads/images/abc123.jpg",
    "filename": "abc123.jpg",
    "size": 1024000,
    "mime_type": "image/jpeg",
    "uploaded_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Upload File

```http
POST /upload/file
```

**Headers:**

```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data:**

- `file`: File (max 10MB, formats: PDF, DOC, DOCX, TXT, ZIP, RAR, MP4, MP3)

#### Delete File (Authenticated)

```http
DELETE /upload/abc123.jpg
```

**Headers:**

```
Authorization: Bearer <token>
```

#### Get Upload Statistics (Authenticated)

```http
GET /upload/stats
```

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
{
  "images": {
    "count": 25,
    "size": 52428800
  },
  "files": {
    "count": 10,
    "size": 104857600
  },
  "total": {
    "count": 35,
    "size": 157286400
  }
}
```

### User Profile (Authenticated)

#### Get Profile

```http
GET /profile
```

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "username": "johndoe",
  "avatar": "/uploads/images/avatar.jpg",
  "bio": "Web developer and content creator",
  "website": "https://johndoe.com",
  "github": "johndoe",
  "twitter": "johndoe",
  "role": "user",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Update Profile

```http
PUT /profile
```

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Updated bio",
  "website": "https://johndoe.com",
  "github": "johndoe",
  "twitter": "johndoe"
}
```

### Admin Endpoints

All admin endpoints require admin role authentication.

#### Create Post (Admin)

```http
POST /admin/posts
```

**Headers:**

```
Authorization: Bearer <admin-token>
```

**Request Body:**

```json
{
  "title": "New Post Title",
  "content": "Post content...",
  "excerpt": "Post excerpt...",
  "slug": "new-post-slug",
  "status": "published",
  "tag_ids": ["1", "2"]
}
```

#### Update Post (Admin)

```http
PUT /admin/posts/1
```

#### Delete Post (Admin)

```http
DELETE /admin/posts/1
```

#### Create Category (Admin)

```http
POST /admin/categories
```

**Request Body:**

```json
{
  "name": "New Category",
  "description": "Category description",
  "color": "#3B82F6",
  "icon": "code"
}
```

#### Update Category (Admin)

```http
PUT /admin/categories/1
```

#### Delete Category (Admin)

```http
DELETE /admin/categories/1
```

#### Create Tag (Admin)

```http
POST /admin/tags
```

**Request Body:**

```json
{
  "name": "New Tag",
  "color": "#EF4444"
}
```

#### Update Tag (Admin)

```http
PUT /admin/tags/1
```

#### Delete Tag (Admin)

```http
DELETE /admin/tags/1
```

#### Get Pending Comments (Admin)

```http
GET /admin/comments/pending
```

#### Approve Comment (Admin)

```http
POST /admin/comments/1/approve
```

#### Reject Comment (Admin)

```http
POST /admin/comments/1/reject
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Validation Error
- `429` - Too Many Requests
- `500` - Internal Server Error

## Validation Rules

### User Registration/Login

- Email: Valid email format
- Password: Minimum 8 characters, must contain uppercase, lowercase, digit, and special character
- Username: 3-30 characters, letters, numbers, underscores, hyphens only
- First/Last Name: 2-50 characters

### Posts

- Title: 3-200 characters
- Content: 10-50,000 characters
- Slug: 3-100 characters, lowercase letters, numbers, hyphens only
- Status: draft, published, or archived

### Comments

- Content: 1-1,000 characters

### Categories

- Name: 2-50 characters
- Description: Maximum 200 characters
- Color: Valid hex color code
- Icon: Maximum 50 characters

### Tags

- Name: 2-30 characters
- Color: Valid hex color code

## Pagination

All list endpoints support pagination with these query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, maximum: 50)

Response includes pagination metadata:

- `total`: Total number of items
- `page`: Current page number
- `limit`: Items per page
- `pages`: Total number of pages

## File Upload Limits

- **Images**: Maximum 5MB, formats: JPEG, PNG, GIF, WebP
- **Files**: Maximum 10MB, formats: PDF, DOC, DOCX, TXT, ZIP, RAR, MP4, MP3

## WebSocket Endpoints

For real-time features, WebSocket connections are available at:

```
ws://localhost:8080/ws
```

### WebSocket Events

#### Client to Server

- `join_room`: Join a specific room (e.g., post comments)
- `leave_room`: Leave a room
- `new_comment`: Send a new comment
- `like_post`: Like/unlike a post

#### Server to Client

- `comment_added`: New comment added
- `comment_updated`: Comment updated
- `comment_deleted`: Comment deleted
- `post_liked`: Post like count updated
- `user_online`: User came online
- `user_offline`: User went offline

## SDKs and Libraries

### JavaScript/TypeScript

```bash
npm install codewithdell-api-client
```

```javascript
import { CodeWithDellAPI } from "codewithdell-api-client";

const api = new CodeWithDellAPI({
  baseURL: "http://localhost:8080/api/v1",
  token: "your-jwt-token",
});

// Get posts
const posts = await api.posts.getAll({ page: 1, limit: 10 });

// Create comment
const comment = await api.comments.create({
  content: "Great post!",
  post_id: 1,
});
```

### Python

```bash
pip install codewithdell-python
```

```python
from codewithdell import CodeWithDellAPI

api = CodeWithDellAPI(
    base_url='http://localhost:8080/api/v1',
    token='your-jwt-token'
)

# Get posts
posts = api.posts.get_all(page=1, limit=10)

# Create comment
comment = api.comments.create(
    content='Great post!',
    post_id=1
)
```

## Support

For API support and questions:

- Email: api-support@codewithdell.com
- Documentation: https://docs.codewithdell.com/api
- GitHub Issues: https://github.com/codewithdell/api/issues
