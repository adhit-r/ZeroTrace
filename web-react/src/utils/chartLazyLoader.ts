/**
 * Lazy load Chart.js to reduce initial bundle size
 * Only loads when charts are actually needed
 */
let chartJsLoaded = false;
let chartJsPromise: Promise<void> | null = null;

export async function loadChartJS(): Promise<void> {
  if (chartJsLoaded) {
    return Promise.resolve();
  }

  if (chartJsPromise) {
    return chartJsPromise;
  }

  chartJsPromise = Promise.all([
    import('chart.js/auto'),
    import('react-chartjs-2'),
  ]).then(([chartModule, reactChartModule]) => {
    chartJsLoaded = true;
    return;
  });

  return chartJsPromise;
}

/**
 * Preload Chart.js when user hovers over chart container
 */
export function preloadChartJS(): void {
  if (!chartJsLoaded && !chartJsPromise) {
    // Start loading in background
    loadChartJS().catch(() => {
      // Silently fail - will retry on actual use
      chartJsPromise = null;
    });
  }
}

