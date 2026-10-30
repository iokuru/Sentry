import React from "react";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "outline" | "danger" | "ghost";
  size?: "default" | "large" | "small";
  fullWidth?: boolean;
  loading?: boolean;
}

export function Button({
  children,
  variant = "primary",
  size = "default",
  fullWidth = false,
  loading = false,
  disabled,
  className = "",
  ...props
}: ButtonProps) {
  const isBusy = loading;
  const isDisabled = disabled || loading;

  return (
    <button
      className={`btn-root btn-${variant} btn-${size} ${fullWidth ? "btn-full" : ""} ${className}`}
      disabled={isDisabled}
      aria-busy={isBusy}
      {...props}
    >
      {loading ? (
        <span className="btn-spinner" aria-hidden="true" />
      ) : null}
      <span className="btn-label">{children}</span>

      <style jsx>{`
        .btn-root {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          gap: 10px;
          border-radius: var(--radius-md);
          font-weight: 650;
          cursor: pointer;
          transition: background-color 0.15s ease, border-color 0.15s ease, transform 0.08s ease, box-shadow 0.15s ease;
          border: 1px solid transparent;
          text-align: center;
          text-decoration: none;
          touch-action: manipulation;
        }

        .btn-root:active:not(:disabled) {
          transform: translateY(1px);
        }

        .btn-root:disabled {
          opacity: 0.65;
          cursor: not-allowed;
          box-shadow: none !important;
        }

        .btn-default {
          min-height: 48px;
          padding: 10px 22px;
          font-size: 0.98rem;
        }

        .btn-large {
          min-height: 56px;
          padding: 14px 28px;
          font-size: 1.1rem;
        }

        .btn-small {
          min-height: 38px;
          padding: 6px 14px;
          font-size: 0.88rem;
        }

        .btn-full {
          width: 100%;
        }

        .btn-primary {
          background-color: var(--color-primary);
          color: #ffffff;
          border-color: var(--color-primary);
          box-shadow: 0 2px 10px rgba(16, 185, 129, 0.25);
        }

        .btn-primary:hover:not(:disabled) {
          background-color: var(--color-primary-hover);
          border-color: var(--color-primary-hover);
          box-shadow: 0 4px 16px rgba(16, 185, 129, 0.35);
        }

        .btn-primary:disabled {
          background-color: var(--button-disabled-bg);
          border-color: var(--button-disabled-border);
          color: var(--button-disabled-text);
        }

        .btn-secondary {
          background-color: var(--surface-raised);
          color: var(--text-main);
          border-color: var(--border-light);
        }

        .btn-secondary:hover:not(:disabled) {
          background-color: var(--bg-muted);
          border-color: var(--text-muted);
        }

        .btn-outline {
          background-color: transparent;
          color: var(--text-main);
          border-color: var(--border-medium);
        }

        .btn-outline:hover:not(:disabled) {
          background-color: var(--surface-raised);
          border-color: var(--color-primary);
        }

        .btn-danger {
          background-color: var(--color-danger);
          color: #ffffff;
          border-color: var(--color-danger);
        }

        .btn-danger:hover:not(:disabled) {
          background-color: #b91c1c;
        }

        .btn-ghost {
          background-color: transparent;
          color: var(--text-muted);
          border-color: transparent;
        }

        .btn-ghost:hover:not(:disabled) {
          background-color: var(--bg-muted);
          color: var(--text-main);
        }

        .btn-spinner {
          width: 18px;
          height: 18px;
          border: 2px solid currentColor;
          border-right-color: transparent;
          border-radius: 50%;
          animation: spin 0.75s linear infinite;
        }

        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </button>
  );
}
