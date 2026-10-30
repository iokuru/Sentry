"use client";

import React, { useEffect, useState } from "react";
import { Check, LoaderCircle, X } from "lucide-react";
import { loadAppConfig } from "@/lib/config";

type ServiceState = "checking" | "ready" | "offline";

const stateCopy: Record<ServiceState, { label: string; detail: string }> = {
  checking: { label: "Checking", detail: "Checking the message check service." },
  ready: { label: "Ready", detail: "The message check service is ready to analyze screenshots." },
  offline: { label: "Offline", detail: "The service is temporarily unavailable. You can still pause and contact someone you trust." },
};

export function ServiceStatus() {
  const [state, setState] = useState<ServiceState>("checking");
  const [expanded, setExpanded] = useState(false);

  useEffect(() => {
    let stopped = false;
    let timer: ReturnType<typeof setTimeout> | undefined;

    const check = async () => {
      if (!navigator.onLine) {
        if (!stopped) setState("offline");
      } else {
        try {
          const config = await loadAppConfig();
          const controller = new AbortController();
          const timeout = setTimeout(() => controller.abort(), 5000);
          const response = await fetch(`${config.apiUrl.replace(/\/$/, "")}/api/health`, {
            signal: controller.signal,
            cache: "no-store",
          });
          clearTimeout(timeout);
          if (!response.ok) throw new Error("Service health check failed");
          const body = await response.json();
          if (body.status !== "healthy") throw new Error("Service is not ready");
          if (!stopped) setState("ready");
        } catch {
          if (!stopped) setState("offline");
        }
      }

      if (!stopped) timer = setTimeout(check, 30_000);
    };

    const refresh = () => {
      if (timer) clearTimeout(timer);
      void check();
    };

    window.addEventListener("online", refresh);
    window.addEventListener("offline", refresh);
    void check();

    return () => {
      stopped = true;
      if (timer) clearTimeout(timer);
      window.removeEventListener("online", refresh);
      window.removeEventListener("offline", refresh);
    };
  }, []);

  const Icon = state === "ready" ? Check : state === "offline" ? X : LoaderCircle;

  return (
    <div className={`service-status status-${state}`}>
      <button
        type="button"
        className="ready-pill"
        aria-expanded={expanded}
        aria-label={`Message check service status: ${stateCopy[state].label}`}
        onClick={() => setExpanded((value) => !value)}
      >
        <span className="ready-dot" aria-hidden="true">
          <Icon size={15} strokeWidth={2.4} />
        </span>
      </button>

      {expanded && (
        <div className="status-detail" role="status">
          <span className="status-detail-light" />
          <div>
            <strong>Message Check Service</strong>
            <p>{stateCopy[state].detail}</p>
          </div>
        </div>
      )}

      <style jsx>{`
        .service-status {
          position: relative;
        }

        .ready-pill {
          display:grid; place-items:center; width:38px; height:38px;
          background:var(--bg-muted);
          border: 1px solid var(--border-light);
          padding:0; border-radius:7px;
          cursor: pointer;
          transition: transform 140ms ease, border-color 140ms ease;
        }

        .ready-pill:hover {
          border-color: var(--text-muted);
        }

        .ready-dot {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 24px;
          height: 24px;
          border-radius: var(--radius-full);
          background-color: var(--status-color);
          color: #ffffff;
          box-shadow:none;
        }

        .status-ready {
          --status-color: #10b981;
          --status-glow: rgba(16, 185, 129, 0.5);
        }

        .status-checking {
          --status-color: #eab308;
          --status-glow: rgba(234, 179, 8, 0.4);
        }

        .status-offline {
          --status-color: #ef4444;
          --status-glow: rgba(239, 68, 68, 0.4);
        }

        .status-checking .ready-dot :global(svg) {
          animation: status-spin 1.2s linear infinite;
        }

        .status-detail {
          position: absolute;
          z-index: 50;
          top: calc(100% + 10px);
          right: 0;
          width: min(300px, calc(100vw - 36px));
          display: grid;
          grid-template-columns: 12px 1fr;
          align-items: start;
          gap: 12px;
          padding: 16px;
          border: 1px solid var(--border-light);
          border-radius: var(--radius-md);
          background: var(--panel-bg);
          color: var(--text-main);
          box-shadow: var(--shadow-lg);
          backdrop-filter: blur(12px);
        }

        .status-detail-light {
          width: 10px;
          height: 10px;
          margin-top: 4px;
          border-radius: 50%;
          background: var(--status-color);
          box-shadow: 0 0 6px var(--status-glow);
        }

        .status-detail strong {
          display: block;
          font-size: 0.94rem;
          color: var(--text-main);
        }

        .status-detail p {
          margin-top: 4px;
          color: var(--text-muted);
          font-size: 0.86rem;
          line-height: 1.45;
        }

        @keyframes status-spin {
          to { transform: rotate(360deg); }
        }

        @media (prefers-reduced-motion: reduce) {
          .status-checking .ready-dot :global(svg) { animation: none; }
        }
      `}</style>
    </div>
  );
}
