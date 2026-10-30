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
          font-weight: 600;
          cursor: pointer;
          transition: background-color 0.15s ease, border-color 0.15s ease, transform 0.05s ease;
          border: 2px solid transparent;
          text-align: center;
          text-decoration: none;
          touch-action: manipulation;
        }

        .btn-root:active:not(:disabled) {
          transform: translateY(1px);
        }

        .btn-root:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .btn-default {
          min-height: 52px;
          padding: 12px 24px;
          font-size: 1.05rem;
        }

        .btn-large {
          min-height: 62px;
          padding: 16px 32px;
          font-size: 1.2rem;
          border-radius: var(--radius-lg);
        }

        .btn-small {
          min-height: 44px;
          padding: 8px 16px;
          font-size: 0.95rem;
        }

        .btn-full {
          width: 100%;
        }

        .btn-primary {
          background-color: var(--color-primary);
          color: var(--text-inverse);
          border-color: var(--color-primary);
        }

        .btn-primary:hover:not(:disabled) {
          background-color: var(--color-primary-hover);
          border-color: var(--color-primary-hover);
        }

        .btn-secondary {
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          border-color: var(--bg-secondary);
        }

        .btn-secondary:hover:not(:disabled) {
          background-color: #d8e7dc;
        }

        .btn-outline {
          background-color: transparent;
          color: var(--text-main);
          border-color: var(--border-light);
        }

        .btn-outline:hover:not(:disabled) {
          background-color: var(--bg-secondary);
          border-color: var(--color-primary);
        }

        .btn-danger {
          background-color: var(--color-danger);
          color: var(--text-inverse);
          border-color: var(--color-danger);
        }

        .btn-danger:hover:not(:disabled) {
          background-color: #991b1b;
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
          width: 20px;
          height: 20px;
          border: 3px solid currentColor;
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
