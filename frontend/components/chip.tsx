import React from "react";

export type ChipVariant = "high" | "medium" | "low" | "neutral" | "primary";

export interface ChipProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: ChipVariant;
}

export function Chip({
  children,
  variant = "neutral",
  className = "",
  ...props
}: ChipProps) {
  return (
    <span className={`chip-root chip-${variant} ${className}`} {...props}>
      {children}

      <style jsx>{`
        .chip-root {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 4px 12px;
          border-radius: var(--radius-full);
          font-size: 0.88rem;
          font-weight: 700;
          letter-spacing: 0.02em;
          text-transform: uppercase;
          line-height: 1.4;
          white-space: nowrap;
        }

        .chip-high {
          background-color: var(--color-danger-subtle);
          color: var(--color-danger);
          border: 1px solid rgba(185, 28, 28, 0.2);
        }

        .chip-medium {
          background-color: var(--color-accent-subtle);
          color: var(--color-accent);
          border: 1px solid rgba(180, 83, 9, 0.2);
        }

        .chip-low {
          background-color: var(--color-safe-subtle);
          color: var(--color-safe);
          border: 1px solid rgba(21, 128, 61, 0.2);
        }

        .chip-neutral {
          background-color: var(--bg-muted);
          color: var(--text-main);
          border: 1px solid var(--border-light);
        }

        .chip-primary {
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          border: 1px solid rgba(29, 107, 79, 0.2);
        }
      `}</style>
    </span>
  );
}
