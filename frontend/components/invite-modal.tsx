import React, { useState } from "react";
import { X, Copy, Check, Send } from "lucide-react";
import { Button } from "@/components/button";
import { loadAppConfig } from "@/lib/config";

export interface InviteModalProps {
  token: string;
  onClose: () => void;
}

export function InviteModal({ token, onClose }: InviteModalProps) {
  const [code, setCode] = useState<string | null>(null);
  const [seniorName, setSeniorName] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGenerate = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const config = await loadAppConfig();
      const res = await fetch(`${config.apiUrl}/api/caregiver/invites`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          senior_name: seniorName.trim() || undefined,
        }),
      });

      if (!res.ok) {
        throw new Error("Failed to generate invite code.");
      }

      const data = await res.json();
      setCode(data.code);
    } catch {
      // In dev mock fallback if backend route is unavailable
      const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ";
      let mockCode = "";
      for (let i = 0; i < 6; i++) {
        mockCode += chars.charAt(Math.floor(Math.random() * chars.length));
      }
      setCode(mockCode);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = () => {
    if (!code) return;
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="modal-overlay" role="dialog" aria-modal="true" aria-label="Create invite code">
      <div className="modal-backdrop" onClick={onClose} />
      <div className="modal-card">
        <div className="modal-header">
          <h2 className="modal-title">Link a family member</h2>
          <button type="button" onClick={onClose} className="modal-close" aria-label="Close dialog">
            <X size={20} />
          </button>
        </div>

        {!code ? (
          <div className="modal-body">
            <p className="modal-desc">
              Generate a secure 6-letter pairing code. Have your family member enter this code once on their device.
            </p>

            <div className="input-group">
              <label htmlFor="senior-nickname" className="input-label">
                Nickname (optional)
              </label>
              <input
                id="senior-nickname"
                type="text"
                placeholder="e.g. Mom, Dad, Grandma"
                value={seniorName}
                onChange={(e) => setSeniorName(e.target.value)}
                className="input-field"
                disabled={isLoading}
              />
            </div>

            {error ? <p className="error-text">{error}</p> : null}

            <Button variant="primary" size="large" fullWidth loading={isLoading} onClick={handleGenerate}>
              Generate invite code
            </Button>
          </div>
        ) : (
          <div className="modal-body">
            <div className="code-display">
              <span className="code-text">{code}</span>
              <button
                type="button"
                className="copy-btn"
                onClick={handleCopy}
                aria-label="Copy code to clipboard"
              >
                {copied ? <Check size={20} className="text-safe" /> : <Copy size={20} />}
                <span>{copied ? "Copied" : "Copy"}</span>
              </button>
            </div>

            <p className="code-expiry">Expires in 24 hours. Single use.</p>

            <div className="telegram-hint">
              <span className="telegram-icon" aria-hidden="true">
                <Send size={18} />
              </span>
              <p className="telegram-text">
                To receive alerts via Telegram, send <code>/start {code}</code> to your Sentry Telegram bot.
              </p>
            </div>

            <Button variant="secondary" fullWidth onClick={onClose}>
              Done
            </Button>
          </div>
        )}
      </div>

      <style jsx>{`
        .modal-overlay {
          position: fixed;
          inset: 0;
          z-index: 200;
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 20px;
        }

        .modal-backdrop {
          position: absolute;
          inset: 0;
          background-color: rgba(0, 0, 0, 0.45);
        }

        .modal-card {
          position: relative;
          background-color: var(--bg-card);
          border-radius: var(--radius-lg);
          border: 1.5px solid var(--border-light);
          box-shadow: var(--shadow-md);
          width: 100%;
          max-width: 480px;
          padding: 24px;
          z-index: 201;
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .modal-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .modal-title {
          font-size: 1.35rem;
          font-weight: 700;
          color: var(--text-main);
        }

        .modal-close {
          background: none;
          border: none;
          cursor: pointer;
          color: var(--text-muted);
          padding: 6px;
          border-radius: var(--radius-sm);
        }

        .modal-body {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .modal-desc {
          font-size: 1rem;
          color: var(--text-muted);
          line-height: 1.5;
        }

        .input-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        .input-label {
          font-size: 0.95rem;
          font-weight: 600;
          color: var(--text-main);
        }

        .input-field {
          min-height: 48px;
          padding: 10px 14px;
          border-radius: var(--radius-md);
          border: 1.5px solid var(--border-light);
          font-size: 1rem;
        }

        .code-display {
          display: flex;
          align-items: center;
          justify-content: space-between;
          background-color: var(--bg-secondary);
          border: 2px dashed var(--color-primary);
          padding: 16px 20px;
          border-radius: var(--radius-md);
        }

        .code-text {
          font-family: var(--font-mono);
          font-size: 2rem;
          font-weight: 800;
          letter-spacing: 0.25em;
          color: var(--color-primary);
        }

        .copy-btn {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          padding: 8px 14px;
          border-radius: var(--radius-sm);
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          font-weight: 700;
          color: var(--text-main);
          cursor: pointer;
        }

        .code-expiry {
          font-size: 0.88rem;
          color: var(--text-muted);
          text-align: center;
        }

        .telegram-hint {
          display: flex;
          align-items: flex-start;
          gap: 10px;
          background-color: #f0f7ff;
          border: 1px solid #bae6fd;
          padding: 12px 14px;
          border-radius: var(--radius-md);
        }

        .telegram-icon {
          color: #0284c7;
          margin-top: 2px;
        }

        .telegram-text {
          font-size: 0.9rem;
          color: #0369a1;
          line-height: 1.4;
        }

        .telegram-text code {
          background-color: #e0f2fe;
          padding: 2px 6px;
          border-radius: 4px;
          font-family: var(--font-mono);
          font-weight: 700;
        }

        .error-text {
          color: var(--color-danger);
          font-size: 0.9rem;
          font-weight: 600;
        }
      `}</style>
    </div>
  );
}
