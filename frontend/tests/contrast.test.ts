import { describe, it, expect } from "vitest";
import fs from "node:fs";
import path from "node:path";

function hexToRgb(hex: string): [number, number, number] {
  const clean = hex.replace("#", "").trim();
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

describe("WCAG AAA contrast verification", () => {
  const cssPath = path.resolve(__dirname, "../styles/tokens.css");
  const css = fs.readFileSync(cssPath, "utf-8");

  function getToken(name: string): string {
    const regex = new RegExp(`${name}:\\s*(#[a-fA-F0-9]{6})`);
    const match = css.match(regex);
    if (!match) throw new Error(`Token ${name} not found in tokens.css`);
    return match[1];
  }

  it("normal text tokens meet WCAG AAA 7:1 against page background", () => {
    const bgPage = getToken("--bg-page");
    const bgCard = getToken("--bg-card");
    const textMain = getToken("--text-main");
    const textMuted = getToken("--text-muted");

    expect(contrastRatio(textMain, bgPage)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textMain, bgCard)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textMuted, bgPage)).toBeGreaterThanOrEqual(7.0);
    expect(contrastRatio(textMuted, bgCard)).toBeGreaterThanOrEqual(7.0);
  });

  it("status alert chip tokens meet WCAG AAA", () => {
    const dangerText = getToken("--color-danger");
    const dangerBg = getToken("--color-danger-subtle");
    expect(contrastRatio(dangerText, dangerBg)).toBeGreaterThanOrEqual(7.0);

    const accentText = getToken("--color-accent");
    const accentBg = getToken("--color-accent-subtle");
    expect(contrastRatio(accentText, accentBg)).toBeGreaterThanOrEqual(7.0);

    const safeText = getToken("--color-safe");
    const safeBg = getToken("--color-safe-subtle");
    expect(contrastRatio(safeText, safeBg)).toBeGreaterThanOrEqual(7.0);
  });
});
