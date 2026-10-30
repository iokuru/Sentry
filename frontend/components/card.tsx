import React from "react";

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "highlight" | "warning" | "danger" | "subtle";
}

export function Card({
  children,
  variant = "default",
  className = "",
  ...props
}: CardProps) {
  return (
    <div className={`card-root card-${variant} ${className}`} {...props}>
      {children}

      <style jsx>{`
        .card-root {
          background-color: var(--bg-card);
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 24px;
          box-shadow: var(--shadow-sm);
        }

        .card-highlight {
          border-color: var(--color-primary);
          background-color: #f6faf7;
        }

        .card-warning {
          border-color: var(--color-accent);
          background-color: #fffbf0;
        }

        .card-danger {
          border-color: var(--color-danger);
          background-color: #fef7f7;
        }

        .card-subtle {
          background-color: var(--bg-muted);
          border-color: transparent;
        }
      `}</style>
    </div>
  );
}

export function CardTitle({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <h3 className={`card-title ${className}`}>
      {children}
      <style jsx>{`
        .card-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 10px;
        }
      `}</style>
    </h3>
  );
}
