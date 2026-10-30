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

  const handleFile = useCallback(
    (file: File) => {
      if (disabled) return;
      onFileSelect(file);
    },
    [disabled, onFileSelect]
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFile(e.target.files[0]);
    }
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

      const items = e.clipboardData.items;
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.indexOf("image") !== -1) {
          const file = item.getAsFile();
          if (file) {
            handleFile(file);
            return;
          }
        }
      }

      if (onTextPaste) {
        const text = e.clipboardData.getData("text");
        if (text && text.trim().length > 0) {
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
      onClick={handleChooseFile}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      role="button"
      tabIndex={disabled ? -1 : 0}
      aria-label="Upload screenshot or photo to check"
      onKeyDown={(e) => {
        if ((e.key === "Enter" || e.key === " ") && !disabled) {
          e.preventDefault();
          handleChooseFile();
        }
      }}
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
        <div className="dropzone-icon-ring" aria-hidden="true">
          <UploadCloud size={38} strokeWidth={2.2} />
        </div>

        <div className="dropzone-text">
          <p className="dropzone-primary-text">
            {isDragging ? "Drop your image now" : "Tap here to choose a screenshot or photo"}
          </p>
          <p className="dropzone-secondary-text">or paste with Ctrl+V / Cmd+V</p>
        </div>

        <div className="dropzone-actions">
          <button
            type="button"
            className="dropzone-camera-btn"
            onClick={handleCameraCapture}
            disabled={disabled}
            aria-label="Take a photo with camera"
          >
            <Camera size={20} strokeWidth={2.2} />
            <span>Take photo</span>
          </button>
        </div>
      </div>

      <style jsx>{`
        .dropzone-box {
          border: 2.5px dashed var(--border-light);
          border-radius: var(--radius-lg);
          background-color: var(--bg-card);
          padding: 44px 24px;
          text-align: center;
          cursor: pointer;
          transition: border-color 0.2s ease, background-color 0.2s ease;
          user-select: none;
        }

        .dropzone-box:hover:not(.dropzone-disabled),
        .dropzone-box:focus-visible {
          border-color: var(--color-primary);
          background-color: #f7faf8;
        }

        .dropzone-active {
          border-color: var(--color-primary);
          background-color: var(--bg-secondary);
        }

        .dropzone-disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .dropzone-content {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 16px;
        }

        .dropzone-icon-ring {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 72px;
          height: 72px;
          border-radius: var(--radius-full);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
        }

        .dropzone-primary-text {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 4px;
        }

        .dropzone-secondary-text {
          font-size: 1rem;
          color: var(--text-muted);
        }

        .dropzone-actions {
          margin-top: 8px;
        }

        .dropzone-camera-btn {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          padding: 10px 20px;
          border-radius: var(--radius-md);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          font-weight: 600;
          font-size: 1rem;
          border: 1px solid var(--border-light);
          cursor: pointer;
          min-height: 48px;
        }

        .dropzone-camera-btn:hover:not(:disabled) {
          background-color: #dbebe0;
        }
      `}</style>
    </div>
  );
}
