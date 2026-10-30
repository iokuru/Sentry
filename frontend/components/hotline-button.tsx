import React from "react";
import { PhoneCall } from "lucide-react";

export interface HotlineButtonProps {
  country: string;
  hotlineName?: string;
  phoneNumber?: string;
  officialUrl?: string;
}

const defaultHotlines: Record<string, { name: string; phone: string; numberLabel: string }> = {
  US: { name: "US FTC Fraud Hotline", phone: "+18773824357", numberLabel: "1-877-382-4357" },
  UK: { name: "UK Action Fraud", phone: "+443001232040", numberLabel: "0300 123 2040" },
  IN: { name: "National Cyber Crime Helpline", phone: "1930", numberLabel: "1930" },
  AU: { name: "Australian Scamwatch Hotline", phone: "+611300795995", numberLabel: "1300 795 995" },
};

export function HotlineButton({
  country,
  hotlineName,
  phoneNumber,
}: HotlineButtonProps) {
  const normCountry = (country || "US").toUpperCase();
  const fallback = defaultHotlines[normCountry] || defaultHotlines.US;

  const displayName = hotlineName || fallback.name;
  const telNumber = phoneNumber || fallback.phone;
  const visibleNumber = fallback.numberLabel;

  return (
    <div className="hotline-card">
      <div className="hotline-header">
        <h2 className="hotline-title">3. Official fraud hotline</h2>
        <p className="hotline-desc">
          Tap the button below to speak directly with official fraud specialists in your country.
        </p>
      </div>

      <a
        href={`tel:${telNumber}`}
        className="hotline-dial-btn"
        aria-label={`Call ${displayName} at ${visibleNumber}`}
      >
        <span className="hotline-dial-icon" aria-hidden="true">
          <PhoneCall size={26} strokeWidth={2.4} />
        </span>
        <div className="hotline-dial-text">
          <span className="hotline-agency">{displayName}</span>
          <span className="hotline-phone-number">Dial {visibleNumber}</span>
        </div>
      </a>

      <style jsx>{`
        .hotline-card {
          background-color: var(--bg-card);
          border: 2px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
          box-shadow: var(--shadow-sm);
        }

        .hotline-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: var(--text-main);
          letter-spacing: -0.01em;
          margin-bottom: 6px;
        }

        .hotline-desc {
          font-size: 1.05rem;
          color: var(--text-muted);
          line-height: 1.5;
        }

        .hotline-dial-btn {
          display: flex;
          align-items: center;
          gap: 16px;
          background-color: var(--color-primary);
          color: var(--text-inverse);
          padding: 16px 24px;
          border-radius: var(--radius-md);
          text-decoration: none;
          min-height: 64px;
          transition: background-color 0.15s ease, transform 0.05s ease;
        }

        .hotline-dial-btn:hover {
          background-color: var(--color-primary-hover);
        }

        .hotline-dial-btn:active {
          transform: translateY(1px);
        }

        .hotline-dial-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 44px;
          height: 44px;
          border-radius: var(--radius-full);
          background-color: rgba(255, 255, 255, 0.2);
          flex-shrink: 0;
        }

        .hotline-dial-text {
          display: flex;
          flex-direction: column;
        }

        .hotline-agency {
          font-size: 1.15rem;
          font-weight: 700;
          line-height: 1.3;
        }

        .hotline-phone-number {
          font-size: 1rem;
          font-weight: 500;
          opacity: 0.92;
        }
      `}</style>
    </div>
  );
}
