/**
 * Request deduplication utility
 * Prevents duplicate API calls within a specified timeframe
 * Integrates with axios interceptors
 */

interface PendingRequest {
  promise: Promise<any>;
  timestamp: number;
}

class RequestDeduplicator {
  private pendingRequests: Map<string, PendingRequest> = new Map();
  private readonly defaultTtl: number = 1000; // 1 second default TTL

  /**
   * Generate a cache key from request config
   */
  private getCacheKey(url: string, method: string, params?: any, data?: any): string {
    const paramsStr = params ? JSON.stringify(params) : '';
    const dataStr = data ? JSON.stringify(data) : '';
    return `${method.toUpperCase()}:${url}:${paramsStr}:${dataStr}`;
  }

  /**
   * Check if a request is pending and return it, or create a new one
   */
  deduplicate<T>(
    url: string,
    method: string,
    requestFn: () => Promise<T>,
    ttl: number = this.defaultTtl
  ): Promise<T> {
    const cacheKey = this.getCacheKey(url, method);

    // Check if there's a pending request
    const pending = this.pendingRequests.get(cacheKey);
    if (pending) {
      const age = Date.now() - pending.timestamp;
      if (age < ttl) {
        // Return existing pending request
        return pending.promise;
      } else {
        // Request is too old, remove it
        this.pendingRequests.delete(cacheKey);
      }
    }

    // Create new request
    const promise = requestFn()
      .then((response) => {
        // Remove from pending after completion
        this.pendingRequests.delete(cacheKey);
        return response;
      })
      .catch((error) => {
        // Remove from pending on error
        this.pendingRequests.delete(cacheKey);
        throw error;
      });

    // Store pending request
    this.pendingRequests.set(cacheKey, {
      promise,
      timestamp: Date.now(),
    });

    return promise;
  }

  /**
   * Clear all pending requests
   */
  clear(): void {
    this.pendingRequests.clear();
  }

  /**
   * Clear expired requests
   */
  clearExpired(ttl: number = this.defaultTtl): void {
    const now = Date.now();
    for (const [key, request] of this.pendingRequests.entries()) {
      if (now - request.timestamp > ttl) {
        this.pendingRequests.delete(key);
      }
    }
  }

  /**
   * Get number of pending requests
   */
  getPendingCount(): number {
    return this.pendingRequests.size;
  }
}

// Singleton instance
export const requestDeduplicator = new RequestDeduplicator();

/**
 * Batch multiple requests together
 */
export async function batchRequests<T>(
  requests: Array<() => Promise<T>>,
  maxConcurrent: number = 5
): Promise<T[]> {
  const results: T[] = [];
  const executing: Promise<void>[] = [];

  for (const request of requests) {
    const promise = request().then((result) => {
      results.push(result);
      executing.splice(executing.indexOf(promise), 1);
    });

    executing.push(promise);

    if (executing.length >= maxConcurrent) {
      await Promise.race(executing);
    }
  }

  await Promise.all(executing);
  return results;
}

/**
 * Queue requests with rate limiting
 */
export class RequestQueue {
  private queue: Array<() => Promise<any>> = [];
  private processing: boolean = false;
  private readonly delay: number;

  constructor(delay: number = 100) {
    this.delay = delay;
  }

  async add<T>(requestFn: () => Promise<T>): Promise<T> {
    return new Promise((resolve, reject) => {
      this.queue.push(async () => {
        try {
          const result = await requestFn();
          resolve(result);
        } catch (error) {
          reject(error);
        }
      });

      this.process();
    });
  }

  private async process(): Promise<void> {
    if (this.processing || this.queue.length === 0) {
      return;
    }

    this.processing = true;

    while (this.queue.length > 0) {
      const request = this.queue.shift();
      if (request) {
        await request();
        if (this.queue.length > 0) {
          await new Promise((resolve) => setTimeout(resolve, this.delay));
        }
      }
    }

    this.processing = false;
  }
}

