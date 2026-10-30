import React, { useState, useEffect } from "react";
import { WifiOff, Phone } from "lucide-react";

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
        An internet connection is required to check screenshots. If you need urgent assistance, you can still call official phone helplines directly.
      </p>

      <a href="tel:911" className="offline-dial">
        <Phone size={18} aria-hidden="true" />
        <span>Call emergency services</span>
      </a>

      <style jsx>{`
        .offline-card {
          background-color: #fff7ed;
          border: 2px solid #ea580c;
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
          background-color: #ffedd5;
          color: #ea580c;
        }

        .offline-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: #9a3412;
        }

        .offline-body {
          font-size: 1.05rem;
          color: var(--text-main);
          line-height: 1.5;
        }

        .offline-dial {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          margin-top: 4px;
          padding: 10px 18px;
          border-radius: var(--radius-md);
          background-color: #ea580c;
          color: #ffffff;
          font-weight: 700;
          text-decoration: none;
          width: fit-content;
        }

        .offline-dial:hover {
          background-color: #c2410c;
        }
      `}</style>
    </div>
  );
}
