import fs from "fs";
import path from "path";
import { createWorker } from "tesseract.js";

function computeSimilarity(s1, s2) {
  const norm1 = s1.toLowerCase().replace(/[^a-z0-9]/g, "");
  const norm2 = s2.toLowerCase().replace(/[^a-z0-9]/g, "");
  if (!norm1 && !norm2) return 1.0;
  if (!norm1 || !norm2) return 0.0;

  const bigrams = (str) => {
    const set = new Map();
    for (let i = 0; i < str.length - 1; i++) {
      const b = str.substring(i, i + 2);
      set.set(b, (set.get(b) || 0) + 1);
    }
    return set;
  };

  const b1 = bigrams(norm1);
  const b2 = bigrams(norm2);
  let intersection = 0;
  for (const [k, count] of b1.entries()) {
    if (b2.has(k)) {
      intersection += Math.min(count, b2.get(k));
    }
  }

  const total = (norm1.length - 1) + (norm2.length - 1);
  return total > 0 ? (2.0 * intersection) / total : (norm1 === norm2 ? 1.0 : 0.0);
}

async function main() {
  const rootDir = process.cwd();
  const screenshotsDir = path.join(rootDir, "testdata", "screenshots");
  const manifestPath = path.join(screenshotsDir, "manifest.json");

  if (!fs.existsSync(manifestPath)) {
    console.error("Manifest not found at", manifestPath);
    process.exit(1);
  }

  const rawData = fs.readFileSync(manifestPath, "utf-8").replace(/^\uFEFF/, "");
  const manifest = JSON.parse(rawData);
  console.log(`Starting OCR evaluation on ${manifest.length} screenshot fixtures...`);

  const worker = await createWorker("eng");
  const results = [];

  for (const item of manifest) {
    const imgPath = path.join(screenshotsDir, item.file);
    const start = Date.now();
    const ret = await worker.recognize(imgPath);
    const elapsed = Date.now() - start;
    const recognizedText = (ret.data.text || "").trim();

    let sim = 0;
    if (item.is_blank) {
      sim = recognizedText.length < 20 ? 1.0 : 0.5;
    } else {
      sim = computeSimilarity(item.expected_text, recognizedText);
    }

    results.push({
      id: item.id,
      file: item.file,
      expected: item.expected_text,
      recognized: recognizedText,
      similarity: sim,
      durationMs: elapsed,
      isBlank: item.is_blank,
    });

    const status = sim >= 0.70 ? "PASS" : "WARN";
    console.log(`[${status}] ${item.id} - Sim: ${(sim * 100).toFixed(1)}% (${elapsed}ms)`);
  }

  await worker.terminate();

  const textResults = results.filter((r) => !r.isBlank);
  const avgSim = textResults.reduce((acc, r) => acc + r.similarity, 0) / textResults.length;
  const passed = results.filter((r) => r.similarity >= 0.70).length;

  console.log("\n================ OCR EVALUATION SUMMARY ================");
  console.log(`Total fixtures:       ${results.length}`);
  console.log(`Text fixtures:        ${textResults.length}`);
  console.log(`Blank/photo fixtures: ${results.filter((r) => r.isBlank).length}`);
  console.log(`Passing fixtures:     ${passed}/${results.length} (${((passed / results.length) * 100).toFixed(1)}%)`);
  console.log(`Avg Text Similarity:  ${(avgSim * 100).toFixed(1)}%`);

  fs.writeFileSync(
    path.join(screenshotsDir, "ocr_report.json"),
    JSON.stringify(results, null, 2),
    "utf-8"
  );
}

main().catch((err) => {
  console.error("OCR evaluation failed:", err);
  process.exit(1);
});
