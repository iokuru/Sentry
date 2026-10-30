export interface OcrProgress {
  status: string;
  progress: number;
}

export async function extractTextInBrowser(
  imageSource: Blob | File | string,
  onProgress?: (progress: OcrProgress) => void
): Promise<string> {
  // lazy load tesseract only when browser ocr is actually triggered
  const { createWorker } = await import("tesseract.js");

  const worker = await createWorker("eng", 1, {
    workerPath: "/tesseract/worker.min.js",
    corePath: "/tesseract/tesseract-core-lstm.wasm.js",
    langPath: "/tesseract",
    logger: (m) => {
      if (onProgress && m.status) {
        onProgress({
          status: m.status,
          progress: typeof m.progress === "number" ? m.progress : 0,
        });
      }
    },
  });

  try {
    const ret = await worker.recognize(imageSource);
    return ret.data.text || "";
  } finally {
    await worker.terminate();
  }
}
