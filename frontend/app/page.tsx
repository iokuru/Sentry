"use client";

import React, { useState } from "react";
import { Shell } from "@/components/shell";

export default function SeniorPage() {
  const [hasStarted, setHasStarted] = useState(false);

  return (
    <Shell activeSection="senior">
      <div className="senior-flow">
        <header className="senior-header">
          <h1 className="senior-title">Check a message or screenshot</h1>
          <p className="senior-subtitle">
            If something feels off, drop an image or paste text below. We will help you figure out what to do.
          </p>
        </header>

        <div className="senior-workspace">
          {/* Main intake area will be mounted here */}
        </div>
      </div>

      <style jsx>{`
        .senior-flow {
          display: flex;
          flex-direction: column;
          gap: 28px;
        }

        .senior-header {
          text-align: center;
          display: flex;
          flex-direction: column;
          gap: 12px;
          padding: 8px 0;
        }

        .senior-title {
          font-size: 2.2rem;
          font-weight: 800;
          color: var(--text-main);
          letter-spacing: -0.02em;
        }

        .senior-subtitle {
          font-size: 1.15rem;
          color: var(--text-muted);
          max-width: 580px;
          margin: 0 auto;
          line-height: 1.5;
        }

        .senior-workspace {
          display: flex;
          flex-direction: column;
          gap: 24px;
        }

        @media (max-width: 640px) {
          .senior-title {
            font-size: 1.8rem;
          }
          .senior-subtitle {
            font-size: 1.05rem;
          }
        }
      `}</style>
    </Shell>
  );
}
