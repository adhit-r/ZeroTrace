import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { registerServiceWorker } from './utils/serviceWorker'
import { performanceMonitor } from './utils/performance'

// Register service worker for PWA
if (import.meta.env.PROD) {
  registerServiceWorker().catch(console.error);
}

// Initialize performance monitoring
performanceMonitor.measureWebVitals();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
