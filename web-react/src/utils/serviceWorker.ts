/**
 * Service Worker registration and management
 * Provides offline support and caching strategies
 */
import { Workbox } from 'workbox-window';

let wb: Workbox | null = null;
let registration: ServiceWorkerRegistration | null = null;

/**
 * Register service worker
 */
export async function registerServiceWorker(): Promise<ServiceWorkerRegistration | null> {
  if ('serviceWorker' in navigator) {
    try {
      wb = new Workbox('/sw.js', { type: 'module' });

      // Handle service worker updates
      wb.addEventListener('controlling', () => {
        window.location.reload();
      });

      // Register the service worker
      registration = (await wb.register()) || null;

      console.log('Service Worker registered successfully');

      // Check for updates periodically
      setInterval(() => {
        if (registration) {
          registration.update();
        }
      }, 60000); // Check every minute

      return registration;
    } catch (error) {
      console.error('Service Worker registration failed:', error);
      return null;
    }
  }

  return null;
}

/**
 * Unregister service worker
 */
export async function unregisterServiceWorker(): Promise<boolean> {
  if ('serviceWorker' in navigator && registration) {
    try {
      const success = await registration.unregister();
      if (success) {
        console.log('Service Worker unregistered successfully');
      }
      return success;
    } catch (error) {
      console.error('Service Worker unregistration failed:', error);
      return false;
    }
  }
  return false;
}

/**
 * Check if service worker is supported
 */
export function isServiceWorkerSupported(): boolean {
  return 'serviceWorker' in navigator;
}

/**
 * Get service worker registration
 */
export function getServiceWorkerRegistration(): ServiceWorkerRegistration | null {
  return registration;
}

