import React, { useState } from "react";
import { Chip } from "@/components/chip";
import { Clock, User } from "lucide-react";

export interface IncidentItem {
  incident_id: string;
  senior_id?: string;
  senior_name?: string;
  status: "NEW" | "ACKNOWLEDGED" | "RESOLVED";
  severity: "HIGH" | "MEDIUM" | "LOW";
  signals: string[];
  hold_off: string;
  safe_alternative: string;
  created_at: string;
  thumbnail_url?: string;
}

export interface IncidentCardProps {
  incident: IncidentItem;
  onUpdateStatus: (id: string, newStatus: "ACKNOWLEDGED" | "RESOLVED") => Promise<void>;
  showThumbnails?: boolean;
}

export function IncidentCard({
  incident,
  onUpdateStatus,
  showThumbnails = true,
}: IncidentCardProps) {
  const [isUpdating, setIsUpdating] = useState(false);

  const formatTimeAgo = (dateStr: string) => {
    try {
      const now = Date.now();
      const past = new Date(dateStr).getTime();
      const diffMinutes = Math.floor((now - past) / (1000 * 60));

      if (diffMinutes < 1) return "Just now";
      if (diffMinutes < 60) return `${diffMinutes}m ago`;
      const diffHours = Math.floor(diffMinutes / 60);
      if (diffHours < 24) return `${diffHours}h ago`;
      const diffDays = Math.floor(diffHours / 24);
      return `${diffDays}d ago`;
    } catch {
      return "Recently";
    }
  };

  const handleStatus = async (status: "ACKNOWLEDGED" | "RESOLVED") => {
    setIsUpdating(true);
    try {
      await onUpdateStatus(incident.incident_id, status);
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <div className={`incident-card-root incident-${incident.status.toLowerCase()}`}>
      <div className="incident-header">
        <div className="incident-pills">
          <Chip
            variant={
              incident.severity === "HIGH"
                ? "high"
                : incident.severity === "LOW"
                ? "low"
                : "medium"
            }
          >
            {incident.severity} RISK
          </Chip>
          <span className={`status-pill status-${incident.status.toLowerCase()}`}>
            {incident.status}
          </span>
        </div>

        <div className="incident-time">
          <Clock size={15} aria-hidden="true" />
          <span>{formatTimeAgo(incident.created_at)}</span>
        </div>
      </div>

      <div className="incident-body">
        {showThumbnails && incident.thumbnail_url ? (
          <div className="incident-thumbnail-box">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={incident.thumbnail_url}
              alt="Screenshot thumbnail"
              className="incident-thumb-img"
            />
          </div>
        ) : null}

        <div className="incident-info">
          {incident.senior_name ? (
            <div className="incident-senior-tag">
              <User size={14} aria-hidden="true" />
              <span>{incident.senior_name}</span>
            </div>
          ) : null}

          <p className="incident-hold-off">{incident.hold_off}</p>

          {incident.signals && incident.signals.length > 0 ? (
            <div className="incident-signals">
              {incident.signals.map((sig) => (
                <span key={sig} className="signal-tag">
                  #{sig}
                </span>
              ))}
            </div>
          ) : null}
        </div>
      </div>

      <div className="incident-actions">
        {incident.status === "NEW" ? (
          <button
            type="button"
            className="action-btn ack-btn"
            disabled={isUpdating}
            onClick={() => handleStatus("ACKNOWLEDGED")}
          >
            Acknowledge
          </button>
        ) : null}

        {incident.status !== "RESOLVED" ? (
          <button
            type="button"
            className="action-btn resolve-btn"
            disabled={isUpdating}
            onClick={() => handleStatus("RESOLVED")}
          >
            Mark resolved
          </button>
        ) : (
          <span className="resolved-notice">Resolved</span>
        )}
      </div>

      <style jsx>{`
        .incident-card-root {
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 22px;
          display: flex;
          flex-direction: column;
          gap: 16px;
          box-shadow: var(--shadow-sm);
          transition: border-color 160ms ease;
        }

        .incident-card-root:hover {
          border-color: var(--border-medium);
        }

        .incident-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 12px;
        }

        .incident-pills {
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .status-pill {
          font-size: 0.78rem;
          font-weight: 750;
          letter-spacing: 0.04em;
          padding: 3px 10px;
          border-radius: var(--radius-full);
          text-transform: uppercase;
        }

        .status-new {
          background-color: var(--danger-bg);
          color: var(--danger-title);
          border: 1px solid rgba(239, 68, 68, 0.25);
        }

        .status-acknowledged {
          background-color: var(--warning-bg);
          color: var(--warning-title);
          border: 1px solid rgba(245, 158, 11, 0.25);
        }

        .status-resolved {
          background-color: var(--safe-bg);
          color: var(--safe-title);
          border: 1px solid rgba(16, 185, 129, 0.25);
        }

        .incident-time {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 0.86rem;
          color: var(--text-muted);
        }

        .incident-body {
          display: flex;
          gap: 16px;
        }

        .incident-thumbnail-box {
          width: 80px;
          height: 80px;
          border-radius: var(--radius-md);
          overflow: hidden;
          background-color: var(--bg-page);
          border: 1px solid var(--border-light);
          flex-shrink: 0;
        }

        .incident-thumb-img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        .incident-info {
          display: flex;
          flex-direction: column;
          gap: 8px;
          flex: 1;
        }

        .incident-senior-tag {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 0.86rem;
          font-weight: 650;
          color: var(--color-primary);
        }

        .incident-hold-off {
          font-size: 1.05rem;
          font-weight: 650;
          color: var(--text-main);
          line-height: 1.5;
        }

        .incident-signals {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }

        .signal-tag {
          font-size: 0.78rem;
          font-weight: 600;
          color: var(--text-muted);
          background-color: var(--surface-raised);
          border: 1px solid var(--border-light);
          padding: 2px 8px;
          border-radius: var(--radius-full);
        }

        .incident-actions {
          display: flex;
          align-items: center;
          gap: 10px;
          border-top: 1px solid var(--border-light);
          padding-top: 14px;
        }

        .action-btn {
          padding: 8px 16px;
          border-radius: var(--radius-full);
          font-size: 0.88rem;
          font-weight: 650;
          cursor: pointer;
          border: 1px solid transparent;
          min-height: 38px;
          transition: background-color 140ms ease;
        }

        .ack-btn {
          background-color: var(--warning-bg);
          color: var(--warning-title);
          border-color: var(--warning-border);
        }

        .ack-btn:hover:not(:disabled) {
          background-color: var(--surface-raised);
        }

        .resolve-btn {
          background-color: var(--surface-raised);
          color: var(--text-main);
          border-color: var(--border-light);
        }

        .resolve-btn:hover:not(:disabled) {
          background-color: var(--bg-muted);
          border-color: var(--text-muted);
        }

        .resolved-notice {
          font-size: 0.9rem;
          font-weight: 700;
          color: var(--color-primary);
        }
      `}</style>
    </div>
  );
}
