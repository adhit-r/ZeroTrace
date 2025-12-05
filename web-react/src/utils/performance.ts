/**
 * Performance monitoring utilities
 * Tracks Web Vitals and custom performance metrics
 */

interface WebVitals {
  lcp?: number; // Largest Contentful Paint
  fid?: number; // First Input Delay
  cls?: number; // Cumulative Layout Shift
  ttfb?: number; // Time to First Byte
  fcp?: number; // First Contentful Paint
}

class PerformanceMonitor {
  private vitals: WebVitals = {};
  private customMetrics: Map<string, number> = new Map();

  /**
   * Measure Web Vitals
   */
  measureWebVitals() {
    // Largest Contentful Paint (LCP)
    if ('PerformanceObserver' in window) {
      try {
        const lcpObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          const lastEntry = entries[entries.length - 1] as any;
          this.vitals.lcp = lastEntry.renderTime || lastEntry.loadTime;
          this.reportMetric('LCP', this.vitals.lcp);
        });
        lcpObserver.observe({ entryTypes: ['largest-contentful-paint'] });
      } catch (e) {
        console.warn('LCP observation not supported', e);
      }

      // First Input Delay (FID)
      try {
        const fidObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          entries.forEach((entry: any) => {
            if (entry.processingStart && entry.startTime) {
              this.vitals.fid = entry.processingStart - entry.startTime;
              this.reportMetric('FID', this.vitals.fid);
            }
          });
        });
        fidObserver.observe({ entryTypes: ['first-input'] });
      } catch (e) {
        console.warn('FID observation not supported', e);
      }

      // Cumulative Layout Shift (CLS)
      try {
        let clsValue = 0;
        const clsObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          entries.forEach((entry: any) => {
            if (!entry.hadRecentInput) {
              clsValue += entry.value;
              this.vitals.cls = clsValue;
            }
          });
          this.reportMetric('CLS', this.vitals.cls);
        });
        clsObserver.observe({ entryTypes: ['layout-shift'] });
      } catch (e) {
        console.warn('CLS observation not supported', e);
      }

      // Time to First Byte (TTFB)
      try {
        const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
        if (navigation) {
          this.vitals.ttfb = navigation.responseStart - navigation.requestStart;
          this.reportMetric('TTFB', this.vitals.ttfb);
        }
      } catch (e) {
        console.warn('TTFB measurement not supported', e);
      }

      // First Contentful Paint (FCP)
      try {
        const fcpObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          entries.forEach((entry: any) => {
            if (entry.name === 'first-contentful-paint') {
              this.vitals.fcp = entry.startTime;
              this.reportMetric('FCP', this.vitals.fcp);
            }
          });
        });
        fcpObserver.observe({ entryTypes: ['paint'] });
      } catch (e) {
        console.warn('FCP observation not supported', e);
      }
    }
  }

  /**
   * Measure custom metric
   */
  measureCustomMetric(name: string, value: number) {
    this.customMetrics.set(name, value);
    this.reportMetric(name, value);
  }

  /**
   * Start performance mark
   */
  mark(name: string) {
    if ('mark' in performance) {
      performance.mark(name);
    }
  }

  /**
   * Measure between two marks
   */
  measure(name: string, startMark: string, endMark?: string) {
    if ('measure' in performance) {
      try {
        performance.measure(name, startMark, endMark);
        const measure = performance.getEntriesByName(name)[0];
        if (measure) {
          this.measureCustomMetric(name, measure.duration);
        }
      } catch (e) {
        console.warn('Performance measure failed', e);
      }
    }
  }

  /**
   * Report metric (can be extended to send to analytics)
   */
  private reportMetric(name: string, value: number) {
    if (import.meta.env.DEV) {
      console.log(`[Performance] ${name}: ${value.toFixed(2)}ms`);
    }

    // In production, you might want to send this to an analytics service
    // Example: analytics.track('performance_metric', { name, value });
  }

  /**
   * Get all vitals
   */
  getVitals(): WebVitals {
    return { ...this.vitals };
  }

  /**
   * Get custom metrics
   */
  getCustomMetrics(): Record<string, number> {
    return Object.fromEntries(this.customMetrics);
  }

  /**
   * Get performance summary
   */
  getSummary() {
    return {
      vitals: this.getVitals(),
      customMetrics: this.getCustomMetrics(),
    };
  }
}

// Singleton instance
export const performanceMonitor = new PerformanceMonitor();

// Initialize on load
if (typeof window !== 'undefined') {
  if (document.readyState === 'complete') {
    performanceMonitor.measureWebVitals();
  } else {
    window.addEventListener('load', () => {
      performanceMonitor.measureWebVitals();
    });
  }
}

