import React from "react";
import { Check, Clock } from "lucide-react";

export function ReadyBadge() {
  return (
    <div className="ready-pill" role="status" aria-label="System status: ready">
      <span className="ready-dot" aria-hidden="true">
        <Check size={10} strokeWidth={3.5} />
      </span>
      <span className="ready-text">Ready</span>

      <style jsx>{`
        .ready-pill {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          padding: 4px 10px;
          border-radius: var(--radius-full);
          box-shadow: var(--shadow-sm);
        }

        .ready-dot {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 16px;
          height: 16px;
          border-radius: var(--radius-full);
          background-color: #10b981;
          color: #ffffff;
          box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
        }

        .ready-text {
          font-size: 0.84rem;
          font-weight: 700;
          color: var(--text-main);
          letter-spacing: 0.01em;
        }
      `}</style>
    </div>
  );
}

export function BeforeYouReplyBadge() {
  return (
    <div className="before-reply-badge" aria-hidden="true">
      <span className="before-reply-icon">
        <Clock size={13} strokeWidth={2.6} />
      </span>
      <span className="before-reply-text">BEFORE YOU REPLY</span>

      <style jsx>{`
        .before-reply-badge {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          background-color: var(--color-accent-subtle);
          border: 1px solid var(--color-accent-border);
          padding: 6px 14px;
          border-radius: var(--radius-full);
          width: fit-content;
        }

        .before-reply-icon {
          display: inline-flex;
          align-items: center;
          color: var(--color-accent-gold);
        }

        .before-reply-text {
          font-size: 0.76rem;
          font-weight: 800;
          letter-spacing: 0.14em;
          color: var(--color-accent-gold);
          text-transform: uppercase;
        }
      `}</style>
    </div>
  );
}

export function ReassuranceStrip() {
  return (
    <div className="reassurance-strip" aria-label="Privacy and safety assurances">
      <div className="reassurance-item">
        <span className="glowing-dot" aria-hidden="true" />
        <span>Free to use</span>
      </div>
      <div className="reassurance-divider" aria-hidden="true">·</div>
      <div className="reassurance-item">
        <span className="glowing-dot" aria-hidden="true" />
        <span>No account needed</span>
      </div>
      <div className="reassurance-divider" aria-hidden="true">·</div>
      <div className="reassurance-item">
        <span className="glowing-dot" aria-hidden="true" />
        <span>Your screenshot is deleted after its check</span>
      </div>

      <style jsx>{`
        .reassurance-strip {
          display: flex;
          align-items: center;
          justify-content: center;
          flex-wrap: wrap;
          gap: 12px;
          padding: 8px 0;
          color: var(--text-muted);
          font-size: 0.94rem;
          font-weight: 500;
        }

        .reassurance-item {
          display: inline-flex;
          align-items: center;
          gap: 8px;
        }

        .glowing-dot {
          display: inline-block;
          width: 7px;
          height: 7px;
          border-radius: 50%;
          background-color: #10b981;
          box-shadow: 0 0 6px rgba(16, 185, 129, 0.8);
          flex-shrink: 0;
        }

        .reassurance-divider {
          color: var(--text-subtle);
          font-weight: 700;
        }

        @media (max-width: 600px) {
          .reassurance-strip {
            flex-direction: column;
            gap: 8px;
          }
          .reassurance-divider {
            display: none;
          }
        }
      `}</style>
    </div>
  );
}
