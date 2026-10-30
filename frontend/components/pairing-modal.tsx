import React, { useState } from "react";
import { Link2, CheckCircle2 } from "lucide-react";
import { Button } from "@/components/button";
import { InputField } from "@/components/field";
import { redeemPairingCode } from "@/lib/device";

export interface PairingModalProps {
  onClose?: () => void;
}

export function PairingModal({ onClose }: PairingModalProps) {
  const [code, setCode] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [status, setStatus] = useState<"idle" | "success" | "error">("idle");
  const [errorMessage, setErrorMessage] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!code.trim()) return;

    setIsSubmitting(true);
    setStatus("idle");
    setErrorMessage("");

    const res = await redeemPairingCode(code);
    setIsSubmitting(false);

    if (res.success) {
      setStatus("success");
      setTimeout(() => {
        if (onClose) onClose();
      }, 2000);
    } else {
      setStatus("error");
      setErrorMessage(res.message || "Invalid or expired invite code.");
    }
  };

  return (
    <div className="pairing-card">
      <div className="pairing-header">
        <span className="pairing-icon" aria-hidden="true">
          <Link2 size={24} strokeWidth={2.4} />
        </span>
        <div>
          <h2 className="pairing-title">Link with a trusted caregiver</h2>
          <p className="pairing-subtitle">
            Enter the six-letter invite code from your family member or caregiver.
          </p>
        </div>
      </div>

      {status === "success" ? (
        <div className="pairing-success">
          <CheckCircle2 size={32} className="success-icon" />
          <p className="success-text">Successfully linked! Your caregiver will be alerted if high risk is detected.</p>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="pairing-form">
          <InputField
            label="Caregiver invite code"
            placeholder="e.g. ABCDEF"
            value={code}
            onChange={(e) => setCode(e.target.value.toUpperCase())}
            maxLength={8}
            disabled={isSubmitting}
            error={errorMessage || undefined}
          />

          <Button
            type="submit"
            variant="primary"
            fullWidth
            loading={isSubmitting}
            disabled={isSubmitting || !code.trim()}
          >
            Connect caregiver
          </Button>
        </form>
      )}

      <style jsx>{`
        .pairing-card {
          background-color: var(--bg-card);
          border: 1.5px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 24px;
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .pairing-header {
          display: flex;
          align-items: flex-start;
          gap: 14px;
        }

        .pairing-icon {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 44px;
          height: 44px;
          border-radius: var(--radius-md);
          background-color: var(--bg-secondary);
          color: var(--color-primary);
          flex-shrink: 0;
        }

        .pairing-title {
          font-size: 1.25rem;
          font-weight: 700;
          color: var(--text-main);
          margin-bottom: 4px;
        }

        .pairing-subtitle {
          font-size: 0.98rem;
          color: var(--text-muted);
          line-height: 1.4;
        }

        .pairing-form {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .pairing-success {
          display: flex;
          align-items: center;
          gap: 12px;
          background-color: var(--color-safe-subtle);
          border: 1px solid var(--color-safe);
          padding: 16px;
          border-radius: var(--radius-md);
        }

        .success-text {
          font-size: 1.05rem;
          font-weight: 600;
          color: var(--color-safe);
        }
      `}</style>
    </div>
  );
}
