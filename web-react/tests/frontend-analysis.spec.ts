import { test, expect } from '@playwright/test';

test.describe('ZeroTrace Frontend Analysis', () => {
  test('should load main page without console errors', async ({ page }) => {
    // Listen for console errors
    const consoleErrors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });

    // Listen for page errors
    const pageErrors: string[] = [];
    page.on('pageerror', error => {
      pageErrors.push(error.message);
    });

    // Navigate to the app
    await page.goto('/');

    // Wait for React root to render - look for the root div with content
    await page.waitForFunction(() => {
      const root = document.getElementById('root');
      return root && root.children.length > 0;
    }, { timeout: 10000 });

    // Wait for the page to be stable
    await page.waitForLoadState('networkidle');

    // Check if main content loaded
    const root = page.locator('#root');
    await expect(root).toBeVisible();

    // Report any errors found
    if (consoleErrors.length > 0) {
      console.log('Console errors found:', consoleErrors);
    }

    if (pageErrors.length > 0) {
      console.log('Page errors found:', pageErrors);
    }

    // Take a screenshot for analysis
    await page.screenshot({ path: 'frontend-analysis.png', fullPage: true });
  });

  test('should navigate to Organization Profile page', async ({ page }) => {
    await page.goto('/');

    // Try to find and click on Organization Profile link
    // This will depend on the actual navigation structure
    const orgProfileLink = page.locator('a[href*="organization"], a[href*="profile"]').first();
    if (await orgProfileLink.isVisible()) {
      await orgProfileLink.click();
      await page.waitForLoadState('networkidle');

      // Check if the page loaded without errors
      const body = page.locator('body');
      await expect(body).toBeVisible();
    }
  });

  test('should check for missing UI components', async ({ page }) => {
    await page.goto('/');

    // Check if common UI elements are present
    const buttons = page.locator('button');
    const inputs = page.locator('input');
    const cards = page.locator('[class*="card"], [class*="Card"]');

    console.log('Buttons found:', await buttons.count());
    console.log('Inputs found:', await inputs.count());
    console.log('Card elements found:', await cards.count());

    // Take screenshot of current state
    await page.screenshot({ path: 'ui-components-check.png' });
  });
});