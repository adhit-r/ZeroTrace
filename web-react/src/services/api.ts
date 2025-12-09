import axios, { type InternalAxiosRequestConfig, type AxiosError } from 'axios';
import { requestDeduplicator } from '@/utils/requestDeduplication';

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token) {
      // For demo tokens, use a valid format for API
      if (token.startsWith('demo-token-')) {
        config.headers.Authorization = `Bearer demo-valid-token`;
      } else {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }

    // Request logging for debugging
    if (import.meta.env.DEV) {
      console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`, config.data || '');
    }

    // Only deduplicate GET requests by default
    if (config.method?.toLowerCase() === 'get' && !config.headers['X-Skip-Deduplication']) {
      const originalRequest = config;

      // Wrap the request in deduplication
      const deduplicatedRequest = requestDeduplicator.deduplicate(
        config.url || '',
        config.method,
        () => axios(originalRequest),
        2000 // 2 second TTL for GET requests
      );

      // Store the promise for the response interceptor
      (config as any).__deduplicatedPromise = deduplicatedRequest;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle deduplicated requests and logging
api.interceptors.response.use(
  (response) => {
    // Response logging for debugging
    if (import.meta.env.DEV) {
      console.log(`[API] ${response.config.method?.toUpperCase()} ${response.config.url} - Success`, response.data);
    }
    return response;
  },
  (error: AxiosError) => {
    // Error logging for debugging
    if (import.meta.env.DEV) {
      console.error(`[API] ${error.config?.method?.toUpperCase()} ${error.config?.url} - Error:`, error.response?.data || error.message);
    }
    // If this was a deduplicated request, the error is already handled
    return Promise.reject(error);
  }
);

// Retry interceptor with exponential backoff
api.interceptors.response.use(
  undefined,
  async (error: AxiosError) => {
    const config = error.config as InternalAxiosRequestConfig & { __retryCount?: number };

    // Don't retry if retry count is exceeded or if it's not a network error
    if (!config || (config.__retryCount || 0) >= 3 || !error.response) {
      return Promise.reject(error);
    }

    config.__retryCount = config.__retryCount || 0;
    config.__retryCount += 1;

    // Exponential backoff: 1s, 2s, 4s
    const delay = Math.pow(2, config.__retryCount - 1) * 1000;

    await new Promise((resolve) => setTimeout(resolve, delay));

    return api(config);
  }
);

