import React, { useRef, useState, useEffect, useCallback } from "react";
import { UploadCloud, Camera } from "lucide-react";

export interface DropzoneProps {
  onFileSelect: (file: File) => void;
  onTextPaste?: (text: string) => void;
  disabled?: boolean;
}

export function Dropzone({ onFileSelect, onTextPaste, disabled = false }: DropzoneProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const cameraInputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [fileError, setFileError] = useState<string | null>(null);

  const handleFile = useCallback(
    (file: File) => {
      if (disabled) return;
      const extension = file.name.toLowerCase().split(".").pop();
      const allowed = file.type.startsWith("image/") || extension === "heic" || extension === "heif";
      if (!allowed) {
        setFileError("Please choose an image, such as a screenshot or photo.");
        return;
      }
      if (file.size > 20 * 1024 * 1024) {
        setFileError("This image is over 20 MB. Try a screenshot or a smaller copy.");
        return;
      }
      setFileError(null);
      onFileSelect(file);
    },
    [disabled, onFileSelect]
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFile(e.target.files[0]);
    }
    e.target.value = "";
  };

  const handleChooseFile = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleCameraCapture = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!disabled && cameraInputRef.current) {
      cameraInputRef.current.click();
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    if (!disabled) {
      setIsDragging(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (!disabled && e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFile(e.dataTransfer.files[0]);
    }
  };

  useEffect(() => {
    const handleGlobalPaste = (e: ClipboardEvent) => {
      if (disabled) return;
      if (!e.clipboardData) return;
      const target = e.target;
      if (
        target instanceof HTMLElement &&
        (target.isContentEditable || target.closest("input, textarea, select, [contenteditable='true']"))
      ) {
        return;
      }

      const items = e.clipboardData.items;
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.indexOf("image") !== -1) {
          const file = item.getAsFile();
          if (file) {
            e.preventDefault();
            handleFile(file);
            return;
          }
        }
      }

      if (onTextPaste) {
        const text = e.clipboardData.getData("text");
        if (text && text.trim().length > 0) {
          e.preventDefault();
          onTextPaste(text);
        }
      }
    };

    window.addEventListener("paste", handleGlobalPaste);
    return () => {
      window.removeEventListener("paste", handleGlobalPaste);
    };
  }, [disabled, handleFile, onTextPaste]);

  return (
    <div
      className={`dropzone-box ${isDragging ? "dropzone-active" : ""} ${disabled ? "dropzone-disabled" : ""}`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      role="group"
      aria-labelledby="dropzone-title"
    >
      <input
        ref={fileInputRef}
        type="file"
        accept="image/png,image/jpeg,image/webp,image/heic,image/heif"
        onChange={handleInputChange}
        className="sr-only"
        aria-hidden="true"
        tabIndex={-1}
      />

      <input
        ref={cameraInputRef}
        type="file"
        accept="image/*"
        capture="environment"
        onChange={handleInputChange}
        className="sr-only"
        aria-hidden="true"
        tabIndex={-1}
      />

      <div className="dropzone-content">
        <div className="dropzone-icon" aria-hidden="true">
          <UploadCloud size={36} strokeWidth={2.2} />
        </div>

        <div className="dropzone-text">
          <p id="dropzone-title" className="dropzone-primary-text">
            {isDragging ? "Drop your screenshot here" : "Drop a screenshot or take a photo"}
          </p>
          <p className="dropzone-secondary-text">PNG, JPG, WEBP, or HEIC up to 20 MB</p>
        </div>

        <div className="dropzone-actions">
          <button
            type="button"
            className="dropzone-select-btn"
            onClick={handleChooseFile}
            disabled={disabled}
          >
            Choose screenshot
          </button>
          <button
            type="button"
            className="dropzone-camera-btn"
            onClick={handleCameraCapture}
            disabled={disabled}
            aria-label="Take a photo with camera"
          >
            <Camera size={18} strokeWidth={2.2} />
            <span>Take photo</span>
          </button>
        </div>

        <p className="dropzone-paste">or paste directly with Ctrl + V / ⌘ + V</p>
        {fileError && <p className="dropzone-error" role="alert">{fileError}</p>}
      </div>

      <style jsx>{`
        .dropzone-box {
          border: 1.5px dashed var(--border-medium);
          border-radius: var(--radius-lg);
          background-color: var(--surface-translucent);
          backdrop-filter: blur(8px);
          padding: 32px 20px 28px;
          text-align: center;
          cursor: default;
          transition: border-color 180ms ease, background-color 180ms ease, box-shadow 180ms ease;
          user-select: none;
        }

        .dropzone-box:hover:not(.dropzone-disabled),
        .dropzone-box:focus-within {
          border-color: var(--color-primary);
          background-color: var(--surface-raised);
          box-shadow: 0 0 20px rgba(16, 185, 129, 0.08);
        }

        .dropzone-active {
          border-color: var(--color-primary);
          background-color: var(--bg-secondary);
          box-shadow: 0 0 24px rgba(16, 185, 129, 0.15);
        }

        .dropzone-disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .dropzone-content {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 12px;
        }

        .dropzone-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 48px;
          height: 48px;
          border-radius: var(--radius-full);
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          color: var(--color-primary);
        }

        .dropzone-primary-text {
          font-size: 1.12rem;
          font-weight: 700;
          color: var(--text-main);
          letter-spacing: -0.015em;
          margin-bottom: 4px;
        }

        .dropzone-secondary-text {
          font-size: 0.88rem;
          color: var(--text-muted);
        }

        .dropzone-actions {
          margin-top: 6px;
          display: flex;
          justify-content: center;
          flex-wrap: wrap;
          gap: 10px;
        }

        .dropzone-select-btn {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 9px 18px;
          min-height: 42px;
          border: 1px solid var(--color-primary);
          border-radius: var(--radius-full);
          background: var(--color-primary);
          color: #ffffff;
          font-size: 0.94rem;
          font-weight: 650;
          cursor: pointer;
          transition: background 140ms ease, transform 140ms ease, box-shadow 140ms ease;
          box-shadow: 0 2px 8px rgba(16, 185, 129, 0.25);
        }

        .dropzone-select-btn:hover:not(:disabled) {
          background: var(--color-primary-hover);
          transform: translateY(-1px);
        }

        .dropzone-camera-btn {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          padding: 9px 18px;
          border-radius: var(--radius-full);
          background-color: var(--surface-raised);
          color: var(--text-main);
          font-weight: 600;
          font-size: 0.94rem;
          border: 1px solid var(--border-light);
          cursor: pointer;
          min-height: 42px;
          transition: background-color 140ms ease, border-color 140ms ease;
        }

        .dropzone-camera-btn:hover:not(:disabled) {
          background-color: var(--bg-muted);
          border-color: var(--text-muted);
        }

        .dropzone-paste {
          color: var(--text-muted);
          font-size: 0.86rem;
          margin-top: 2px;
        }

        .dropzone-error {
          color: var(--color-danger);
          font-size: 0.88rem;
          font-weight: 600;
        }
      `}</style>
    </div>
  );
}
