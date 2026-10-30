import React, { useState } from "react";
import { Check, Clock } from "lucide-react";

export interface IncidentActionsProps {
  currentStatus: "NEW" | "ACKNOWLEDGED" | "RESOLVED";
  onTransition: (newStatus: "ACKNOWLEDGED" | "RESOLVED") => Promise<void>;
}

export function IncidentActions({
  currentStatus,
  onTransition,
}: IncidentActionsProps) {
  const [loadingStatus, setLoadingStatus] = useState<string | null>(null);

  const handleAction = async (status: "ACKNOWLEDGED" | "RESOLVED") => {
    setLoadingStatus(status);
    try {
      await onTransition(status);
    } finally {
      setLoadingStatus(null);
    }
  };

  return (
    <div className="actions-wrapper">
      {currentStatus === "NEW" ? (
        <button
          type="button"
          className="btn-ack"
          disabled={Boolean(loadingStatus)}
          onClick={() => handleAction("ACKNOWLEDGED")}
        >
          <Clock size={16} aria-hidden="true" />
          <span>{loadingStatus === "ACKNOWLEDGED" ? "Updating..." : "Acknowledge"}</span>
        </button>
      ) : null}

      {currentStatus !== "RESOLVED" ? (
        <button
          type="button"
          className="btn-resolve"
          disabled={Boolean(loadingStatus)}
          onClick={() => handleAction("RESOLVED")}
        >
          <Check size={16} aria-hidden="true" />
          <span>{loadingStatus === "RESOLVED" ? "Resolving..." : "Mark resolved"}</span>
        </button>
      ) : (
        <span className="badge-resolved">
          <Check size={14} aria-hidden="true" />
          <span>Resolved</span>
        </span>
      )}

      <style jsx>{`
        .actions-wrapper {
          display: flex;
          align-items: center;
          gap: 10px;
        }

        button {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 16px;
          border-radius: var(--radius-sm);
          font-size: 0.92rem;
          font-weight: 600;
          cursor: pointer;
          border: 1px solid transparent;
          transition: background-color 0.15s ease;
        }

        .btn-ack {
          background-color: #fef3c7;
          color: #92400e;
          border-color: #fde68a;
        }

        .btn-ack:hover:not(:disabled) {
          background-color: #fde68a;
        }

        .btn-resolve {
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          border-color: #cce3d4;
        }

        .btn-resolve:hover:not(:disabled) {
          background-color: #d8eadd;
        }

        button:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .badge-resolved {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 0.92rem;
          font-weight: 700;
          color: var(--color-safe);
        }
      `}</style>
    </div>
  );
}
