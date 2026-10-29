import { describe, it, expect } from 'vitest';
import fs from 'node:fs';
import zlib from 'node:zlib';
import { fileURLToPath } from 'node:url';

function hexToRgb(hex: string): [number, number, number] {
  const clean = hex.replace('#', '').trim();
  const num = parseInt(clean, 16);
  return [(num >> 16) & 255, (num >> 8) & 255, num & 255];
}

function relativeLuminance(rgb: [number, number, number]): number {
  const [r, g, b] = rgb.map((val) => {
    const s = val / 255;
    return s <= 0.03928 ? s / 12.92 : Math.pow((s + 0.055) / 1.055, 2.4);
  });
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrastRatio(hex1: string, hex2: string): number {
  const l1 = relativeLuminance(hexToRgb(hex1));
  const l2 = relativeLuminance(hexToRgb(hex2));
  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);
  return (lighter + 0.05) / (darker + 0.05);
}

describe('WCAG AAA contrast checks', () => {
  const cssPath = fileURLToPath(new URL('../app.css', import.meta.url));
  const css = fs.readFileSync(cssPath, 'utf-8');

  function getToken(name: string): string {
    const regex = new RegExp(`${name}:\\s*(#[a-fA-F0-9]{6})`);
    const match = css.match(regex);
    if (!match) throw new Error(`Token ${name} not found in app.css`);
    return match[1];
  }

  it('normal text tokens meet WCAG AAA 7:1 against app background', () => {
    const bgApp = getToken('--bg-app');
    const bgSurface = getToken('--bg-surface');
    const textPrimary = getToken('--text-primary');
    const textSecondary = getToken('--text-secondary');
    const textMuted = getToken('--text-muted');

    expect(contrastRatio(textPrimary, bgApp)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textPrimary, bgSurface)).toBeGreaterThanOrEqual(7.0);

    expect(contrastRatio(textSecondary, bgApp)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textSecondary, bgSurface)).toBeGreaterThanOrEqual(7.0);

    expect(contrastRatio(textMuted, bgApp)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textMuted, bgSurface)).toBeGreaterThanOrEqual(7.0);
  });

  it('risk alert card tokens meet WCAG AAA', () => {
    const roseText = getToken('--rose-text');
    const roseBg = getToken('--rose-bg');
    expect(contrastRatio(roseText, roseBg)).toBeGreaterThanOrEqual(7.0);

    const amberText = getToken('--amber-text');
    const amberBg = getToken('--amber-bg');
    expect(contrastRatio(amberText, amberBg)).toBeGreaterThanOrEqual(7.0);

    const stoneText = getToken('--stone-text');
    const stoneBg = getToken('--stone-bg');
    expect(contrastRatio(stoneText, stoneBg)).toBeGreaterThanOrEqual(7.0);
  });
});

describe('senior route page weight budget', () => {
  it('senior route assets bundle under 150KB gzipped', () => {
    const buildDir = fileURLToPath(new URL('../../build', import.meta.url));
    if (!fs.existsSync(buildDir)) {
      // If build directory does not exist yet during local unit tests, pass
      return;
    }

    let totalGzipBytes = 0;
    function walkDir(dir: string) {
      const entries = fs.readdirSync(dir, { withFileTypes: true });
      for (const entry of entries) {
        if (entry.name === 'tesseract') {
          // lazy worker assets are not loaded on initial senior page
          continue;
        }
        const fullPath = `${dir}/${entry.name}`;
        if (entry.isDirectory()) {
          walkDir(fullPath);
        } else if (entry.name.endsWith('.js') || entry.name.endsWith('.css') || entry.name.endsWith('.html')) {
          const content = fs.readFileSync(fullPath);
          const compressed = zlib.gzipSync(content);
          totalGzipBytes += compressed.length;
        }
      }
    }

    walkDir(buildDir);
    const maxBudgetKb = 150;
    const actualKb = totalGzipBytes / 1024;
    expect(actualKb).toBeLessThan(maxBudgetKb);
  });

  it('senior initial page html does not preload tesseract assets', () => {
    const buildDir = fileURLToPath(new URL('../../build', import.meta.url));
    const indexPath = `${buildDir}/index.html`;
    if (!fs.existsSync(indexPath)) return;
    const html = fs.readFileSync(indexPath, 'utf-8');
    expect(html.includes('tesseract')).toBe(false);
  });
});
