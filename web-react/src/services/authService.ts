import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    // For demo tokens, use a valid format for API
    if (token.startsWith('demo-token-')) {
      config.headers.Authorization = `Bearer demo-valid-token`;
    } else {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
    company_id: string;
  };
  token: string;
}

export const authService = {
  async login(email: string, password: string): Promise<LoginResponse> {
    // Demo mode - accept any credentials
    if (email && password) {
      return {
        user: {
          id: 'demo-user-1',
          email: email,
          name: 'Demo User',
          role: 'admin',
          company_id: 'demo-company-1',
        },
        token: 'demo-token-' + Date.now(),
      };
    }
    
    // Fallback to API call
    const response = await api.post('/api/v1/auth/login', { email, password });
    return response.data.data;
  },

  async register(userData: {
    email: string;
    password: string;
    name: string;
    company_id: string;
  }) {
    const response = await api.post('/api/v1/auth/register', userData);
    return response.data.data;
  },

  async validateToken(token: string) {
    const response = await api.get('/api/v1/auth/validate', {
      headers: { Authorization: `Bearer ${token}` },
    });
    return response.data.data;
  },
};
