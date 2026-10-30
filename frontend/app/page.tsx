"use client";

import React, { useState } from "react";
import { Shell } from "@/components/shell";
import { Dropzone } from "@/components/dropzone";
import { PreviewCard } from "@/components/preview";
import { ContextNote } from "@/components/context-note";
import { CheckButton } from "@/components/check-button";
import { downscaleImage } from "@/lib/downscale";
import { requestPresignedUrl, uploadImageToS3, submitIntake } from "@/lib/api";

export default function SeniorPage() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [contextText, setContextText] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleFileSelect = async (file: File) => {
    try {
      setErrorMessage(null);
      const url = URL.createObjectURL(file);
      setSelectedFile(file);
      setPreviewUrl(url);
    } catch {
      setErrorMessage("Could not load image. Please try another photo.");
    }
  };

  const handleClearFile = () => {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    setSelectedFile(null);
    setPreviewUrl(null);
    setErrorMessage(null);
  };

  const handleCheck = async () => {
    if (!selectedFile && !contextText.trim()) return;

    setIsLoading(true);
    setErrorMessage(null);

    try {
      let objectKey: string | undefined;

      if (selectedFile) {
        const { blob } = await downscaleImage(selectedFile);
        const presign = await requestPresignedUrl(blob.type || "image/jpeg");
        await uploadImageToS3(presign.upload_url, blob);
        objectKey = presign.object_key;
      }

      await submitIntake({
        object_key: objectKey,
        context_text: contextText.trim() || undefined,
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Something went wrong. Please try again.";
      setErrorMessage(msg);
    } finally {
      setIsLoading(false);
    }
  };

  const hasInput = Boolean(selectedFile || contextText.trim());

  return (
    <Shell activeSection="senior">
      <div className="senior-flow">
        <header className="senior-header">
          <h1 className="senior-title">Check a message or screenshot</h1>
          <p className="senior-subtitle">
            If something feels off, drop an image or paste text below. We will help you figure out what to do.
          </p>
        </header>

        <div className="senior-workspace">
          {!selectedFile ? (
            <Dropzone onFileSelect={handleFileSelect} disabled={isLoading} />
          ) : (
            previewUrl && (
              <PreviewCard
                file={selectedFile}
                previewUrl={previewUrl}
                onClear={handleClearFile}
                disabled={isLoading}
              />
            )
          )}

          <ContextNote
            value={contextText}
            onChange={setContextText}
            disabled={isLoading}
          />

          {errorMessage ? (
            <div className="error-banner" role="alert">
              {errorMessage}
            </div>
          ) : null}

          <CheckButton
            onClick={handleCheck}
            loading={isLoading}
            disabled={isLoading}
            hasInput={hasInput}
          />
        </div>
      </div>

      <style jsx>{`
        .senior-flow {
          display: flex;
          flex-direction: column;
          gap: 28px;
        }

        .senior-header {
          text-align: center;
          display: flex;
          flex-direction: column;
          gap: 12px;
          padding: 8px 0;
        }

        .senior-title {
          font-size: 2.2rem;
          font-weight: 800;
          color: var(--text-main);
          letter-spacing: -0.02em;
        }

        .senior-subtitle {
          font-size: 1.15rem;
          color: var(--text-muted);
          max-width: 580px;
          margin: 0 auto;
          line-height: 1.5;
        }

        .senior-workspace {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .error-banner {
          background-color: var(--color-danger-subtle);
          border: 1.5px solid var(--color-danger);
          color: var(--color-danger);
          padding: 14px 18px;
          border-radius: var(--radius-md);
          font-weight: 600;
          font-size: 1rem;
        }

        @media (max-width: 640px) {
          .senior-title {
            font-size: 1.8rem;
          }
          .senior-subtitle {
            font-size: 1.05rem;
          }
        }
      `}</style>
    </Shell>
  );
}
