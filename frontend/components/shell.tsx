"use client";

import React from "react";
import Link from "next/link";
import { MessageSquareText, UsersRound } from "lucide-react";
import { ThemeToggle } from "@/components/theme-toggle";

interface ShellProps {
  children: React.ReactNode;
  activeSection?: "landing" | "senior" | "caregiver";
}

export function Shell({ children, activeSection = "senior" }: ShellProps) {
  const isSenior = activeSection === "senior";
  const isCaregiver = activeSection === "caregiver";

  return (
    <div className="app-shell">
      <aside className="app-rail">
        <Link href="/" className="rail-brand">Sentry</Link>
        <nav className="rail-nav" aria-label="Main navigation">
          <h2 className="rail-heading">Workspace</h2>
          <Link href="/check" className={isSenior ? "rail-link rail-link-active" : "rail-link"} aria-current={isSenior ? "page" : undefined}>
            <MessageSquareText size={19} aria-hidden="true" />
            <span>Message check</span>
          </Link>
          <Link href="/caregiver" className={isCaregiver ? "rail-link rail-link-active" : "rail-link"} aria-current={isCaregiver ? "page" : undefined}>
            <UsersRound size={19} aria-hidden="true" />
            <span>For family</span>
          </Link>
        </nav>
        <div className="rail-bottom">
          <p>Need another pair of eyes?</p>
          <Link href={isSenior ? "/caregiver" : "/check"} className="rail-help-link">
            {isSenior ? "Open the family space" : "Open message check"}
          </Link>
        </div>
      </aside>

      <div className="app-main-column">
        <header className="app-topbar">
          <div className="topbar-title">{isSenior ? "Message check" : isCaregiver ? "Family space" : ""}</div>
          <div className="topbar-controls">
            <ThemeToggle />
          </div>
        </header>
        <main className="app-main" id="main-content">
          <div className="shell-container">{children}</div>
        </main>
      </div>

      <style jsx>{`
        .app-shell { min-height:100vh; display:grid; grid-template-columns:258px minmax(0,1fr); background:var(--bg-page); color:var(--text-main); }
        .app-rail { position:sticky; top:0; display:flex; flex-direction:column; width:258px; height:100vh; padding:20px 16px; background:var(--chrome-bg); border-right:1px solid var(--border-light); }
        :global(.rail-brand) { display:inline-flex; align-items:center; min-height:60px; min-width:60px; margin:0 10px 20px; color:var(--text-main)!important; font-size:20px; font-weight:700; letter-spacing:-.025em; text-decoration:none; }
        .rail-nav { display:grid; gap:5px; }
        .rail-heading { margin:0 10px 9px; color:var(--text-main); font-size:15px; font-weight:650; }
        :global(.rail-link) { display:inline-flex; align-items:center; gap:12px; min-height:60px; min-width:60px; padding:0 14px; border:1px solid transparent; border-radius:7px; color:var(--text-muted)!important; font-size:14px; font-weight:540; text-decoration:none; transition:background 140ms ease,color 140ms ease,border-color 140ms ease; }
        :global(.rail-link svg) { flex:none; }
        :global(.rail-link:hover) { background:var(--bg-muted); color:var(--text-main)!important; }
        :global(.rail-link-active) { background:var(--bg-muted); border-color:var(--border-light); color:var(--text-main)!important; }
        :global(.rail-link-active svg) { color:var(--color-primary); }
        .rail-bottom { margin-top:auto; padding:16px 10px 2px; border-top:1px solid var(--border-light); }
        .rail-bottom p { color:var(--text-muted); font-size:14px; line-height:1.4; }
        :global(.rail-help-link) { display:inline-block; margin-top:8px; color:var(--color-primary); font-size:14px; font-weight:600; text-decoration:none; }
        :global(.rail-help-link:hover) { text-decoration:underline; }
        .app-main-column { min-width:0; }
        .app-topbar { position:sticky; top:0; z-index:30; display:flex; align-items:center; justify-content:space-between; min-height:60px; padding:0 clamp(20px,3vw,48px); border-bottom:1px solid var(--border-light); background:var(--bg-page); }
        .topbar-title { color:var(--text-main); font-size:15px; font-weight:600; }
        .topbar-controls { display:flex; align-items:center; justify-content:flex-end; min-width:60px; min-height:60px; }
        .app-main { min-width:0; }
        .shell-container { width:100%; max-width:1500px; margin:0 auto; padding:24px clamp(20px,3.2vw,48px) 64px; }
        @media(max-width:900px) { .app-shell { grid-template-columns:220px minmax(0,1fr); } .app-rail { width:220px; padding-inline:12px; } .shell-container { padding:24px 20px 52px; } }
        @media(max-width:700px) {
          .app-shell { display:flex; flex-direction:column; }
          .app-rail { position:static; width:100%; height:auto; padding:6px 14px; border-right:0; border-bottom:1px solid var(--border-light); flex-direction:row; align-items:center; justify-content:space-between; }
          :global(.rail-brand) { margin:0; font-size:18px; min-height:60px; }
          .rail-nav { display:flex; align-items:center; gap:6px; }
          .rail-heading { display:none; }
          :global(.rail-link) { min-height:60px; min-width:60px; padding:0 12px; font-size:14px; }
          :global(.rail-link svg) { width:18px; height:18px; }
          .rail-bottom { display:none; }
          .app-topbar { position:static; min-height:60px; padding:0 16px; }
          .topbar-title { font-size:15px; }
          .shell-container { padding:20px 16px 40px; }
        }
        @media(max-width:400px) { :global(.rail-link) { gap:6px; padding:0 8px; font-size:13px; } }
      `}</style>
    </div>
  );
}
