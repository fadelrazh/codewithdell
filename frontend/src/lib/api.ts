import {
  ApiResponse,
  AuthResponse,
  CreatePostRequest,
  LoginRequest,
  Post,
  PostsResponse,
  RegisterRequest,
  UpdatePostRequest,
  UpdateProfileRequest,
  User,
  Category,
  Tag,
  Comment,
  CreateCommentRequest,
  UpdateCommentRequest,
  SearchRequest,
  SearchResponse,
  AnalyticsResponse,
  UploadResponse,
} from '@/types/api';

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const API_ENDPOINTS = {
  // Auth
  AUTH: {
    LOGIN: `${API_BASE_URL}/api/v1/auth/login`,
    REGISTER: `${API_BASE_URL}/api/v1/auth/register`,
    REFRESH: `${API_BASE_URL}/api/v1/auth/refresh`,
  },
  
  // Posts
  POSTS: {
    LIST: `${API_BASE_URL}/api/v1/posts`,
    BY_SLUG: (slug: string) => `${API_BASE_URL}/api/v1/posts/${slug}`,
    CREATE: `${API_BASE_URL}/api/v1/admin/posts`,
    UPDATE: (id: number) => `${API_BASE_URL}/api/v1/admin/posts/${id}`,
    DELETE: (id: number) => `${API_BASE_URL}/api/v1/admin/posts/${id}`,
  },
  
  // Categories
  CATEGORIES: {
    LIST: `${API_BASE_URL}/api/v1/categories`,
    BY_SLUG: (slug: string) => `${API_BASE_URL}/api/v1/categories/${slug}`,
    POSTS: (slug: string) => `${API_BASE_URL}/api/v1/categories/${slug}/posts`,
    PROJECTS: (slug: string) => `${API_BASE_URL}/api/v1/categories/${slug}/projects`,
    CREATE: `${API_BASE_URL}/api/v1/admin/categories`,
    UPDATE: (id: number) => `${API_BASE_URL}/api/v1/admin/categories/${id}`,
    DELETE: (id: number) => `${API_BASE_URL}/api/v1/admin/categories/${id}`,
  },
  
  // Tags
  TAGS: {
    LIST: `${API_BASE_URL}/api/v1/tags`,
    POPULAR: `${API_BASE_URL}/api/v1/tags/popular`,
    BY_SLUG: (slug: string) => `${API_BASE_URL}/api/v1/tags/${slug}`,
    POSTS: (slug: string) => `${API_BASE_URL}/api/v1/tags/${slug}/posts`,
    PROJECTS: (slug: string) => `${API_BASE_URL}/api/v1/tags/${slug}/projects`,
    CREATE: `${API_BASE_URL}/api/v1/admin/tags`,
    UPDATE: (id: number) => `${API_BASE_URL}/api/v1/admin/tags/${id}`,
    DELETE: (id: number) => `${API_BASE_URL}/api/v1/admin/tags/${id}`,
  },
  
  // Comments
  COMMENTS: {
    LIST: `${API_BASE_URL}/api/v1/comments`,
    CREATE: `${API_BASE_URL}/api/v1/comments`,
    UPDATE: (id: number) => `${API_BASE_URL}/api/v1/comments/${id}`,
    DELETE: (id: number) => `${API_BASE_URL}/api/v1/comments/${id}`,
    PENDING: `${API_BASE_URL}/api/v1/admin/comments/pending`,
    APPROVE: (id: number) => `${API_BASE_URL}/api/v1/admin/comments/${id}/approve`,
    REJECT: (id: number) => `${API_BASE_URL}/api/v1/admin/comments/${id}/reject`,
  },
  
  // Interactions
  INTERACTIONS: {
    LIKE_POST: (id: number) => `${API_BASE_URL}/api/v1/interactions/posts/${id}/like`,
    UNLIKE_POST: (id: number) => `${API_BASE_URL}/api/v1/interactions/posts/${id}/like`,
    BOOKMARK_POST: (id: number) => `${API_BASE_URL}/api/v1/interactions/posts/${id}/bookmark`,
    REMOVE_BOOKMARK: (id: number) => `${API_BASE_URL}/api/v1/interactions/posts/${id}/bookmark`,
    CHECK_INTERACTION: (id: number) => `${API_BASE_URL}/api/v1/interactions/posts/${id}/check`,
    USER_LIKES: `${API_BASE_URL}/api/v1/interactions/likes`,
    USER_BOOKMARKS: `${API_BASE_URL}/api/v1/interactions/bookmarks`,
  },
  
  // Search
  SEARCH: {
    SEARCH: `${API_BASE_URL}/api/v1/search`,
    SUGGESTIONS: `${API_BASE_URL}/api/v1/search/suggestions`,
    STATS: `${API_BASE_URL}/api/v1/search/stats`,
  },
  
  // Analytics
  ANALYTICS: {
    OVERVIEW: `${API_BASE_URL}/api/v1/analytics`,
    POST_STATS: (id: number) => `${API_BASE_URL}/api/v1/analytics/posts/${id}`,
    USER_STATS: (id: number) => `${API_BASE_URL}/api/v1/analytics/users/${id}`,
  },
  
  // Upload
  UPLOAD: {
    IMAGE: `${API_BASE_URL}/api/v1/upload/image`,
    FILE: `${API_BASE_URL}/api/v1/upload/file`,
    DELETE: (filename: string) => `${API_BASE_URL}/api/v1/upload/${filename}`,
    STATS: `${API_BASE_URL}/api/v1/upload/stats`,
  },
  
  // Profile
  PROFILE: {
    GET: `${API_BASE_URL}/api/v1/profile`,
    UPDATE: `${API_BASE_URL}/api/v1/profile`,
  },
  
  // Health
  HEALTH: `${API_BASE_URL}/api/v1/test`,
};

