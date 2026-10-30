import React from "react";
import { TextAreaField } from "@/components/field";

export interface ContextNoteProps {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
}

export function ContextNote({
  value,
  onChange,
  disabled = false,
}: ContextNoteProps) {
  return (
    <div className="context-root">
      <TextAreaField
        label="What happened? (optional)"
        hint="Type or paste any message, text, or details you were sent."
        placeholder="For example: A text message saying my parcel is on hold, asking me to tap a link and pay a fee."
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        rows={3}
      />

      <style jsx>{`
        .context-root {
          background-color: var(--bg-card);
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 20px;
        }
      `}</style>
    </div>
  );
}
