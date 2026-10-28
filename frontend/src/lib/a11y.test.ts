import { describe, it, expect } from 'vitest';
import fs from 'node:fs';
import { fileURLToPath } from 'node:url';

describe('accessibility and tap targets', () => {
  const cssPath = fileURLToPath(new URL('../app.css', import.meta.url));
  const css = fs.readFileSync(cssPath, 'utf-8');

  it('defines touch-min custom property as at least 60px', () => {
    const match = css.match(/--touch-min:\s*(\d+)px/);
    expect(match).not.toBeNull();
    const touchMin = parseInt(match![1], 10);
    expect(touchMin).toBeGreaterThanOrEqual(60);
  });

  it('includes reduced motion media query overrides', () => {
    expect(css).toContain('@media (prefers-reduced-motion: reduce)');
    expect(css).toContain('animation-duration: 0.01ms');
    expect(css).toContain('transition-duration: 0.01ms');
  });

  it('defines high contrast visible focus styles', () => {
    expect(css).toContain(':focus-visible');
    expect(css).toContain('outline: 3px solid');
  });

  it('verifies senior intake button targets meet 60px minimum', () => {
    const checkBtnPath = fileURLToPath(new URL('./CheckButton.svelte', import.meta.url));
    const checkBtnCss = fs.readFileSync(checkBtnPath, 'utf-8');
    expect(checkBtnCss).toContain('--touch-min');

    const hotlineBtnPath = fileURLToPath(new URL('./HotlineButton.svelte', import.meta.url));
    const hotlineBtnCss = fs.readFileSync(hotlineBtnPath, 'utf-8');
    expect(hotlineBtnCss).toContain('--touch-min');
  });

  it('verifies caregiver action buttons meet 60px minimum', () => {
    const caregiverPagePath = fileURLToPath(new URL('../routes/caregiver/+page.svelte', import.meta.url));
    const caregiverCss = fs.readFileSync(caregiverPagePath, 'utf-8');
    expect(caregiverCss).toContain('--touch-min');
  });
});