export const fetchAPI = async (url: string, options?: RequestInit) => {
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...options?.headers,
    },
    mode: 'cors',
    ...options,
  });

  if (!response.ok) {
    throw new Error(`API request failed: ${response.status} ${response.statusText}`);
  }

  return response.json();
};

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    // Add auth token if available
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: `Bearer ${token}`,
        };
      }
    }

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Auth endpoints
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    return this.request<AuthResponse>('/api/v1/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  // Posts endpoints
  async getPosts(page = 1, limit = 10, filters?: {
    category?: string;
    tags?: string[];
    author?: string;
    status?: string;
  }): Promise<PostsResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    
    if (filters?.category) params.append('category', filters.category);
    if (filters?.tags) params.append('tags', filters.tags.join(','));
    if (filters?.author) params.append('author', filters.author);
    if (filters?.status) params.append('status', filters.status);
    
    return this.request<PostsResponse>(`/api/v1/posts?${params}`);
  }

  async getPostBySlug(slug: string): Promise<Post> {
    return this.request<Post>(`/api/v1/posts/${slug}`);
  }

  // Categories endpoints
  async getCategories(): Promise<{ categories: Category[]; total: number }> {
    return this.request<{ categories: Category[]; total: number }>('/api/v1/categories');
  }

  async getCategoryBySlug(slug: string): Promise<Category> {
    return this.request<Category>(`/api/v1/categories/${slug}`);
  }

  async getCategoryPosts(slug: string): Promise<{ category: Category; posts: Post[]; total: number }> {
    return this.request<{ category: Category; posts: Post[]; total: number }>(`/api/v1/categories/${slug}/posts`);
  }

  // Tags endpoints
  async getTags(): Promise<{ tags: Tag[]; total: number }> {
    return this.request<{ tags: Tag[]; total: number }>('/api/v1/tags');
  }

  async getPopularTags(): Promise<{ tags: Tag[]; total: number }> {
    return this.request<{ tags: Tag[]; total: number }>('/api/v1/tags/popular');
  }

  async getTagBySlug(slug: string): Promise<Tag> {
    return this.request<Tag>(`/api/v1/tags/${slug}`);
  }

  async getTagPosts(slug: string): Promise<{ tag: Tag; posts: Post[]; total: number }> {
    return this.request<{ tag: Tag; posts: Post[]; total: number }>(`/api/v1/tags/${slug}/posts`);
  }

  // Comments endpoints
  async getComments(postId?: number, projectId?: number): Promise<{ comments: Comment[]; total: number }> {
    const params = new URLSearchParams();
    if (postId) params.append('post_id', postId.toString());
    if (projectId) params.append('project_id', projectId.toString());
    
    return this.request<{ comments: Comment[]; total: number }>(`/api/v1/comments?${params}`);
  }

  async createComment(data: CreateCommentRequest): Promise<{ message: string; comment: Comment }> {
    return this.request<{ message: string; comment: Comment }>('/api/v1/comments', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateComment(id: number, data: UpdateCommentRequest): Promise<{ message: string; comment: Comment }> {
    return this.request<{ message: string; comment: Comment }>(`/api/v1/comments/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteComment(id: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/comments/${id}`, {
      method: 'DELETE',
    });
  }

  // Interactions endpoints
  async likePost(postId: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/interactions/posts/${postId}/like`, {
      method: 'POST',
    });
  }

  async unlikePost(postId: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/interactions/posts/${postId}/like`, {
      method: 'DELETE',
    });
  }

  async bookmarkPost(postId: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/interactions/posts/${postId}/bookmark`, {
      method: 'POST',
    });
  }

  async removeBookmark(postId: number): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/interactions/posts/${postId}/bookmark`, {
      method: 'DELETE',
    });
  }

  async checkUserInteraction(postId: number): Promise<{ is_liked: boolean; is_bookmarked: boolean }> {
    return this.request<{ is_liked: boolean; is_bookmarked: boolean }>(`/api/v1/interactions/posts/${postId}/check`);
  }

  async getUserLikes(): Promise<{ posts: Post[]; total: number }> {
    return this.request<{ posts: Post[]; total: number }>('/api/v1/interactions/likes');
  }

  async getUserBookmarks(): Promise<{ posts: Post[]; total: number }> {
    return this.request<{ posts: Post[]; total: number }>('/api/v1/interactions/bookmarks');
  }

  // Search endpoints
  async search(query: SearchRequest): Promise<SearchResponse> {
    const params = new URLSearchParams();
    if (query.query) params.append('q', query.query);
    if (query.type) params.append('type', query.type);
    if (query.category) params.append('category', query.category);
    if (query.tags) params.append('tags', query.tags.join(','));
    if (query.author) params.append('author', query.author);
    if (query.status) params.append('status', query.status);
    if (query.sort_by) params.append('sort_by', query.sort_by);
    if (query.sort_order) params.append('sort_order', query.sort_order);
    if (query.page) params.append('page', query.page.toString());
    if (query.limit) params.append('limit', query.limit.toString());
    
    return this.request<SearchResponse>(`/api/v1/search?${params}`);
  }

  async getSearchSuggestions(query: string): Promise<{ suggestions: string[] }> {
    return this.request<{ suggestions: string[] }>(`/api/v1/search/suggestions?q=${encodeURIComponent(query)}`);
  }

  async getSearchStats(): Promise<{
    total_posts: number;
    total_projects: number;
    total_tags: number;
    total_categories: number;
    popular_tags: Tag[];
  }> {
    return this.request<{
      total_posts: number;
      total_projects: number;
      total_tags: number;
      total_categories: number;
      popular_tags: Tag[];
    }>('/api/v1/search/stats');
  }

  // Analytics endpoints
  async getAnalytics(): Promise<AnalyticsResponse> {
    return this.request<AnalyticsResponse>('/api/v1/analytics');
  }

  async getPostStats(postId: number): Promise<{
    post_id: number;
    title: string;
    view_count: number;
    like_count: number;
    bookmark_count: number;
    comment_count: number;
    recent_comments: Comment[];
  }> {
    return this.request<{
      post_id: number;
      title: string;
      view_count: number;
      like_count: number;
      bookmark_count: number;
      comment_count: number;
      recent_comments: Comment[];
    }>(`/api/v1/analytics/posts/${postId}`);
  }

  // Upload endpoints
  async uploadImage(file: File): Promise<{ message: string; data: UploadResponse }> {
    const formData = new FormData();
    formData.append('image', file);
    
    const token = localStorage.getItem('token');
    const response = await fetch(`${this.baseURL}/api/v1/upload/image`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`);
    }
    
    return response.json();
  }

  async uploadFile(file: File): Promise<{ message: string; data: UploadResponse }> {
    const formData = new FormData();
    formData.append('file', file);
    
    const token = localStorage.getItem('token');
    const response = await fetch(`${this.baseURL}/api/v1/upload/file`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`);
    }
    
    return response.json();
  }

  async deleteFile(filename: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/api/v1/upload/${filename}`, {
      method: 'DELETE',
    });
  }

  async getUploadStats(): Promise<{
    images: { count: number; size: number };
    files: { count: number; size: number };
    total: { count: number; size: number };
  }> {
    return this.request<{
      images: { count: number; size: number };
      files: { count: number; size: number };
      total: { count: number; size: number };
    }>('/api/v1/upload/stats');
  }

  // Protected endpoints
  async getProfile(): Promise<User> {
    return this.request<User>('/api/v1/profile');
  }

  async updateProfile(data: UpdateProfileRequest): Promise<User> {
    return this.request<User>('/api/v1/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Admin endpoints
  async createPost(data: CreatePostRequest): Promise<Post> {
    return this.request<Post>('/api/v1/admin/posts', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updatePost(id: number, data: UpdatePostRequest): Promise<Post> {
    return this.request<Post>(`/api/v1/admin/posts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deletePost(id: number): Promise<void> {
    return this.request<void>(`/api/v1/admin/posts/${id}`, {
      method: 'DELETE',
    });
  }

  // Health check
  async healthCheck(): Promise<{ message: string; service: string; timestamp: string }> {
    return this.request('/api/v1/test');
  }
}

export const apiClient = new ApiClient(API_BASE_URL); 