"use client";

import React, { useState, useEffect, useMemo } from "react";
import { Shell } from "@/components/shell";
import { Button } from "@/components/button";
import { InputField } from "@/components/field";
import { IncidentCard, type IncidentItem } from "@/components/incident-card";
import { InviteModal } from "@/components/invite-modal";
import { Plus, LogOut, Search, Image as ImageIcon, ShieldAlert } from "lucide-react";
import { getStoredSession, signInCaregiver, signOutCaregiver, type AuthSession } from "@/lib/auth";
import { loadAppConfig } from "@/lib/config";

export default function CaregiverPage() {
  const [session, setSession] = useState<AuthSession | null>(null);
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [loginError, setLoginError] = useState<string | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  const [incidents, setIncidents] = useState<IncidentItem[]>([]);
  const [statusFilter, setStatusFilter] = useState<"ALL" | "NEW" | "ACKNOWLEDGED" | "RESOLVED">("ALL");
  const [searchQuery, setSearchQuery] = useState("");
  const [showThumbnails, setShowThumbnails] = useState(true);
  const [showInviteModal, setShowInviteModal] = useState(false);

  useEffect(() => {
    const s = getStoredSession();
    if (s) setSession(s);
  }, []);

  useEffect(() => {
    if (!session) return;

    let isCancelled = false;

    const fetchFeed = async () => {
      try {
        const config = await loadAppConfig();
        const res = await fetch(`${config.apiUrl}/api/caregiver/incidents`, {
          headers: {
            Authorization: `Bearer ${session.token}`,
            Accept: "application/json",
          },
        });

        if (res.ok && !isCancelled) {
          const data = await res.json();
          setIncidents(data.incidents || []);
        }
      } catch {
        // network or mock
      }
    };

    fetchFeed();
    const timer = setInterval(fetchFeed, 8000);
    return () => {
      isCancelled = true;
      clearInterval(timer);
    };
  }, [session]);

  const openCount = useMemo(() => {
    return incidents.filter((i) => i.status === "NEW").length;
  }, [incidents]);

  useEffect(() => {
    if (typeof document !== "undefined") {
      document.title = openCount > 0 ? `(${openCount}) Caregiver Feed - Sentry` : "Caregiver Feed - Sentry";
    }
  }, [openCount]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginEmail.trim() || !loginPassword) return;

    setIsLoggingIn(true);
    setLoginError(null);

    const res = await signInCaregiver(loginEmail, loginPassword);
    setIsLoggingIn(false);

    if (res.success && res.session) {
      setSession(res.session);
    } else {
      setLoginError(res.error || "Sign in failed. Please verify your credentials.");
    }
  };

  const handleSignOut = () => {
    signOutCaregiver();
    setSession(null);
  };

  const handleUpdateStatus = async (id: string, newStatus: "ACKNOWLEDGED" | "RESOLVED") => {
    // optimistic update
    setIncidents((prev) =>
      prev.map((inc) => (inc.incident_id === id ? { ...inc, status: newStatus } : inc))
    );

    if (!session) return;

    try {
      const config = await loadAppConfig();
      await fetch(`${config.apiUrl}/api/caregiver/incidents/${id}/status`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session.token}`,
        },
        body: JSON.stringify({ status: newStatus }),
      });
    } catch {
      // transient failure
    }
  };

  const filteredIncidents = useMemo(() => {
    return incidents.filter((item) => {
      if (statusFilter !== "ALL" && item.status !== statusFilter) return false;
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const textMatch = (item.hold_off || "").toLowerCase().includes(q);
        const nameMatch = (item.senior_name || "").toLowerCase().includes(q);
        const sigMatch = (item.signals || []).some((s) => s.toLowerCase().includes(q));
        if (!textMatch && !nameMatch && !sigMatch) return false;
      }
      return true;
    });
  }, [incidents, statusFilter, searchQuery]);

  return (
    <Shell activeSection="caregiver">
      <div className="caregiver-layout">
        {!session ? (
          <div className="login-card">
            <div className="login-header">
              <span className="login-shield" aria-hidden="true">
                <ShieldAlert size={32} />
              </span>
              <h1 className="login-title">Caregiver Sign In</h1>
              <p className="login-sub">Sign in to view alerts from your linked family members.</p>
            </div>

            <form onSubmit={handleLogin} className="login-form">
              <InputField
                label="Email address"
                type="email"
                placeholder="caregiver@example.com"
                value={loginEmail}
                onChange={(e) => setLoginEmail(e.target.value)}
                required
              />

              <InputField
                label="Password"
                type="password"
                placeholder="Enter your password"
                value={loginPassword}
                onChange={(e) => setLoginPassword(e.target.value)}
                required
              />

              {loginError ? <p className="login-error" role="alert">{loginError}</p> : null}

              <Button type="submit" variant="primary" size="large" fullWidth loading={isLoggingIn}>
                Sign in to caregiver feed
              </Button>
            </form>
          </div>
        ) : (
          <>
            <header className="caregiver-header">
              <div className="caregiver-header-text">
                <h1 className="caregiver-title">
                  Caregiver Feed {openCount > 0 ? <span className="open-badge">{openCount} new</span> : null}
                </h1>
                <p className="caregiver-subtitle">Signed in as {session.email}</p>
              </div>

              <div className="caregiver-header-actions">
                <Button variant="primary" size="default" onClick={() => setShowInviteModal(true)}>
                  <Plus size={18} aria-hidden="true" />
                  <span>Create invite</span>
                </Button>
                <Button variant="ghost" size="default" onClick={handleSignOut}>
                  <LogOut size={18} aria-hidden="true" />
                  <span>Sign out</span>
                </Button>
              </div>
            </header>

            <div className="filter-bar">
              <div className="status-tabs" role="tablist">
                {(["ALL", "NEW", "ACKNOWLEDGED", "RESOLVED"] as const).map((tab) => (
                  <button
                    key={tab}
                    type="button"
                    role="tab"
                    aria-selected={statusFilter === tab}
                    className={`tab-btn ${statusFilter === tab ? "tab-active" : ""}`}
                    onClick={() => setStatusFilter(tab)}
                  >
                    {tab}
                  </button>
                ))}
              </div>

              <div className="search-and-toggles">
                <div className="search-box">
                  <Search size={18} className="search-icon" aria-hidden="true" />
                  <input
                    type="search"
                    placeholder="Search incidents or signals..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="search-input"
                    aria-label="Search incidents"
                  />
                </div>

                <button
                  type="button"
                  className={`thumb-toggle ${showThumbnails ? "thumb-active" : ""}`}
                  onClick={() => setShowThumbnails(!showThumbnails)}
                  aria-label="Toggle screenshot thumbnails"
                >
                  <ImageIcon size={18} />
                  <span>Thumbnails</span>
                </button>
              </div>
            </div>

            <div className="incidents-feed">
              {filteredIncidents.length === 0 ? (
                <div className="empty-feed">
                  <p className="empty-title">No incidents match this filter.</p>
                  <p className="empty-sub">When your linked seniors check suspicious messages, they will appear here.</p>
                </div>
              ) : (
                filteredIncidents.map((inc) => (
                  <IncidentCard
                    key={inc.incident_id}
                    incident={inc}
                    onUpdateStatus={handleUpdateStatus}
                    showThumbnails={showThumbnails}
                  />
                ))
              )}
            </div>

            {showInviteModal && (
              <InviteModal token={session.token} onClose={() => setShowInviteModal(false)} />
            )}
          </>
        )}
      </div>

      <style jsx>{`
        .caregiver-layout {
          display: flex;
          flex-direction: column;
          gap: 24px;
        }

        .login-card {
          max-width: 460px;
          margin: 40px auto;
          width: 100%;
          background-color: var(--bg-card);
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 32px 28px;
          display: flex;
          flex-direction: column;
          gap: 24px;
          box-shadow: var(--shadow-sm);
        }

        .login-header {
          text-align: center;
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
        }

        .login-shield {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 56px;
          height: 56px;
          border-radius: var(--radius-full);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          margin-bottom: 8px;
        }

        .login-title {
          font-size: 1.8rem;
          font-weight: 800;
          color: var(--text-main);
        }

        .login-sub {
          font-size: 1rem;
          color: var(--text-muted);
        }

        .login-form {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .login-error {
          color: var(--color-danger);
          font-size: 0.92rem;
          font-weight: 600;
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
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .open-badge {
          font-size: 0.88rem;
          background-color: var(--color-danger-subtle);
          color: var(--color-danger);
          padding: 3px 10px;
          border-radius: var(--radius-full);
          border: 1px solid rgba(185, 28, 28, 0.2);
        }

        .caregiver-subtitle {
          font-size: 1rem;
          color: var(--text-muted);
          margin-top: 4px;
        }

        .caregiver-header-actions {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .filter-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 16px;
        }

        .status-tabs {
          display: flex;
          background-color: var(--bg-muted);
          padding: 4px;
          border-radius: var(--radius-md);
          gap: 4px;
        }

        .tab-btn {
          padding: 8px 14px;
          border-radius: var(--radius-sm);
          font-size: 0.88rem;
          font-weight: 700;
          border: none;
          background: transparent;
          color: var(--text-muted);
          cursor: pointer;
        }

        .tab-active {
          background-color: var(--bg-card);
          color: var(--color-primary);
          box-shadow: var(--shadow-sm);
        }

        .search-and-toggles {
          display: flex;
          align-items: center;
          gap: 12px;
          flex-wrap: wrap;
        }

        .search-box {
          display: flex;
          align-items: center;
          gap: 8px;
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          padding: 8px 12px;
          border-radius: var(--radius-md);
        }

        .search-icon {
          color: var(--text-muted);
        }

        .search-input {
          border: none;
          outline: none;
          font-size: 0.95rem;
          color: var(--text-main);
          width: 200px;
        }

        .thumb-toggle {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 14px;
          border-radius: var(--radius-md);
          border: 1px solid var(--border-light);
          background-color: var(--bg-card);
          color: var(--text-muted);
          font-size: 0.9rem;
          font-weight: 600;
          cursor: pointer;
        }

        .thumb-active {
          border-color: var(--color-primary);
          color: var(--color-primary);
          background-color: var(--bg-secondary);
        }

        .incidents-feed {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .empty-feed {
          background-color: var(--bg-card);
          border: 1.5px dashed var(--border-light);
          border-radius: var(--radius-lg);
          padding: 48px 24px;
          text-align: center;
        }

        .empty-title {
          font-size: 1.2rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 6px;
        }

        .empty-sub {
          font-size: 1rem;
          color: var(--text-muted);
        }
      `}</style>
    </Shell>
  );
}
