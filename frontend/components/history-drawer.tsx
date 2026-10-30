import React, { useState, useEffect } from "react";
import { History, X, Trash2 } from "lucide-react";
import { getCheckHistory, clearCheckHistory, type HistoryEntry } from "@/lib/storage";
import { Chip } from "@/components/chip";

export function HistoryDrawer() {
  const [isOpen, setIsOpen] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);

  useEffect(() => {
    if (isOpen) {
      setHistory(getCheckHistory());
    }
  }, [isOpen]);

  const handleClear = () => {
    clearCheckHistory();
    setHistory([]);
  };

  const formatDate = (iso: string) => {
    try {
      const d = new Date(iso);
      return d.toLocaleDateString(undefined, {
        month: "short",
        day: "numeric",
        hour: "numeric",
        minute: "numeric",
      });
    } catch {
      return "Recent";
    }
  };

  return (
    <div className="history-wrapper">
      <button
        type="button"
        onClick={() => setIsOpen(true)}
        className="history-toggle-btn"
        aria-label="View checks you have done on this phone"
      >
        <History size={18} aria-hidden="true" />
        <span>Checks you have done on this phone</span>
      </button>

      {isOpen ? (
        <div className="history-overlay" role="dialog" aria-modal="true" aria-label="Recent check history">
          <div className="history-backdrop" onClick={() => setIsOpen(false)} />
          <div className="history-sheet">
            <div className="history-sheet-header">
              <h2 className="history-sheet-title">Your recent checks</h2>
              <button
                type="button"
                onClick={() => setIsOpen(false)}
                className="history-close-btn"
                aria-label="Close check history"
              >
                <X size={20} />
              </button>
            </div>

            <div className="history-sheet-content">
              {history.length === 0 ? (
                <p className="history-empty">No saved checks found on this phone.</p>
              ) : (
                <div className="history-list">
                  {history.map((item) => (
                    <div key={item.id} className="history-item">
                      <div className="history-item-top">
                        <span className="history-date">{formatDate(item.timestamp)}</span>
                        <Chip
                          variant={
                            item.severity === "HIGH"
                              ? "high"
                              : item.severity === "LOW"
                              ? "low"
                              : "medium"
                          }
                        >
                          {item.severity} RISK
                        </Chip>
                      </div>
                      <p className="history-summary">{item.summary}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {history.length > 0 ? (
              <div className="history-sheet-footer">
                <button type="button" onClick={handleClear} className="history-clear-btn">
                  <Trash2 size={16} />
                  <span>Clear history on this device</span>
                </button>
              </div>
            ) : null}
          </div>
        </div>
      ) : null}

      <style jsx>{`
        .history-wrapper {
          display: flex;
          justify-content: center;
          margin-top: 12px;
        }

        .history-toggle-btn {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          background: none;
          border: none;
          color: var(--text-muted);
          font-size: 0.95rem;
          font-weight: 600;
          cursor: pointer;
          padding: 8px 12px;
          border-radius: var(--radius-sm);
        }

        .history-toggle-btn:hover {
          color: var(--color-primary);
        }

        .history-overlay {
          position: fixed;
          inset: 0;
          z-index: 100;
          display: flex;
          justify-content: flex-end;
        }

        .history-backdrop {
          position: absolute;
          inset: 0;
          background-color: rgba(0, 0, 0, 0.4);
        }

        .history-sheet {
          position: relative;
          width: 100%;
          max-width: 440px;
          height: 100%;
          background-color: var(--bg-card);
          box-shadow: -4px 0 20px rgba(0, 0, 0, 0.15);
          display: flex;
          flex-direction: column;
          z-index: 101;
        }

        .history-sheet-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 20px 24px;
          border-bottom: 1px solid var(--border-light);
        }

        .history-sheet-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--text-main);
        }

        .history-close-btn {
          background: none;
          border: none;
          cursor: pointer;
          color: var(--text-muted);
          padding: 6px;
          border-radius: var(--radius-sm);
        }

        .history-sheet-content {
          flex: 1;
          overflow-y: auto;
          padding: 20px 24px;
        }

        .history-empty {
          color: var(--text-muted);
          font-size: 1rem;
          text-align: center;
          padding-top: 40px;
        }

        .history-list {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .history-item {
          background-color: var(--bg-page);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-md);
          padding: 16px;
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .history-item-top {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .history-date {
          font-size: 0.88rem;
          color: var(--text-muted);
          font-weight: 500;
        }

        .history-summary {
          font-size: 0.98rem;
          color: var(--text-main);
          font-weight: 600;
          line-height: 1.4;
        }

        .history-sheet-footer {
          padding: 16px 24px;
          border-top: 1px solid var(--border-light);
          display: flex;
          justify-content: center;
        }

        .history-clear-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          background: none;
          border: none;
          color: var(--color-danger);
          font-size: 0.9rem;
          font-weight: 600;
          cursor: pointer;
          padding: 6px 12px;
        }
      `}</style>
    </div>
  );
}
