import React from "react";
import { Loader2 } from "lucide-react";

export interface ProgressCardProps {
  phase: "uploading" | "analyzing";
}

export function ProgressCard({ phase }: ProgressCardProps) {
  const isUploading = phase === "uploading";

  return (
    <div
      className="progress-card"
      role="status"
      aria-live="polite"
      aria-label={isUploading ? "Uploading image" : "Analyzing message"}
    >
      <div className="progress-spinner" aria-hidden="true">
        <Loader2 size={40} className="animate-spin" />
      </div>

      <div className="progress-text">
        <h2 className="progress-title">
          {isUploading ? "Uploading screenshot securely" : "Checking message for scam signals"}
        </h2>
        <p className="progress-sub">
          {isUploading
            ? "Preparing your photo for a quick check. This only takes a moment."
            : "Reviewing text, urgency patterns, and payment requests to provide safe guidance."}
        </p>
      </div>

      <style jsx>{`
        .progress-card {
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 48px 24px;
          display: flex;
          flex-direction: column;
          align-items: center;
          text-align: center;
          gap: 20px;
          box-shadow: var(--shadow-sm);
          max-width: 580px;
          margin: 30px auto;
        }

        .progress-spinner {
          display: grid;
          place-items: center;
          color: var(--color-primary);
          animation: spin 1s linear infinite;
        }

        .progress-title {
          font-size: 1.4rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 8px;
        }

        .progress-sub {
          font-size: 1.05rem;
          color: var(--text-muted);
          max-width: 480px;
          margin: 0 auto;
          line-height: 1.5;
        }

        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
