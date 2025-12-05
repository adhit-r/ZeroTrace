#!/usr/bin/env node
/**
 * Script to download Inter font files for self-hosting
 * Run: node scripts/download-fonts.js
 */

import https from 'https';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const fontsDir = path.join(__dirname, '../public/fonts');
const weights = [300, 400, 500, 600, 700];

// Ensure fonts directory exists
if (!fs.existsSync(fontsDir)) {
  fs.mkdirSync(fontsDir, { recursive: true });
}

// Inter font CDN URLs (using Google Fonts API)
const downloadFont = (weight) => {
  return new Promise((resolve, reject) => {
    const url = `https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_fvQtMwCp50KnMw2boKoduKmMEVuLyfAZ9hiJ-Ek-_EeA.woff2`;
    const filename = path.join(fontsDir, `inter-${weight}.woff2`);
    
    // For now, we'll use a placeholder approach
    // In production, download actual font files from Google Fonts or use a font service
    console.log(`Note: Font file for weight ${weight} should be downloaded from:`);
    console.log(`https://fonts.google.com/specimen/Inter`);
    console.log(`Or use: npm install @fontsource/inter`);
    console.log(`Then copy from node_modules/@fontsource/inter/files/`);
    
    resolve();
  });
};

// Alternative: Use @fontsource/inter package
console.log('To self-host Inter font:');
console.log('1. Run: npm install @fontsource/inter');
console.log('2. Copy font files from node_modules/@fontsource/inter/files/ to public/fonts/');
console.log('3. Or use the optimized Google Fonts link with display=swap');

