export async function downscaleImage(
  file: File,
  maxDimension = 1600,
  quality = 0.85
): Promise<{ blob: Blob; previewUrl: string }> {
  // reject unreasonably large input before processing
  if (file.size > 20 * 1024 * 1024) {
    throw new Error('Image exceeds 20MB limit. Please select a smaller screenshot.');
  }

  return new Promise((resolve, reject) => {
    const objectUrl = URL.createObjectURL(file);
    const img = new Image();

    img.onload = () => {
      URL.revokeObjectURL(objectUrl);
      let { width, height } = img;

      if (width > maxDimension || height > maxDimension) {
        if (width > height) {
          height = Math.round((height * maxDimension) / width);
          width = maxDimension;
        } else {
          width = Math.round((width * maxDimension) / height);
          height = maxDimension;
        }
      }

      const canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;

      const ctx = canvas.getContext('2d');
      if (!ctx) {
        resolve({ blob: file, previewUrl: URL.createObjectURL(file) });
        return;
      }

      ctx.drawImage(img, 0, 0, width, height);
      canvas.toBlob(
        (blob) => {
          if (!blob) {
            resolve({ blob: file, previewUrl: URL.createObjectURL(file) });
            return;
          }
          // ensure blob is within 4mb textract limit
          if (blob.size > 4 * 1024 * 1024) {
            canvas.toBlob(
              (smallerBlob) => {
                const finalBlob = smallerBlob || blob;
                resolve({ blob: finalBlob, previewUrl: URL.createObjectURL(finalBlob) });
              },
              'image/jpeg',
              0.65
            );
            return;
          }
          const previewUrl = URL.createObjectURL(blob);
          resolve({ blob, previewUrl });
        },
        'image/jpeg',
        quality
      );
    };

    img.onerror = () => {
      URL.revokeObjectURL(objectUrl);
      reject(new Error('failed to load image for downscaling'));
    };

    img.src = objectUrl;
  });
}
