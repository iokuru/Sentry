import React, { useState } from "react";
import { Chip } from "@/components/chip";
import { Clock, User, Image as ImageIcon } from "lucide-react";

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
          <Clock size={16} aria-hidden="true" />
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
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 20px;
          display: flex;
          flex-direction: column;
          gap: 16px;
          box-shadow: var(--shadow-sm);
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
          font-size: 0.82rem;
          font-weight: 700;
          padding: 3px 10px;
          border-radius: var(--radius-full);
          text-transform: uppercase;
        }

        .status-new {
          background-color: #fee2e2;
          color: #991b1b;
        }

        .status-acknowledged {
          background-color: #fef3c7;
          color: #92400e;
        }

        .status-resolved {
          background-color: #dcfce7;
          color: #166534;
        }

        .incident-time {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 0.9rem;
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
          font-size: 0.88rem;
          font-weight: 600;
          color: var(--color-primary);
        }

        .incident-hold-off {
          font-size: 1.05rem;
          font-weight: 600;
          color: var(--text-main);
          line-height: 1.45;
        }

        .incident-signals {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }

        .signal-tag {
          font-size: 0.82rem;
          font-weight: 600;
          color: var(--text-muted);
          background-color: var(--bg-muted);
          padding: 2px 8px;
          border-radius: var(--radius-sm);
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
          border-radius: var(--radius-sm);
          font-size: 0.95rem;
          font-weight: 600;
          cursor: pointer;
          border: none;
          min-height: 40px;
        }

        .ack-btn {
          background-color: #fef3c7;
          color: #92400e;
        }

        .ack-btn:hover:not(:disabled) {
          background-color: #fde68a;
        }

        .resolve-btn {
          background-color: var(--bg-secondary);
          color: var(--color-primary);
        }

        .resolve-btn:hover:not(:disabled) {
          background-color: #dbebe0;
        }

        .resolved-notice {
          font-size: 0.95rem;
          font-weight: 600;
          color: var(--color-safe);
        }
      `}</style>
    </div>
  );
}
