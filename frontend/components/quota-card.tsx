import React from "react";
import { Clock } from "lucide-react";
import { HotlineButton } from "@/components/hotline-button";

export interface QuotaCardProps {
  country: string;
}

export function QuotaCard({ country }: QuotaCardProps) {
  return (
    <div className="quota-card" role="status">
      <div className="quota-header">
        <span className="quota-icon" aria-hidden="true">
          <Clock size={28} strokeWidth={2.4} />
        </span>
        <div>
          <h2 className="quota-title">Daily check limit reached</h2>
          <p className="quota-desc">
            To ensure Sentry remains completely free for everyone, checks are limited each day. Daily checks will reset tomorrow morning.
          </p>
        </div>
      </div>

      <div className="quota-reassurance">
        <p className="quota-reassurance-text">
          If you are worried about a message or phone call right now, official fraud specialists are available by phone.
        </p>

        <HotlineButton country={country} />
      </div>

      <style jsx>{`
        .quota-card {
          background-color: var(--warning-bg);
          border: 1px solid var(--warning-border);
          border-radius: var(--radius-lg);
          padding: 28px 24px;
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .quota-header {
          display: flex;
          align-items: flex-start;
          gap: 16px;
        }

        .quota-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 48px;
          height: 48px;
          border-radius: var(--radius-md);
          background-color: var(--warning-icon-bg);
          color: var(--warning-icon-ink);
          flex-shrink: 0;
        }

        .quota-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: var(--warning-title);
          margin-bottom: 6px;
        }

        .quota-desc {
          font-size: 1.05rem;
          color: var(--text-main);
          line-height: 1.5;
        }

        .quota-reassurance {
          background-color: var(--surface-raised);
          border-radius: var(--radius-md);
          padding: 16px;
          border: 1px solid var(--warning-border);
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        .quota-reassurance-text {
          font-size: 1rem;
          color: var(--text-muted);
        }

      `}</style>
    </div>
  );
}
