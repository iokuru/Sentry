import React from "react";
import { Button } from "@/components/button";
import { ShieldCheck } from "lucide-react";

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
        aria-label={loading ? "Checking screenshot for scams" : "Check this screenshot now"}
      >
        <span className="check-btn-inner">
          {!loading ? <ShieldCheck size={26} strokeWidth={2.4} aria-hidden="true" /> : null}
          <span>{loading ? "Checking your message safely..." : "Check this now"}</span>
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
      `}</style>
    </div>
  );
}
