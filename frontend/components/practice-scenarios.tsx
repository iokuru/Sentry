"use client";

import React from "react";
import { ShieldAlert, Package, Laptop, Landmark, ArrowRight } from "lucide-react";

export interface Scenario {
  id: string;
  category: string;
  icon: React.ComponentType<{ size?: number; className?: string }>;
  title: string;
  message: string;
  tagColor: string;
}

export const sampleScenarios: Scenario[] = [
  {
    id: "bank-fraud",
    category: "Bank Impersonation",
    icon: Landmark,
    title: "Suspicious Debit Alert",
    message: "ALERT: Unauthorized charge of $489.90 at Apple Store. Reply NO or click sentry-bank-verify.com to lock card immediately.",
    tagColor: "#f59e0b",
  },
  {
    id: "package-fee",
    category: "Delivery Fraud",
    icon: Package,
    title: "Unpaid Parcel Fee",
    message: "USPS: Your package has an unpaid handling fee of $1.85. Update your address within 12 hours: bit.ly/usps-pkg-82 to avoid return.",
    tagColor: "#3b82f6",
  },
  {
    id: "tech-support",
    category: "Tech Support Scam",
    icon: Laptop,
    title: "Urgent Virus Warning",
    message: "MICROSOFT DEFENDER: Trojan horse virus detected on your computer. Call toll-free 1-800-555-0199 immediately. Do not restart.",
    tagColor: "#ef4444",
  },
  {
    id: "tax-refund",
    category: "Government Agency",
    icon: ShieldAlert,
    title: "Unclaimed Tax Refund",
    message: "IRS Notification: You have an unclaimed tax rebate of $1,240. Enter your direct deposit information before midnight to claim.",
    tagColor: "#10b981",
  },
];

interface PracticeScenariosProps {
  onSelectScenario: (message: string) => void;
}

export function PracticeScenarios({ onSelectScenario }: PracticeScenariosProps) {
  return (
    <div className="scenarios-wrap">
      <div className="scenarios-header">
        <span className="scenarios-subtitle">Want to test it out first?</span>
        <h3 className="scenarios-title">Try one of these real-world examples</h3>
      </div>

      <div className="scenarios-grid">
        {sampleScenarios.map((s) => {
          const Icon = s.icon;
          return (
            <button
              key={s.id}
              type="button"
              className="scenario-card"
              onClick={() => onSelectScenario(s.message)}
            >
              <div className="scenario-top">
                <span className="scenario-icon" style={{ color: s.tagColor }}>
                  <Icon size={18} />
                </span>
                <span className="scenario-tag" style={{ borderColor: `${s.tagColor}40`, color: s.tagColor }}>
                  {s.category}
                </span>
              </div>

              <strong className="scenario-name">{s.title}</strong>
              <p className="scenario-preview">“{s.message}”</p>

              <div className="scenario-action">
                <span>Check this example</span>
                <ArrowRight size={14} />
              </div>
            </button>
          );
        })}
      </div>

      <style jsx>{`
        .scenarios-wrap {
          margin-top: 36px;
          padding-top: 28px;
          border-top: 1px solid var(--border-light);
        }

        .scenarios-header {
          margin-bottom: 20px;
          text-align: left;
        }

        .scenarios-subtitle {
          display: block;
          font-size: 0.82rem;
          font-weight: 700;
          letter-spacing: 0.08em;
          text-transform: uppercase;
          color: var(--text-muted);
          margin-bottom: 4px;
        }

        .scenarios-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--text-main);
          letter-spacing: -0.02em;
        }

        .scenarios-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
          gap: 16px;
        }

        .scenario-card {
          display: flex;
          flex-direction: column;
          align-items: flex-start;
          text-align: left;
          padding: 18px;
          border-radius: var(--radius-md);
          background: var(--bg-card);
          border: 1px solid var(--border-light);
          cursor: pointer;
          transition: transform 160ms ease, border-color 160ms ease, box-shadow 160ms ease;
          box-shadow: var(--shadow-sm);
        }

        .scenario-card:hover {
          transform: translateY(-2px);
          border-color: var(--color-primary);
          box-shadow: var(--shadow-md);
        }

        .scenario-top {
          display: flex;
          align-items: center;
          gap: 8px;
          width: 100%;
          margin-bottom: 12px;
        }

        .scenario-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
        }

        .scenario-tag {
          font-size: 0.74rem;
          font-weight: 750;
          text-transform: uppercase;
          letter-spacing: 0.06em;
          padding: 2px 7px;
          border-radius: var(--radius-full);
          border: 1px solid;
          background: var(--surface-raised);
        }

        .scenario-name {
          font-size: 0.98rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 6px;
        }

        .scenario-preview {
          font-size: 0.86rem;
          color: var(--text-muted);
          line-height: 1.45;
          margin: 0 0 16px 0;
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .scenario-action {
          margin-top: auto;
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 0.84rem;
          font-weight: 700;
          color: var(--color-primary);
        }

        @media (max-width: 600px) {
          .scenarios-grid {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </div>
  );
}
