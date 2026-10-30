"use client";

import React, { useState } from "react";
import { Shell } from "@/components/shell";
import { Button } from "@/components/button";
import { Users, Plus, LogOut } from "lucide-react";

export default function CaregiverPage() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  return (
    <Shell activeSection="caregiver">
      <div className="caregiver-layout">
        <header className="caregiver-header">
          <div className="caregiver-header-text">
            <h1 className="caregiver-title">Caregiver Feed</h1>
            <p className="caregiver-subtitle">
              Monitor suspicious incidents flagged by your linked family members.
            </p>
          </div>

          {isAuthenticated ? (
            <div className="caregiver-header-actions">
              <Button variant="primary" size="default">
                <Plus size={18} aria-hidden="true" />
                <span>Create invite</span>
              </Button>
              <Button variant="ghost" size="default" onClick={() => setIsAuthenticated(false)}>
                <LogOut size={18} aria-hidden="true" />
                <span>Sign out</span>
              </Button>
            </div>
          ) : null}
        </header>

        <div className="caregiver-content">
          {/* Incidents feed and authentication will be wired here */}
        </div>
      </div>

      <style jsx>{`
        .caregiver-layout {
          display: flex;
          flex-direction: column;
          gap: 28px;
        }

        .caregiver-header {
          display: flex;
          align-items: flex-start;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 16px;
          border-bottom: 1px solid var(--border-light);
          padding-bottom: 20px;
        }

        .caregiver-title {
          font-size: 2rem;
          font-weight: 800;
          color: var(--text-main);
          letter-spacing: -0.02em;
        }

        .caregiver-subtitle {
          font-size: 1.05rem;
          color: var(--text-muted);
          margin-top: 4px;
        }

        .caregiver-header-actions {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .caregiver-content {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }
      `}</style>
    </Shell>
  );
}
