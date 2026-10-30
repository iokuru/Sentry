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
    <>
      <button type="button" onClick={toggleTheme} className="theme-toggle-btn" aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`} title={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}>
        {theme === "dark" ? <Sun size={17} strokeWidth={1.8} /> : <Moon size={17} strokeWidth={1.8} />}
      </button>
      <style jsx>{`
        .theme-toggle-btn {
          display:grid; place-items:center; width:38px; height:38px;
          background-color:var(--bg-muted);
          border: 1px solid var(--border-light);
          padding:0; border-radius:7px;
          color: var(--text-muted);
          cursor: pointer;
          transition: background-color 0.15s ease, color 0.15s ease, border-color 0.15s ease;
        }

        .theme-toggle-btn:hover {
          background-color: var(--bg-secondary);
          color: var(--text-main);
          border-color: var(--border-focus);
        }

      `}</style>
    </>
  );
}
