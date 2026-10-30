import React from "react";
import { CheckCircle2 } from "lucide-react";

export interface AlternativeCardProps {
  action: string;
}

export function AlternativeCard({ action }: AlternativeCardProps) {
  return (
    <div className="alternative-card">
      <div className="alternative-header">
        <span className="alternative-icon" aria-hidden="true">
          <CheckCircle2 size={24} strokeWidth={2.4} />
        </span>
        <h2 className="alternative-title">2. A safe alternative action</h2>
      </div>

      <div className="alternative-body">
        <p className="alternative-action">{action}</p>
      </div>

      <style jsx>{`
        .alternative-card {
          background-color: var(--safe-bg);
          border: 1px solid var(--safe-border);
          border-radius: var(--radius-lg);
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 12px;
          box-shadow: var(--shadow-sm);
        }

        .alternative-header {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .alternative-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 42px;
          height: 42px;
          border-radius: var(--radius-md);
          background-color: var(--safe-icon-bg);
          color: var(--safe-icon-ink);
          flex-shrink: 0;
        }

        .alternative-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: var(--safe-title);
          letter-spacing: -0.01em;
        }

        .alternative-body {
          padding-left: 54px;
        }

        .alternative-action {
          font-size: 1.15rem;
          line-height: 1.6;
          color: var(--text-main);
          font-weight: 600;
        }

        @media (max-width: 600px) {
          .alternative-body {
            padding-left: 0;
          }
          .alternative-title {
            font-size: 1.2rem;
          }
          .alternative-action {
            font-size: 1.05rem;
          }
        }
      `}</style>
    </div>
  );
}
