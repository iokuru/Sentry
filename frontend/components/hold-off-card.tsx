import React from "react";
import { Hand } from "lucide-react";

export interface HoldOffCardProps {
  advice: string;
}

export function HoldOffCard({ advice }: HoldOffCardProps) {
  return (
    <div className="hold-card">
      <div className="hold-header">
        <span className="hold-icon" aria-hidden="true">
          <Hand size={24} strokeWidth={2.4} />
        </span>
        <h2 className="hold-title">1. Hold off on this right now</h2>
      </div>

      <div className="hold-body">
        <p className="hold-advice">{advice}</p>
      </div>

      <style jsx>{`
        .hold-card {
          background-color: #fffaf0;
          border: 2px solid #f59e0b;
          border-radius: var(--radius-lg);
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 12px;
          box-shadow: var(--shadow-sm);
        }

        .hold-header {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .hold-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 42px;
          height: 42px;
          border-radius: var(--radius-md);
          background-color: #fef3c7;
          color: #b45309;
          flex-shrink: 0;
        }

        .hold-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: #92400e;
          letter-spacing: -0.01em;
        }

        .hold-body {
          padding-left: 54px;
        }

        .hold-advice {
          font-size: 1.15rem;
          line-height: 1.6;
          color: var(--text-main);
          font-weight: 600;
        }

        @media (max-width: 600px) {
          .hold-body {
            padding-left: 0;
          }
          .hold-title {
            font-size: 1.2rem;
          }
          .hold-advice {
            font-size: 1.05rem;
          }
        }
      `}</style>
    </div>
  );
}
