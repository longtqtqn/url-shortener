import axios, { AxiosError } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Enhanced error helper function with detailed network error handling
export const getErrorMessage = (error: unknown): string => {
  if (error instanceof AxiosError) {
    // Check if we have a response from the server
    if (error.response) {
      // Server responded with error status
      const status = error.response.status;
      const data = error.response.data;
      
      // Try to get structured error message first
      if (data?.error?.message) {
        return data.error.message;
      }
      if (data?.error) {
        return typeof data.error === 'string' ? data.error : 'Server error occurred';
      }
      
      // Fallback to user-friendly HTTP status messages
      switch (status) {
        case 400:
          return 'Please check your input and try again.';
        case 401:
          return 'Please log in again to continue.';
        case 403:
          return 'You don\'t have permission for this action.';
        case 404:
          return 'The requested item was not found.';
        case 409:
          return 'This item already exists. Please try a different option.';
        case 429:
          return 'Too many requests. Please wait a moment and try again.';
        case 500:
        case 502:
        case 503:
        case 504:
          return 'Service is temporarily unavailable. Please try again later.';
        default:
          return 'Something went wrong. Please try again later.';
      }
    } else if (error.request) {
      // Request was made but no response received - Backend is down or network issues
      if (error.code === 'ECONNREFUSED' || error.code === 'ERR_NETWORK') {
        return 'Service is temporarily unavailable. Please try again in a few minutes.';
      }
      if (error.code === 'ENOTFOUND') {
        return 'Unable to connect. Please check your internet connection.';
      }
      if (error.code === 'ETIMEDOUT' || error.code === 'ECONNABORTED') {
        return 'Request timed out. Please try again.';
      }
      if (error.message.toLowerCase().includes('network')) {
        return 'Connection problem. Please check your internet and try again.';
      }
      return 'Service is currently unavailable. Please try again later.';
    } else {
      // Something else happened
      return 'Something went wrong. Please try again.';
    }
  }
  
  if (error instanceof Error) {
    return 'Something went wrong. Please try again.';
  }
  
  if (typeof error === 'string') {
    return error;
  }
  
  return 'Something went wrong. Please try again.';
};

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
  
  return config;
});

// Add response interceptor to handle common errors globally
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Add additional context for network errors
    if (error.code === 'ERR_NETWORK' && !error.response) {
      console.error('Backend server appears to be unavailable');
    }
    return Promise.reject(error);
  }
);

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
  api_key?: string; // For register response
  api_keys?: Array<{
    key: string;
    createdAt: string;
  }>; // For login response
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

  // Delete API key (JWT auth required)
  deleteApiKey: async (apiKey: string): Promise<{message: string}> => {
    const response = await api.delete('/api/api-key', {
      data: { api_key: apiKey }
    });
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
  deleteLink: async (shortCode: string, apiKey?: string): Promise<void> => {
    const config = apiKey ? {
      headers: {
        'X-API-KEY': apiKey
      }
    } : {};
    
    // Always use JWT route but pass API key in header for additional validation
    await api.delete(`/api/links/${shortCode}`, config);
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

  // API keys list management
  storeApiKeys: (apiKeys: Array<{key: string, createdAt: string}>): void => {
    localStorage.setItem('apiKeys', JSON.stringify(apiKeys));
    // Set the first API key as the current one if none is set
    if (!apiService.hasApiKey() && apiKeys.length > 0) {
      apiService.storeApiKey(apiKeys[0].key);
    }
  },

  getApiKeys: (): Array<{key: string, createdAt: string}> => {
    const stored = localStorage.getItem('apiKeys');
    return stored ? JSON.parse(stored) : [];
  },

  setSelectedApiKey: (apiKey: string): void => {
    apiService.storeApiKey(apiKey);
  },

  // Remove API key from stored list
  removeApiKeyFromList: (apiKeyToRemove: string): void => {
    const currentApiKeys = apiService.getApiKeys();
    const updatedApiKeys = currentApiKeys.filter(apiKey => apiKey.key !== apiKeyToRemove);
    localStorage.setItem('apiKeys', JSON.stringify(updatedApiKeys));
    
    // If the removed key was the currently selected one, clear it
    if (apiService.getApiKey() === apiKeyToRemove) {
      apiService.removeApiKey();
      // If there are other keys, select the first one
      if (updatedApiKeys.length > 0) {
        apiService.setSelectedApiKey(updatedApiKeys[0].key);
      }
    }
  },

  // Check if user is authenticated (either token or API key)
  isAuthenticated: (): boolean => {
    return apiService.hasToken() || apiService.hasApiKey();
  },

  // Logout (remove both token and API key)
  logout: (): void => {
    localStorage.removeItem('token');
    localStorage.removeItem('apiKey');
    localStorage.removeItem('apiKeys');
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
