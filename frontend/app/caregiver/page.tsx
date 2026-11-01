"use client";

import React, { useState, useEffect, useMemo } from "react";
import { Shell } from "@/components/shell";
import { Button } from "@/components/button";
import { InputField } from "@/components/field";
import { IncidentCard, type IncidentItem } from "@/components/incident-card";
import { InviteModal } from "@/components/invite-modal";
import { Plus, LogOut, Search, Image as ImageIcon, Users } from "lucide-react";
import {
  getStoredSession,
  signInCaregiver,
  signUpCaregiver,
  confirmSignUpCaregiver,
  resendConfirmationCode,
  signOutCaregiver,
  type AuthSession,
} from "@/lib/auth";
import { loadAppConfig } from "@/lib/config";
import { normalizeIncident } from "@/lib/api";

export default function CaregiverPage() {
  const [session, setSession] = useState<AuthSession | null>(null);
  const [authMode, setAuthMode] = useState<"signin" | "signup" | "confirm">("signin");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [confirmCode, setConfirmCode] = useState("");
  const [authSuccessMsg, setAuthSuccessMsg] = useState<string | null>(null);
  const [loginError, setLoginError] = useState<string | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [isResending, setIsResending] = useState(false);

  const [incidents, setIncidents] = useState<IncidentItem[]>([]);
  const [statusFilter, setStatusFilter] = useState<"ALL" | "NEW" | "ACKNOWLEDGED" | "RESOLVED">("ALL");
  const [searchQuery, setSearchQuery] = useState("");
  const [showThumbnails, setShowThumbnails] = useState(true);
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [feedError, setFeedError] = useState<string | null>(null);

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

        if (!res.ok) throw new Error("Could not load family alerts.");
        if (!isCancelled) {
          const data = await res.json();
          const rows = Array.isArray(data) ? data : data.incidents || [];
          setIncidents(rows.map((raw: any) => {
            const incident = normalizeIncident(raw);
            return {
              ...incident,
              status: incident.status === "ANALYZING" ? "NEW" : incident.status,
              senior_id: raw.senior_id,
              senior_name: raw.senior_name || raw.senior_id,
              thumbnail_url: raw.thumbnail_url ? new URL(raw.thumbnail_url, config.apiUrl).toString() : undefined,
            };
          }));
          setFeedError(null);
        }
      } catch {
        if (!isCancelled) setFeedError("Family alerts could not be loaded. Check your connection and try again.");
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
      document.title = openCount > 0 ? `(${openCount}) Caregiver Feed — Sentry` : "Caregiver Feed — Sentry";
    }
  }, [openCount]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginEmail.trim() || !loginPassword) return;

    setIsLoggingIn(true);
    setLoginError(null);
    setAuthSuccessMsg(null);

    const res = await signInCaregiver(loginEmail, loginPassword);
    setIsLoggingIn(false);

    if (res.success && res.session) {
      setSession(res.session);
    } else {
      setLoginError(res.error || "Sign in failed. Please verify your credentials.");
    }
  };

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginEmail.trim() || !loginPassword) return;

    setIsLoggingIn(true);
    setLoginError(null);
    setAuthSuccessMsg(null);

    const res = await signUpCaregiver(loginEmail, loginPassword);
    setIsLoggingIn(false);

    if (res.success) {
      if (res.needsConfirmation) {
        setAuthMode("confirm");
        setAuthSuccessMsg(`A verification code was sent to ${loginEmail.trim()}. Enter it below to activate your account.`);
      } else {
        const loginRes = await signInCaregiver(loginEmail, loginPassword);
        if (loginRes.success && loginRes.session) {
          setSession(loginRes.session);
        }
      }
    } else {
      setLoginError(res.error || "Sign up failed. Password must be at least 8 characters with numbers and special symbols.");
    }
  };

  const handleConfirmSignUp = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginEmail.trim() || !confirmCode.trim()) return;

    setIsLoggingIn(true);
    setLoginError(null);

    const res = await confirmSignUpCaregiver(loginEmail, confirmCode);
    if (!res.success) {
      setIsLoggingIn(false);
      setLoginError(res.error || "Verification code is invalid or expired.");
      return;
    }

    const loginRes = await signInCaregiver(loginEmail, loginPassword);
    setIsLoggingIn(false);

    if (loginRes.success && loginRes.session) {
      setSession(loginRes.session);
    } else {
      setAuthMode("signin");
      setAuthSuccessMsg("Account verified successfully. Please sign in.");
    }
  };

  const handleResendCode = async () => {
    if (!loginEmail.trim()) return;
    setIsResending(true);
    setLoginError(null);
    const res = await resendConfirmationCode(loginEmail);
    setIsResending(false);
    if (res.success) {
      setAuthSuccessMsg(`A new code was sent to ${loginEmail.trim()}.`);
    } else {
      setLoginError(res.error || "Could not resend code.");
    }
  };

  const handleSignOut = () => {
    signOutCaregiver();
    setSession(null);
    setLoginEmail("");
    setLoginPassword("");
    setConfirmCode("");
    setAuthSuccessMsg(null);
    setLoginError(null);
    setAuthMode("signin");
  };

  const handleUpdateStatus = async (id: string, newStatus: "ACKNOWLEDGED" | "RESOLVED") => {
    if (!session) return;

    try {
      const config = await loadAppConfig();
      const res = await fetch(`${config.apiUrl}/api/caregiver/status`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session.token}`,
        },
        body: JSON.stringify({ incident_id: id, status: newStatus.toLowerCase() }),
      });
      if (!res.ok) throw new Error("Could not update alert status.");
      setIncidents((prev) => prev.map((inc) => (inc.incident_id === id ? { ...inc, status: newStatus } : inc)));
      setFeedError(null);
    } catch {
      setFeedError("That alert could not be updated. Please try again.");
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
              <span className="login-icon-pill" aria-hidden="true">
                <Users size={24} />
              </span>
              <h1 className="login-title">
                {authMode === "confirm"
                  ? "Verify your account"
                  : authMode === "signup"
                  ? "Create caregiver account"
                  : "Caregiver portal"}
              </h1>
              <p className="login-sub">
                {authMode === "confirm"
                  ? "Enter the 6-digit confirmation code sent to your email."
                  : authMode === "signup"
                  ? "Sign up to receive alerts and look after your family members."
                  : "Sign in to monitor alerts and pair with senior devices."}
              </p>
            </div>

            {authMode !== "confirm" && (
              <div className="auth-tab-bar" role="tablist">
                <button
                  type="button"
                  role="tab"
                  aria-selected={authMode === "signin"}
                  className={`auth-tab-btn ${authMode === "signin" ? "auth-tab-active" : ""}`}
                  onClick={() => {
                    setAuthMode("signin");
                    setLoginError(null);
                    setAuthSuccessMsg(null);
                  }}
                >
                  Sign in
                </button>
                <button
                  type="button"
                  role="tab"
                  aria-selected={authMode === "signup"}
                  className={`auth-tab-btn ${authMode === "signup" ? "auth-tab-active" : ""}`}
                  onClick={() => {
                    setAuthMode("signup");
                    setLoginError(null);
                    setAuthSuccessMsg(null);
                  }}
                >
                  Create account
                </button>
              </div>
            )}

            {authSuccessMsg && <p className="auth-success" role="status">{authSuccessMsg}</p>}
            {loginError && <p className="login-error" role="alert">{loginError}</p>}

            {authMode === "confirm" ? (
              <form onSubmit={handleConfirmSignUp} className="login-form">
                <InputField
                  label="Email address"
                  type="email"
                  value={loginEmail}
                  onChange={(e) => setLoginEmail(e.target.value)}
                  disabled
                />
                <InputField
                  label="Verification code"
                  type="text"
                  placeholder="Enter 6-digit code"
                  value={confirmCode}
                  onChange={(e) => setConfirmCode(e.target.value)}
                  required
                />
                <Button type="submit" variant="primary" size="large" fullWidth loading={isLoggingIn}>
                  Confirm and sign in
                </Button>
                <div className="auth-footer-actions">
                  <button type="button" className="auth-link-btn" onClick={handleResendCode} disabled={isResending}>
                    {isResending ? "Resending..." : "Resend code"}
                  </button>
                  <span className="auth-dot-sep">·</span>
                  <button type="button" className="auth-link-btn" onClick={() => setAuthMode("signin")}>
                    Back to sign in
                  </button>
                </div>
              </form>
            ) : (
              <form onSubmit={authMode === "signup" ? handleSignUp : handleLogin} className="login-form">
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
                  placeholder={authMode === "signup" ? "Min 8 chars with numbers and symbols" : "Enter your password"}
                  value={loginPassword}
                  onChange={(e) => setLoginPassword(e.target.value)}
                  required
                />
                <Button type="submit" variant="primary" size="large" fullWidth loading={isLoggingIn}>
                  {authMode === "signup" ? "Create account" : "Sign in to caregiver feed"}
                </Button>
                <p className="auth-toggle-note">
                  {authMode === "signup" ? (
                    <>Already have an account? <button type="button" className="auth-link-btn" onClick={() => setAuthMode("signin")}>Sign in</button></>
                  ) : (
                    <>New to Sentry? <button type="button" className="auth-link-btn" onClick={() => setAuthMode("signup")}>Create an account</button></>
                  )}
                </p>
              </form>
            )}
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
                  <Plus size={16} aria-hidden="true" />
                  <span>Create invite</span>
                </Button>
                <Button variant="ghost" size="default" onClick={handleSignOut}>
                  <LogOut size={16} aria-hidden="true" />
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
                  <Search size={16} className="search-icon" aria-hidden="true" />
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
                  <ImageIcon size={16} />
                  <span>Thumbnails</span>
                </button>
              </div>
            </div>

            {feedError ? <p className="feed-error" role="status">{feedError}</p> : null}
            <div className="incidents-feed">
              {filteredIncidents.length === 0 ? (
                <div className="empty-feed">
                  <p className="empty-title">No incidents match this filter.</p>
                  <p className="empty-sub">When your linked family members check suspicious messages, high-risk alerts appear here.</p>
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
          gap: 28px;
          max-width: 900px;
          margin: 0 auto;
          width: 100%;
        }

        .login-card {
          max-width: 440px;
          margin: 40px auto;
          width: 100%;
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-xl);
          padding: 36px 32px;
          display: flex;
          flex-direction: column;
          gap: 24px;
          box-shadow: var(--shadow-lg);
          backdrop-filter: blur(12px);
        }

        .login-header {
          text-align: center;
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 8px;
        }

        .login-icon-pill {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 48px;
          height: 48px;
          border-radius: var(--radius-full);
          background-color: var(--surface-raised);
          border: 1px solid var(--border-light);
          color: var(--color-primary);
          margin-bottom: 4px;
        }

        .login-title {
          font-size: 1.65rem;
          font-weight: 850;
          letter-spacing: -0.03em;
          color: var(--text-main);
        }

        .login-sub {
          font-size: 0.94rem;
          color: var(--text-muted);
          line-height: 1.45;
        }

        .login-form {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .login-error {
          color: var(--color-danger);
          font-size: 0.88rem;
          font-weight: 600;
        }

        .auth-tab-bar {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 6px;
          background: var(--surface-raised);
          padding: 4px;
          border-radius: var(--radius-md);
          border: 1px solid var(--border-light);
        }

        .auth-tab-btn {
          border: none;
          background: transparent;
          color: var(--text-muted);
          padding: 8px 12px;
          font-size: 0.92rem;
          font-weight: 600;
          border-radius: calc(var(--radius-md) - 2px);
          cursor: pointer;
          transition: background 150ms ease, color 150ms ease;
        }

        .auth-tab-btn:hover {
          color: var(--text-main);
        }

        .auth-tab-active {
          background: var(--bg-card);
          color: var(--text-main);
          box-shadow: var(--shadow-sm);
        }

        .auth-success {
          background: var(--safe-bg);
          border: 1px solid var(--safe-border);
          color: var(--safe-title);
          padding: 10px 14px;
          border-radius: var(--radius-md);
          font-size: 0.88rem;
          line-height: 1.45;
        }

        .auth-footer-actions {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 10px;
          margin-top: 4px;
        }

        .auth-toggle-note {
          text-align: center;
          color: var(--text-muted);
          font-size: 0.88rem;
          margin-top: 4px;
        }

        .auth-link-btn {
          background: none;
          border: none;
          color: var(--color-primary);
          font-size: inherit;
          font-weight: 650;
          cursor: pointer;
          padding: 0;
          text-decoration: underline;
        }

        .auth-link-btn:hover {
          opacity: 0.85;
        }

        .auth-dot-sep {
          color: var(--text-subtle);
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
          font-size: 1.85rem;
          font-weight: 850;
          letter-spacing: -0.03em;
          color: var(--text-main);
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .open-badge {
          font-size: 0.82rem;
          font-weight: 750;
          background-color: var(--danger-bg);
          color: var(--danger-title);
          padding: 3px 10px;
          border-radius: var(--radius-full);
          border: 1px solid rgba(239, 68, 68, 0.3);
        }

        .caregiver-subtitle {
          font-size: 0.92rem;
          color: var(--text-muted);
          margin-top: 4px;
        }

        .caregiver-header-actions {
          display: flex;
          align-items: center;
          gap: 10px;
        }

        .filter-bar {
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 16px;
        }

        .status-tabs {
          display: inline-flex;
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-full);
          padding: 3px;
          gap: 3px;
        }

        .tab-btn {
          padding: 6px 14px;
          border-radius: var(--radius-full);
          font-size: 0.84rem;
          font-weight: 650;
          border: none;
          background: transparent;
          color: var(--text-muted);
          cursor: pointer;
          transition: color 140ms ease, background-color 140ms ease;
        }

        .tab-btn:hover {
          color: var(--text-main);
        }

        .tab-active {
          background-color: var(--bg-card);
          color: var(--text-main);
          box-shadow: var(--shadow-sm);
        }

        .search-and-toggles {
          display: flex;
          align-items: center;
          gap: 10px;
          flex-wrap: wrap;
        }

        .search-box {
          display: flex;
          align-items: center;
          gap: 8px;
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          padding: 6px 12px;
          border-radius: var(--radius-full);
        }

        .search-icon {
          color: var(--text-muted);
        }

        .search-input {
          border: none;
          outline: none;
          font-size: 0.88rem;
          color: var(--text-main);
          background: transparent;
          width: 180px;
        }

        .thumb-toggle {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 6px 12px;
          border-radius: var(--radius-full);
          border: 1px solid var(--border-light);
          background-color: var(--surface-raised);
          color: var(--text-muted);
          font-size: 0.86rem;
          font-weight: 600;
          cursor: pointer;
          transition: border-color 140ms ease, color 140ms ease;
        }

        .thumb-toggle:hover {
          color: var(--text-main);
        }

        .thumb-active {
          border-color: var(--color-primary);
          color: var(--color-primary);
        }

        .incidents-feed {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .feed-error { color:var(--color-danger); font-size:0.9rem; }

        .empty-feed {
          background-color: var(--bg-card);
          border: 1px dashed var(--border-medium);
          border-radius: var(--radius-lg);
          padding: 48px 24px;
          text-align: center;
        }

        .empty-title {
          font-size: 1.15rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 6px;
        }

        .empty-sub {
          font-size: 0.94rem;
          color: var(--text-muted);
        }

        @media (max-width: 600px) {
          .caregiver-header {
            flex-direction: column;
            align-items: stretch;
          }
          .filter-bar {
            flex-direction: column;
            align-items: stretch;
          }
          .search-box {
            width: 100%;
          }
          .search-input {
            width: 100%;
          }
        }
      `}</style>
    </Shell>
  );
}
