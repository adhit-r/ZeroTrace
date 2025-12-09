/**
 * Utility for lazy loading components with better error handling
 */
import { lazy, type ComponentType } from 'react';

/**
 * Lazy load a component with retry logic
 */
export function lazyLoad<T extends ComponentType<any>>(
  importFunc: () => Promise<{ default: T }>,
  retries = 3,
  delay = 1000
): React.LazyExoticComponent<T> {
  return lazy(() => {
    return new Promise<{ default: T }>((resolve, reject) => {
      const attemptImport = (attempt: number) => {
        importFunc()
          .then(resolve)
          .catch((error) => {
            if (attempt < retries) {
              setTimeout(() => attemptImport(attempt + 1), delay * attempt);
            } else {
              reject(error);
            }
          });
      };
      attemptImport(1);
    });
  });
}

/**
 * Lazy load with preload hint
 */
export function lazyLoadWithPreload<T extends ComponentType<any>>(
  importFunc: () => Promise<{ default: T }>
): React.LazyExoticComponent<T> & { preload: () => Promise<void> } {
  const LazyComponent = lazy(importFunc) as React.LazyExoticComponent<T> & {
    preload: () => Promise<void>;
  };

  LazyComponent.preload = async () => {
    await importFunc();
  };

  return LazyComponent;
}

