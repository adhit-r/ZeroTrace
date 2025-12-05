/**
 * Hook for using Web Workers
 * Simplifies communication with workers
 */
import { useRef, useCallback, useEffect } from 'react';

interface WorkerMessage {
  type: string;
  data: any;
  id: string;
}

interface WorkerResponse {
  id: string;
  type: string;
  result?: any;
  error?: string;
  success: boolean;
}

export function useWebWorker(workerPath: string) {
  const workerRef = useRef<Worker | null>(null);
  const pendingRequestsRef = useRef<Map<string, { resolve: (value: any) => void; reject: (error: any) => void }>>(new Map());

  useEffect(() => {
    // Create worker
    workerRef.current = new Worker(workerPath, { type: 'module' });

    // Handle messages from worker
    workerRef.current.onmessage = (event: MessageEvent<WorkerResponse>) => {
      const { id, result, error, success } = event.data;
      const pending = pendingRequestsRef.current.get(id);

      if (pending) {
        if (success) {
          pending.resolve(result);
        } else {
          pending.reject(new Error(error || 'Worker error'));
        }
        pendingRequestsRef.current.delete(id);
      }
    };

    // Handle errors
    workerRef.current.onerror = (error) => {
      console.error('Worker error:', error);
      // Reject all pending requests
      pendingRequestsRef.current.forEach(({ reject }) => {
        reject(error);
      });
      pendingRequestsRef.current.clear();
    };

    // Cleanup
    return () => {
      if (workerRef.current) {
        workerRef.current.terminate();
        workerRef.current = null;
      }
      pendingRequestsRef.current.clear();
    };
  }, [workerPath]);

  const postMessage = useCallback(<T,>(type: string, data: any): Promise<T> => {
    return new Promise((resolve, reject) => {
      if (!workerRef.current) {
        reject(new Error('Worker not initialized'));
        return;
      }

      const id = `${Date.now()}-${Math.random()}`;
      pendingRequestsRef.current.set(id, { resolve, reject });

      workerRef.current.postMessage({
        type,
        data,
        id,
      } as WorkerMessage);

      // Timeout after 30 seconds
      setTimeout(() => {
        if (pendingRequestsRef.current.has(id)) {
          pendingRequestsRef.current.delete(id);
          reject(new Error('Worker request timeout'));
        }
      }, 30000);
    });
  }, []);

  return { postMessage };
}

