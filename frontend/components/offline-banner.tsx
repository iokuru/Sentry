import React, { useState, useEffect } from "react";
import { WifiOff } from "lucide-react";

export function OfflineBanner() {
  const [isOffline, setIsOffline] = useState(false);

  useEffect(() => {
    const updateStatus = () => {
      setIsOffline(!navigator.onLine);
    };

    updateStatus();
    window.addEventListener("online", updateStatus);
    window.addEventListener("offline", updateStatus);

    return () => {
      window.removeEventListener("online", updateStatus);
      window.removeEventListener("offline", updateStatus);
    };
  }, []);

  if (!isOffline) return null;

  return (
    <div className="offline-card" role="alert">
      <div className="offline-header">
        <span className="offline-icon" aria-hidden="true">
          <WifiOff size={24} strokeWidth={2.4} />
        </span>
        <h2 className="offline-title">You are currently offline</h2>
      </div>

      <p className="offline-body">
        An internet connection is required to check this message. If you need urgent help, call your local emergency number or someone you trust.
      </p>

      <style jsx>{`
        .offline-card {
          background-color: var(--warning-bg);
          border: 1px solid var(--warning-border);
          border-radius: var(--radius-lg);
          padding: 20px 24px;
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        .offline-header {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .offline-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 38px;
          height: 38px;
          border-radius: var(--radius-md);
          background-color: var(--warning-icon-bg);
          color: var(--warning-icon-ink);
        }

        .offline-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--warning-title);
        }

        .offline-body {
          font-size: 1.05rem;
          color: var(--text-main);
          line-height: 1.5;
        }

      `}</style>
    </div>
  );
}
