import React from "react";
import { AlertTriangle, CreditCard, ShieldAlert, Monitor, PhoneCall, ExternalLink } from "lucide-react";

export interface SignalChipProps {
  signal: string;
}

const signalLabels: Record<string, { label: string; icon: React.ReactNode }> = {
  payment_pressure: { label: "Payment Demand", icon: <CreditCard size={14} /> },
  impersonation: { label: "Impersonation", icon: <ShieldAlert size={14} /> },
  remote_access: { label: "Remote Access", icon: <Monitor size={14} /> },
  urgency: { label: "High Urgency", icon: <AlertTriangle size={14} /> },
  fake_link: { label: "Suspicious Link", icon: <ExternalLink size={14} /> },
  phone_scam: { label: "Unsolicited Call", icon: <PhoneCall size={14} /> },
};

export function SignalChip({ signal }: SignalChipProps) {
  const norm = signal.toLowerCase().replace(/-/g, "_");
  const config = signalLabels[norm] || {
    label: signal.replace(/_/g, " "),
    icon: <AlertTriangle size={14} />,
  };

  return (
    <span className="signal-chip">
      <span className="signal-chip-icon" aria-hidden="true">
        {config.icon}
      </span>
      <span className="signal-chip-text">{config.label}</span>

      <style jsx>{`
        .signal-chip {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          padding: 4px 10px;
          border-radius: var(--radius-sm);
          font-size: 0.82rem;
          font-weight: 600;
          line-height: 1.3;
        }

        .signal-chip-icon {
          display: inline-flex;
          align-items: center;
        }

        .signal-chip-text {
          text-transform: capitalize;
        }
      `}</style>
    </span>
  );
}
