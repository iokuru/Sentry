import React, { useId } from "react";

export interface InputFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  hint?: string;
  error?: string;
}

export function InputField({
  label,
  hint,
  error,
  id,
  className = "",
  ...props
}: InputFieldProps) {
  const generatedId = useId();
  const fieldId = id || generatedId;
  const hintId = hint ? `${fieldId}-hint` : undefined;
  const errorId = error ? `${fieldId}-error` : undefined;

  return (
    <div className={`field-group ${className}`}>
      <label htmlFor={fieldId} className="field-label">
        {label}
      </label>
      {hint ? (
        <span id={hintId} className="field-hint">
          {hint}
        </span>
      ) : null}
      <input
        id={fieldId}
        aria-describedby={[hintId, errorId].filter(Boolean).join(" ") || undefined}
        aria-invalid={error ? "true" : undefined}
        className={`field-control ${error ? "field-error" : ""}`}
        {...props}
      />
      {error ? (
        <span id={errorId} className="field-error-msg" role="alert">
          {error}
        </span>
      ) : null}

      <style jsx>{`
        .field-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        .field-label {
          font-size: 1.05rem;
          font-weight: 600;
          color: var(--text-main);
        }

        .field-hint {
          font-size: 0.92rem;
          color: var(--text-muted);
        }

        .field-control {
          min-height: 52px;
          padding: 12px 16px;
          border-radius: var(--radius-md);
          border: 1.5px solid var(--border-light);
          background-color: var(--bg-card);
          color: var(--text-main);
          font-size: 1.05rem;
          transition: border-color 0.15s ease;
        }

        .field-control:focus {
          border-color: var(--color-primary);
        }

        .field-error {
          border-color: var(--color-danger);
        }

        .field-error-msg {
          font-size: 0.92rem;
          color: var(--color-danger);
          font-weight: 600;
        }
      `}</style>
    </div>
  );
}

export interface TextAreaFieldProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label: string;
  hint?: string;
  error?: string;
}

export function TextAreaField({
  label,
  hint,
  error,
  id,
  className = "",
  ...props
}: TextAreaFieldProps) {
  const generatedId = useId();
  const fieldId = id || generatedId;
  const hintId = hint ? `${fieldId}-hint` : undefined;
  const errorId = error ? `${fieldId}-error` : undefined;

  return (
    <div className={`field-group ${className}`}>
      <label htmlFor={fieldId} className="field-label">
        {label}
      </label>
      {hint ? (
        <span id={hintId} className="field-hint">
          {hint}
        </span>
      ) : null}
      <textarea
        id={fieldId}
        aria-describedby={[hintId, errorId].filter(Boolean).join(" ") || undefined}
        aria-invalid={error ? "true" : undefined}
        className={`field-control field-textarea ${error ? "field-error" : ""}`}
        rows={4}
        {...props}
      />
      {error ? (
        <span id={errorId} className="field-error-msg" role="alert">
          {error}
        </span>
      ) : null}

      <style jsx>{`
        .field-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        .field-label {
          font-size: 1.05rem;
          font-weight: 600;
          color: var(--text-main);
        }

        .field-hint {
          font-size: 0.92rem;
          color: var(--text-muted);
        }

        .field-control {
          padding: 14px 16px;
          border-radius: var(--radius-md);
          border: 1.5px solid var(--border-light);
          background-color: var(--bg-card);
          color: var(--text-main);
          font-size: 1.05rem;
          line-height: 1.5;
          resize: vertical;
        }

        .field-control:focus {
          border-color: var(--color-primary);
        }

        .field-error {
          border-color: var(--color-danger);
        }

        .field-error-msg {
          font-size: 0.92rem;
          color: var(--color-danger);
          font-weight: 600;
        }
      `}</style>
    </div>
  );
}
