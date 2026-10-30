import React from "react";

export interface ContextNoteProps {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  mode?: "paste" | "context";
}

export function ContextNote({ value, onChange, disabled = false, mode = "paste" }: ContextNoteProps) {
  const isPaste = mode === "paste";
  return (
    <div className="context-root">
      <label htmlFor="message-context" className="context-label">
        {isPaste ? "Paste the message text" : "Add extra context (optional)"}
      </label>
      <textarea
        id="message-context"
        className="context-textarea"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={
          isPaste
            ? "Paste the text or link you received. Leave out personal account numbers."
            : "For example: 'They said my account is suspended and asked me to buy gift cards.'"
        }
        disabled={disabled}
        rows={isPaste ? 4 : 3}
      />
      <p className="context-hint">
        {isPaste
          ? "No personal data is retained. The message is analyzed and cleared."
          : "A brief note helps detect urgency tricks, impersonation, and pressure tactics."}
      </p>

      <style jsx>{`
        .context-root {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .context-label {
          color: var(--text-main);
          font-size: 0.95rem;
          font-weight: 700;
        }

        .context-textarea {
          width: 100%;
          min-height: 100px;
          padding: 14px 16px;
          border: 1px solid var(--border-medium);
          border-radius: var(--radius-md);
          background: var(--surface-translucent);
          color: var(--text-main);
          font-family: inherit;
          font-size: 0.95rem;
          line-height: 1.55;
          resize: vertical;
          transition: border-color 150ms ease, box-shadow 150ms ease, background-color 150ms ease;
        }

        .context-textarea::placeholder {
          color: var(--text-subtle);
        }

        .context-textarea:focus {
          outline: none;
          border-color: var(--color-primary);
          background: var(--surface-raised);
          box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.15);
        }

        .context-textarea:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .context-hint {
          color: var(--text-muted);
          font-size: 0.84rem;
          line-height: 1.45;
        }
      `}</style>
    </div>
  );
}
