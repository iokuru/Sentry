import React from "react";
import { Button } from "@/components/button";
import { ArrowRight, LoaderCircle } from "lucide-react";

export interface CheckButtonProps {
  onClick: () => void;
  loading?: boolean;
  disabled?: boolean;
  hasInput: boolean;
}

export function CheckButton({
  onClick,
  loading = false,
  disabled = false,
  hasInput,
}: CheckButtonProps) {
  return (
    <div className="check-btn-wrapper">
      <Button
        variant="primary"
        size="large"
        fullWidth
        loading={loading}
        disabled={disabled || !hasInput}
        onClick={onClick}
        aria-label={loading ? "Checking this message" : "Check this message"}
      >
        <span className="check-btn-inner">
          {!loading ? <ArrowRight size={20} strokeWidth={2.1} aria-hidden="true" /> : <LoaderCircle size={20} className="checking-icon" aria-hidden="true" />}
          <span>{loading ? "Checking your message…" : "Check this message"}</span>
        </span>
      </Button>

      <style jsx>{`
        .check-btn-wrapper {
          width: 100%;
        }

        .check-btn-inner {
          display: inline-flex;
          align-items: center;
          gap: 12px;
          font-weight: 700;
        }

        .checking-icon { animation: spin 1s linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
        @media (prefers-reduced-motion: reduce) { .checking-icon { animation: none; } }
      `}</style>
    </div>
  );
}
