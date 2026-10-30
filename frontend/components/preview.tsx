import React from "react";
import { Image as ImageIcon, X } from "lucide-react";

export interface PreviewCardProps {
  file: File;
  previewUrl: string;
  onClear: () => void;
  disabled?: boolean;
}

export function PreviewCard({
  file,
  previewUrl,
  onClear,
  disabled = false,
}: PreviewCardProps) {
  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  return (
    <div className="preview-root">
      <div className="preview-image-box">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={previewUrl}
          alt="Selected screenshot for review"
          className="preview-img"
        />
      </div>

      <div className="preview-meta">
        <div className="preview-info">
          <span className="preview-icon" aria-hidden="true">
            <ImageIcon size={20} />
          </span>
          <div className="preview-details">
            <p className="preview-filename">{file.name}</p>
            <p className="preview-filesize">{formatSize(file.size)}</p>
          </div>
        </div>

        <button
          type="button"
          onClick={onClear}
          disabled={disabled}
          className="preview-clear-btn"
          aria-label="Remove selected image"
        >
          <X size={18} strokeWidth={2.4} />
          <span>Remove</span>
        </button>
      </div>

      <style jsx>{`
        .preview-root {
          background-color: var(--bg-card);
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 16px;
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .preview-image-box {
          width: 100%;
          max-height: 380px;
          border-radius: var(--radius-md);
          overflow: hidden;
          background-color: var(--bg-page);
          display: flex;
          align-items: center;
          justify-content: center;
          border: 1px solid var(--border-light);
        }

        .preview-img {
          max-width: 100%;
          max-height: 380px;
          object-fit: contain;
          display: block;
        }

        .preview-meta {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 4px;
        }

        .preview-info {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .preview-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 36px;
          height: 36px;
          border-radius: var(--radius-sm);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
        }

        .preview-filename {
          font-size: 1rem;
          font-weight: 600;
          color: var(--text-main);
          max-width: 300px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .preview-filesize {
          font-size: 0.88rem;
          color: var(--text-muted);
        }

        .preview-clear-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 14px;
          border-radius: var(--radius-sm);
          background-color: var(--bg-muted);
          color: var(--color-danger);
          font-size: 0.95rem;
          font-weight: 600;
          border: none;
          cursor: pointer;
        }

        .preview-clear-btn:hover:not(:disabled) {
          background-color: #fee2e2;
        }
      `}</style>
    </div>
  );
}
