import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add JWT token and/or API key to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  const apiKey = localStorage.getItem('apiKey');
  
  // Add JWT token if available
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  
  // Add API key if available (server can use either JWT or API key)
  if (apiKey) {
    config.headers['X-API-KEY'] = apiKey;
  }
  
  // Debug logging (can be removed in production)
  if (import.meta.env.DEV) {
    console.log('API Request:', {
      url: config.url,
      method: config.method,
      hasToken: !!token,
      hasApiKey: !!apiKey,
      headers: {
        Authorization: config.headers['Authorization'] ? 'Bearer [REDACTED]' : undefined,
        'X-API-KEY': config.headers['X-API-KEY'] ? apiKey : undefined,
      }
    });
  }
  
  return config;
});

// Auth interfaces
export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  message: string;
  token: string;
  api_key?: string; // API key might be included in future
}

export interface CreateApiKeyResponse {
  message: string;
  api_key: string;
}

// Link interfaces
export interface CreateLinkRequest {
  long_url: string;
  short_code?: string;
}

export interface CreateLinkResponse {
  shortened_url: string;
  short_code: string;
  long_url: string;
}

export interface Link {
  shortURL: string;
  longURL: string;
  clickCount: number;
  lastClicked: string | null;
  createdAt: string;
  apiKey?: string;
}

export const apiService = {
  // Authentication methods
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post('/register', data);
    return response.data;
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/login', data);
    return response.data;
  },

  // Create API key (JWT auth required)
  createApiKey: async (): Promise<CreateApiKeyResponse> => {
    const response = await api.post('/api/create-api-key');
    return response.data;
  },

  // Create short link (auth required)
  createLink: async (data: CreateLinkRequest): Promise<CreateLinkResponse> => {
    const response = await api.post('/api/links', data);
    return response.data;
  },

  // Create short link (public, no auth required)
  createLinkPublic: async (data: CreateLinkRequest): Promise<CreateLinkResponse> => {
    const response = await api.post('/shorten', data);
    return response.data;
  },

  // Get user's links (auth required)
  getUserLinks: async (): Promise<Link[]> => {
    const response = await api.get('/api/links');
    return response.data;
  },

  // Delete link (auth required)
  deleteLink: async (shortCode: string): Promise<void> => {
    await api.delete(`/api/links/${shortCode}`);
  },

  // Token management
  hasToken: (): boolean => {
    return !!localStorage.getItem('token');
  },

  getToken: (): string | null => {
    return localStorage.getItem('token');
  },

  storeToken: (token: string): void => {
    localStorage.setItem('token', token);
  },

  removeToken: (): void => {
    localStorage.removeItem('token');
  },

  // API key management (legacy support)
  hasApiKey: (): boolean => {
    return !!localStorage.getItem('apiKey');
  },

  getApiKey: (): string | null => {
    return localStorage.getItem('apiKey');
  },

  storeApiKey: (apiKey: string): void => {
    localStorage.setItem('apiKey', apiKey);
  },

  removeApiKey: (): void => {
    localStorage.removeItem('apiKey');
  },

  // Check if user is authenticated (either token or API key)
  isAuthenticated: (): boolean => {
    return apiService.hasToken() || apiService.hasApiKey();
  },

  // Logout (remove both token and API key)
  logout: (): void => {
    localStorage.removeItem('token');
    localStorage.removeItem('apiKey');
  },

  // Debug utility to check current authentication headers
  getAuthHeaders: (): { hasToken: boolean; hasApiKey: boolean; headers: Record<string, string> } => {
    const token = localStorage.getItem('token');
    const apiKey = localStorage.getItem('apiKey');
    const headers: Record<string, string> = {};

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    if (apiKey) {
      headers['X-API-KEY'] = apiKey;
    }

    return {
      hasToken: !!token,
      hasApiKey: !!apiKey,
      headers
    };
  },
};

export default api;
