/**
 * Lazy load D3.js to reduce initial bundle size
 * Only loads when topology/visualization components are needed
 */
let d3Loaded = false;
let d3Promise: Promise<void> | null = null;

export async function loadD3(): Promise<typeof import('d3')> {
  if (d3Loaded) {
    return import('d3');
  }

  if (d3Promise) {
    await d3Promise;
    return import('d3');
  }

  d3Promise = import('d3').then((d3) => {
    d3Loaded = true;
    return d3;
  });

  return d3Promise;
}

export async function loadD3Geo(): Promise<typeof import('d3-geo')> {
  return import('d3-geo');
}

/**
 * Preload D3 when user navigates to topology page
 */
export function preloadD3(): void {
  if (!d3Loaded && !d3Promise) {
    loadD3().catch(() => {
      d3Promise = null;
    });
  }
}

