export type OCRProgressCallback = (message: string) => void;

let workerInstance: any = null;

export async function extractTextWithBrowserOCR(
  imageSource: Blob | File | string,
  onProgress?: OCRProgressCallback
): Promise<string> {
  if (onProgress) {
    onProgress('Reading screenshot on your device...');
  }

  // dynamically import tesseract.js so it is not bundled on initial page load
  const { createWorker } = await import('tesseract.js');

  if (!workerInstance) {
    workerInstance = await createWorker('eng', 1, {
      workerPath: '/tesseract/worker.min.js',
      corePath: '/tesseract',
      langPath: '/tesseract',
      gzip: true,
      logger: (m: { status?: string; progress?: number }) => {
        if (!onProgress || !m || !m.status) return;
        if (m.status === 'recognizing text') {
          onProgress('Reading screenshot on your device...');
        } else if (m.status === 'loading tesseract core') {
          onProgress('Preparing reader on your device...');
        } else if (m.status === 'loading language traineddata') {
          onProgress('Preparing language model...');
        }
      }
    });
  }

  try {
    const res = await workerInstance.recognize(imageSource);
    return res && res.data && res.data.text ? res.data.text.trim() : '';
  } catch {
    return '';
  }
}

export async function terminateOCRWorker(): Promise<void> {
  if (workerInstance) {
    try {
      await workerInstance.terminate();
    } catch {
      // ignore
    }
    workerInstance = null;
  }
}
