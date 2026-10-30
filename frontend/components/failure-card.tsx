import React from "react";
import { AlertCircle, RotateCcw } from "lucide-react";
import { Button } from "@/components/button";

export interface FailureCardProps {
  errorMessage?: string;
  onRetry: () => void;
}

export function FailureCard({ errorMessage, onRetry }: FailureCardProps) {
  return (
    <div className="failure-card" role="alert">
      <div className="failure-header">
        <span className="failure-icon" aria-hidden="true">
          <AlertCircle size={28} strokeWidth={2.4} />
        </span>
        <div>
          <h2 className="failure-title">Could not complete check</h2>
          <p className="failure-desc">
            {errorMessage || "We ran into an issue checking this screenshot. Please try again."}
          </p>
        </div>
      </div>

      <div className="failure-actions">
        <Button variant="outline" onClick={onRetry}>
          <RotateCcw size={18} aria-hidden="true" />
          <span>Try again</span>
        </Button>
      </div>

      <style jsx>{`
        .failure-card {
          background-color: var(--color-danger-subtle);
          border: 2px solid var(--color-danger);
          border-radius: var(--radius-lg);
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .failure-header {
          display: flex;
          align-items: flex-start;
          gap: 14px;
        }

        .failure-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 44px;
          height: 44px;
          border-radius: var(--radius-md);
          background-color: var(--danger-bg);
          color: var(--color-danger);
          flex-shrink: 0;
        }

        .failure-title {
          font-size: 1.3rem;
          font-weight: 700;
          color: var(--danger-title);
          margin-bottom: 4px;
        }

        .failure-desc {
          font-size: 1.05rem;
          color: var(--text-main);
          line-height: 1.5;
        }

        .failure-actions {
          display: flex;
          gap: 12px;
        }
      `}</style>
    </div>
  );
}
