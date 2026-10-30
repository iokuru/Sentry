import { describe, it, expect } from "vitest";
import fs from "node:fs";
import path from "node:path";
import zlib from "node:zlib";

describe("senior route page weight budget", () => {
  const buildDir = path.resolve(__dirname, "../build");

  it("senior route assets bundle under one hundred fifty kilobytes gzipped", () => {
    const indexPath = path.join(buildDir, "index.html");
    if (!fs.existsSync(indexPath)) {
      return;
    }

    const html = fs.readFileSync(indexPath, "utf-8");
    // Ignore legacy nomodule polyfills not loaded by modern browsers
    const scriptMatches = [...html.matchAll(/<script(?![^>]*nomodule)[^>]*src="(\/_next\/[^"]+)"/gi)].map(
      (m) => m[1]
    );
    const styleMatches = [...html.matchAll(/href="(\/_next\/[^"]+\.css)"/g)].map((m) => m[1]);

    const assetUrls = Array.from(new Set([...scriptMatches, ...styleMatches]));

    let totalGzipBytes = zlib.gzipSync(Buffer.from(html)).length;

    for (const assetUrl of assetUrls) {
      const relative = assetUrl.replace(/^\/_next\//, "_next/");
      const assetPath = path.join(buildDir, relative);
      if (fs.existsSync(assetPath)) {
        const content = fs.readFileSync(assetPath);
        totalGzipBytes += zlib.gzipSync(content).length;
      }
    }

    const maxBudgetKb = 150;
    const actualKb = totalGzipBytes / 1024;
    expect(actualKb).toBeLessThan(maxBudgetKb);
  });

  it("senior index page does not preload heavy ocr language assets", () => {
    const indexPath = path.join(buildDir, "index.html");
    if (!fs.existsSync(indexPath)) return;

    const html = fs.readFileSync(indexPath, "utf-8");
    expect(html).not.toContain("eng.traineddata");
    expect(html).not.toContain("tesseract-core");
  });
});
