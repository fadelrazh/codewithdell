// User types
export interface User {
  id: number;
  uuid: string;
  first_name: string;
  last_name: string;
  email: string;
  username: string;
  avatar?: string;
  bio?: string;
  website?: string;
  github?: string;
  twitter?: string;
  role: 'user' | 'admin';
  status: 'active' | 'inactive' | 'banned';
  created_at: string;
  updated_at: string;
}

// Post types
export interface Post {
  id: number;
  uuid: string;
  title: string;
  slug: string;
  content: string;
  excerpt?: string;
  featured_image?: string;
  status: 'draft' | 'published' | 'archived';
  published_at?: string;
  author_id: number;
  view_count: number;
  like_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
  author: User;
  categories: Category[];
  tags: Tag[];
}

// Category types
export interface Category {
  id: number;
  uuid: string;
  name: string;
  slug: string;
  description?: string;
  color?: string;
  icon?: string;
  created_at: string;
  updated_at: string;
}

// Tag types
export interface Tag {
  id: number;
  uuid: string;
  name: string;
  slug: string;
  color?: string;
  created_at: string;
  updated_at: string;
}

// Comment types
export interface Comment {
  id: number;
  uuid: string;
  content: string;
  user_id: number;
  post_id?: number;
  project_id?: number;
  parent_id?: number;
  status: 'pending' | 'approved' | 'spam';
  created_at: string;
  updated_at: string;
  user: User;
  children: Comment[];
}

// Project types
export interface Project {
  id: number;
  uuid: string;
  title: string;
  slug: string;
  description: string;
  content: string;
  featured_image?: string;
  status: 'draft' | 'published' | 'archived';
  published_at?: string;
  view_count: number;
  like_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
  technologies: Technology[];
  tags: Tag[];
  categories: Category[];
}

// Technology types
export interface Technology {
  id: number;
  uuid: string;
  name: string;
  slug: string;
  description?: string;
  icon?: string;
  color?: string;
  website?: string;
  created_at: string;
  updated_at: string;
}

// Authentication types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  username: string;
}

export interface AuthResponse {
  token: string;
  refresh_token: string;
  user: User;
}

// Post management types
export interface CreatePostRequest {
  title: string;
  content: string;
  excerpt?: string;
  slug?: string;
  status: 'draft' | 'published' | 'archived';
  tag_ids?: string[];
  category_ids?: string[];
}

export interface UpdatePostRequest {
  title?: string;
  content?: string;
  excerpt?: string;
  slug?: string;
  status?: 'draft' | 'published' | 'archived';
  tag_ids?: string[];
  category_ids?: string[];
}

export interface PostsResponse {
  posts: Post[];
  total: number;
  page: number;
  limit: number;
  pages: number;
}

// Comment types
export interface CreateCommentRequest {
  content: string;
  post_id?: number;
  project_id?: number;
  parent_id?: number;
}

export interface UpdateCommentRequest {
  content: string;
}

// Search types
export interface SearchRequest {
  query?: string;
  type?: 'posts' | 'projects' | 'all';
  category?: string;
  tags?: string[];
  author?: string;
  status?: string;
  sort_by?: 'relevance' | 'date' | 'views' | 'likes';
  sort_order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface SearchResponse {
  query: string;
  type: string;
  results: any;
  total: number;
  page: number;
  limit: number;
  pages: number;
}

// Analytics types
export interface AnalyticsResponse {
  overview: {
    total_posts: number;
    total_projects: number;
    total_users: number;
    total_comments: number;
    total_likes: number;
    total_bookmarks: number;
  };
  trends: {
    recent_posts: number;
    recent_projects: number;
    new_users: number;
    recent_comments: number;
    post_growth_rate: number;
  };
  popular: {
    popular_posts: Post[];
    most_liked_posts: Post[];
    popular_projects: Project[];
    popular_tags: Tag[];
  };
  user_stats: {
    active_users: number;
    top_contributors: User[];
    most_engaged_users: User[];
  };
  engagement: {
    avg_views_per_post: number;
    avg_likes_per_post: number;
    avg_comments_per_post: number;
    engagement_rate: number;
    total_views: number;
    total_likes: number;
    total_comments: number;
  };
}

// Upload types
export interface UploadResponse {
  url: string;
  filename: string;
  size: number;
  mime_type: string;
  uploaded_at: string;
}

// Profile types
export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  bio?: string;
  website?: string;
  github?: string;
  twitter?: string;
}

// Generic API response
export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  error?: string;
} 