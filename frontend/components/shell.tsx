import React from "react";
import Link from "next/link";
import { Shield } from "lucide-react";

interface ShellProps {
  children: React.ReactNode;
  activeSection?: "senior" | "caregiver";
}

export function Shell({ children, activeSection = "senior" }: ShellProps) {
  return (
    <div className="shell-root">
      <header className="shell-header" role="banner">
        <div className="shell-header-inner">
          <Link href="/" className="shell-brand" aria-label="Sentry home">
            <span className="shell-shield-icon" aria-hidden="true">
              <Shield size={28} strokeWidth={2.2} />
            </span>
            <span className="shell-title">Sentry</span>
          </Link>

          <nav className="shell-nav" aria-label="Main navigation">
            {activeSection === "senior" ? (
              <Link href="/caregiver" className="shell-nav-link">
                Caregiver portal
              </Link>
            ) : (
              <Link href="/" className="shell-nav-link">
                Senior check
              </Link>
            )}
          </nav>
        </div>
      </header>

      <main className="shell-main" id="main-content">
        <div className="shell-container">{children}</div>
      </main>

      <footer className="shell-footer" role="contentinfo">
        <div className="shell-footer-inner">
          <p className="shell-reassurance">
            Free and private. Uploads are analyzed immediately and never stored permanently.
          </p>
          <p className="shell-emergency">
            If you are in immediate danger or being pressured to send money right now, pause and call emergency services or a trusted family member.
          </p>
        </div>
      </footer>

      <style jsx>{`
        .shell-root {
          min-height: 100vh;
          display: flex;
          flex-direction: column;
          background-color: var(--bg-page);
        }

        .shell-header {
          border-bottom: 1px solid var(--border-light);
          background-color: var(--bg-card);
        }

        .shell-header-inner {
          max-width: 800px;
          margin: 0 auto;
          padding: 16px 20px;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .shell-brand {
          display: inline-flex;
          align-items: center;
          gap: 12px;
          text-decoration: none;
          color: var(--color-primary);
        }

        .shell-shield-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 44px;
          height: 44px;
          border-radius: var(--radius-md);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
        }

        .shell-title {
          font-size: 1.4rem;
          font-weight: 700;
          letter-spacing: -0.02em;
          color: var(--text-main);
        }

        .shell-nav {
          display: flex;
          align-items: center;
        }

        .shell-nav-link {
          font-size: 1rem;
          font-weight: 600;
          color: var(--color-primary);
          padding: 8px 12px;
          border-radius: var(--radius-sm);
        }

        .shell-main {
          flex: 1;
          display: flex;
          flex-direction: column;
        }

        .shell-container {
          max-width: 740px;
          width: 100%;
          margin: 0 auto;
          padding: 24px 20px 48px;
        }

        .shell-footer {
          border-top: 1px solid var(--border-light);
          background-color: var(--bg-card);
          padding: 24px 20px;
          margin-top: auto;
        }

        .shell-footer-inner {
          max-width: 740px;
          margin: 0 auto;
          display: flex;
          flex-direction: column;
          gap: 10px;
          text-align: center;
        }

        .shell-reassurance {
          font-size: 0.95rem;
          color: var(--text-main);
          font-weight: 500;
        }

        .shell-emergency {
          font-size: 0.88rem;
          color: var(--text-muted);
          line-height: 1.5;
        }
      `}</style>
    </div>
  );
}
