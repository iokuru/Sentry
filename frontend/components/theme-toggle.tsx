"use client";

import React, { useEffect, useState } from "react";
import { Sun, Moon } from "lucide-react";

export function ThemeToggle() {
  const [theme, setTheme] = useState<"light" | "dark">("dark");
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    const saved = localStorage.getItem("sentry_theme") as "light" | "dark" | null;
    if (saved) {
      setTheme(saved);
      document.documentElement.setAttribute("data-theme", saved);
    } else {
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      const initial = prefersDark ? "dark" : "light";
      setTheme(initial);
      document.documentElement.setAttribute("data-theme", initial);
    }
  }, []);

  const toggleTheme = () => {
    const next = theme === "dark" ? "light" : "dark";
    setTheme(next);
    document.documentElement.setAttribute("data-theme", next);
    try {
      localStorage.setItem("sentry_theme", next);
    } catch {
      // storage unavailable
    }
  };

  if (!mounted) {
    return (
      <div className="theme-toggle-placeholder" aria-hidden="true">
        <span className="theme-toggle-btn" />
      </div>
    );
  }

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="theme-toggle-btn"
      aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
      title={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
    >
      <span className="theme-icon-box">
        {theme === "dark" ? (
          <Sun size={17} strokeWidth={2.2} className="theme-icon sun-icon" />
        ) : (
          <Moon size={17} strokeWidth={2.2} className="theme-icon moon-icon" />
        )}
      </span>
      <span className="theme-text">{theme === "dark" ? "Light" : "Dark"}</span>

      <style jsx>{`
        .theme-toggle-btn {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          background-color: var(--bg-muted);
          border: 1px solid var(--border-light);
          padding: 6px 12px;
          border-radius: var(--radius-full);
          color: var(--text-muted);
          font-size: 0.88rem;
          font-weight: 600;
          cursor: pointer;
          transition: background-color 0.15s ease, color 0.15s ease, border-color 0.15s ease;
          min-height: 38px;
        }

        .theme-toggle-btn:hover {
          background-color: var(--bg-secondary);
          color: var(--text-main);
          border-color: var(--border-focus);
        }

        .theme-icon-box {
          display: inline-flex;
          align-items: center;
          justify-content: center;
        }

        .theme-text {
          letter-spacing: 0.01em;
        }

        .theme-toggle-placeholder {
          width: 78px;
          height: 38px;
        }
      `}</style>
    </button>
  );
}
