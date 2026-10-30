"use client";

import React, { useState } from "react";
import Link from "next/link";
import { ServiceStatus } from "@/components/service-status";
import { ThemeToggle } from "@/components/theme-toggle";
import { Users, Shield, Menu, X, PhoneCall } from "lucide-react";

interface ShellProps {
  children: React.ReactNode;
  activeSection?: "senior" | "caregiver";
}

export function Shell({ children, activeSection = "senior" }: ShellProps) {
  const isSenior = activeSection === "senior";
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <div className="shell-root">
      <header className="shell-header">
        <div className="shell-header-inner">
          <div className="shell-brand-group">
            <Link href="/" className="shell-brand">
              <span className="brand-name">Sentry</span>
            </Link>
            <div className="brand-status">
              <ServiceStatus />
            </div>
          </div>

          <nav className="shell-nav-desktop" aria-label="Main navigation">
            <Link
              href="/"
              className={`nav-link ${isSenior ? "nav-link-active" : ""}`}
              aria-current={isSenior ? "page" : undefined}
            >
              <span>Check a message</span>
            </Link>
            <Link
              href="/caregiver"
              className={`nav-link ${!isSenior ? "nav-link-active" : ""}`}
              aria-current={!isSenior ? "page" : undefined}
            >
              <Users size={16} />
              <span>Caregiver Circle</span>
            </Link>
          </nav>

          <div className="shell-header-actions">
            <ThemeToggle />
            <button
              type="button"
              className="mobile-menu-toggle"
              aria-label="Toggle navigation menu"
              aria-expanded={mobileMenuOpen}
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? <X size={20} /> : <Menu size={20} />}
            </button>
          </div>
        </div>

        {mobileMenuOpen && (
          <div className="mobile-nav-drawer">
            <Link
              href="/"
              className={`mobile-nav-link ${isSenior ? "mobile-nav-active" : ""}`}
              onClick={() => setMobileMenuOpen(false)}
            >
              <span>Check a message</span>
            </Link>
            <Link
              href="/caregiver"
              className={`mobile-nav-link ${!isSenior ? "mobile-nav-active" : ""}`}
              onClick={() => setMobileMenuOpen(false)}
            >
              <Users size={18} />
              <span>Caregiver Circle</span>
            </Link>
          </div>
        )}
      </header>

      <main className="shell-main" id="main-content">
        <div className="shell-container">{children}</div>
      </main>

      <footer className="shell-footer">
        <div className="footer-inner">
          <div className="footer-col">
            <div className="footer-brand">Sentry</div>
            <p className="footer-desc">
              Immediate, calm guidance before you reply to unexpected messages. Built to protect seniors and families at zero personal cost.
            </p>
          </div>
          <div className="footer-links">
            <div className="footer-link-group">
              <span className="footer-group-title">Official Hotlines</span>
              <a href="https://reportfraud.ftc.gov" target="_blank" rel="noopener noreferrer">US FTC Fraud Line</a>
              <a href="https://www.actionfraud.police.uk" target="_blank" rel="noopener noreferrer">UK Action Fraud</a>
              <a href="https://cybercrime.gov.in" target="_blank" rel="noopener noreferrer">India 1930 Cyber</a>
              <a href="https://www.scamwatch.gov.au" target="_blank" rel="noopener noreferrer">Australia Scamwatch</a>
            </div>
            <div className="footer-link-group">
              <span className="footer-group-title">Family Circle</span>
              <Link href="/caregiver">Caregiver Portal</Link>
              <Link href="/">Check Screenshot</Link>
            </div>
          </div>
        </div>
        <div className="footer-bottom">
          <span>Sentry Security &middot; Zero storage of your screenshots after safety analysis.</span>
        </div>
      </footer>

      <style jsx>{`
        .shell-root {
          min-height: 100vh;
          display: flex;
          flex-direction: column;
          background-color: var(--bg-page);
          color: var(--text-main);
        }

        .shell-header {
          position: sticky;
          top: 0;
          z-index: 100;
          background: var(--chrome-bg);
          backdrop-filter: blur(16px);
          -webkit-backdrop-filter: blur(16px);
          border-bottom: 1px solid var(--border-light);
          transition: background-color 150ms ease, border-color 150ms ease;
        }

        .shell-header-inner {
          max-width: 1200px;
          margin: 0 auto;
          height: 64px;
          padding: 0 24px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 20px;
        }

        .shell-brand-group {
          display: flex;
          align-items: center;
          gap: 16px;
        }

        .shell-brand {
          display: inline-flex;
          align-items: center;
          text-decoration: none;
          color: var(--text-main);
        }

        .brand-name {
          font-size: 1.35rem;
          font-weight: 850;
          letter-spacing: -0.04em;
          background: linear-gradient(135deg, var(--text-main) 30%, var(--text-muted));
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
        }

        .shell-nav-desktop {
          display: flex;
          align-items: center;
          gap: 8px;
        }

        .nav-link {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          padding: 8px 16px;
          border-radius: var(--radius-full);
          font-size: 0.92rem;
          font-weight: 600;
          color: var(--text-muted);
          text-decoration: none;
          transition: color 140ms ease, background-color 140ms ease;
        }

        .nav-link:hover {
          color: var(--text-main);
          background-color: var(--bg-muted);
        }

        .nav-link-active {
          color: var(--text-main);
          background-color: var(--surface-raised);
          box-shadow: var(--shadow-sm);
        }

        .shell-header-actions {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .mobile-menu-toggle {
          display: none;
          align-items: center;
          justify-content: center;
          width: 38px;
          height: 38px;
          border-radius: var(--radius-md);
          background: transparent;
          border: 1px solid var(--border-light);
          color: var(--text-main);
          cursor: pointer;
        }

        .mobile-nav-drawer {
          display: flex;
          flex-direction: column;
          gap: 8px;
          padding: 16px 24px 20px;
          background: var(--bg-card);
          border-bottom: 1px solid var(--border-light);
        }

        .mobile-nav-link {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 12px 16px;
          border-radius: var(--radius-md);
          font-size: 1rem;
          font-weight: 600;
          color: var(--text-muted);
          text-decoration: none;
        }

        .mobile-nav-active {
          color: var(--text-main);
          background: var(--surface-raised);
        }

        .shell-main {
          flex: 1;
        }

        .shell-container {
          max-width: 1200px;
          margin: 0 auto;
          padding: 40px 24px 80px;
        }

        .shell-footer {
          border-top: 1px solid var(--border-light);
          background: var(--bg-card);
          margin-top: auto;
        }

        .footer-inner {
          max-width: 1200px;
          margin: 0 auto;
          padding: 48px 24px 32px;
          display: grid;
          grid-template-columns: 1.4fr 1fr;
          gap: 40px;
        }

        .footer-brand {
          font-size: 1.25rem;
          font-weight: 800;
          letter-spacing: -0.03em;
          color: var(--text-main);
          margin-bottom: 10px;
        }

        .footer-desc {
          font-size: 0.92rem;
          line-height: 1.55;
          color: var(--text-muted);
          max-width: 440px;
        }

        .footer-links {
          display: flex;
          gap: 48px;
          justify-content: flex-end;
        }

        .footer-link-group {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }

        .footer-group-title {
          font-size: 0.8rem;
          font-weight: 750;
          text-transform: uppercase;
          letter-spacing: 0.08em;
          color: var(--text-main);
          margin-bottom: 4px;
        }

        .footer-link-group a {
          font-size: 0.88rem;
          color: var(--text-muted);
          text-decoration: none;
          transition: color 140ms ease;
        }

        .footer-link-group a:hover {
          color: var(--text-main);
        }

        .footer-bottom {
          max-width: 1200px;
          margin: 0 auto;
          padding: 20px 24px;
          border-top: 1px solid var(--border-light);
          font-size: 0.84rem;
          color: var(--text-muted);
          text-align: center;
        }

        @media (max-width: 768px) {
          .shell-nav-desktop {
            display: none;
          }

          .mobile-menu-toggle {
            display: flex;
          }

          .footer-inner {
            grid-template-columns: 1fr;
            gap: 32px;
          }

          .footer-links {
            justify-content: flex-start;
            flex-wrap: wrap;
            gap: 32px;
          }

          .shell-container {
            padding: 24px 16px 60px;
          }
        }
      `}</style>
    </div>
  );
}
