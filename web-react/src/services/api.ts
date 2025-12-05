import axios, { type InternalAxiosRequestConfig } from 'axios';
import { requestDeduplicator } from '@/utils/requestDeduplication';

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080'
});

// Request deduplication interceptor
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
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

// Response interceptor to handle deduplicated requests
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // If this was a deduplicated request, the error is already handled
    return Promise.reject(error);
  }
);

// Retry interceptor with exponential backoff
api.interceptors.response.use(
  undefined,
  async (error) => {
    const config = error.config;

    // Don't retry if retry count is exceeded or if it's not a network error
    if (!config || config.__retryCount >= 3 || !error.response) {
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

