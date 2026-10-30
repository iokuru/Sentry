import React from "react";
import { Clock, PhoneCall } from "lucide-react";

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

        <a href="tel:18773824357" className="quota-call-btn">
          <PhoneCall size={20} strokeWidth={2.4} aria-hidden="true" />
          <span>Call official fraud hotline</span>
        </a>
      </div>

      <style jsx>{`
        .quota-card {
          background-color: #fefce8;
          border: 2px solid #ca8a04;
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
          background-color: #fef08a;
          color: #a16207;
          flex-shrink: 0;
        }

        .quota-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: #854d0e;
          margin-bottom: 6px;
        }

        .quota-desc {
          font-size: 1.05rem;
          color: var(--text-main);
          line-height: 1.5;
        }

        .quota-reassurance {
          background-color: #ffffff;
          border-radius: var(--radius-md);
          padding: 16px;
          border: 1px solid #fde047;
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        .quota-reassurance-text {
          font-size: 1rem;
          color: var(--text-muted);
        }

        .quota-call-btn {
          display: inline-flex;
          align-items: center;
          gap: 10px;
          padding: 12px 20px;
          background-color: var(--color-primary);
          color: var(--text-inverse);
          border-radius: var(--radius-md);
          font-weight: 700;
          text-decoration: none;
          width: fit-content;
        }

        .quota-call-btn:hover {
          background-color: var(--color-primary-hover);
        }
      `}</style>
    </div>
  );
}
