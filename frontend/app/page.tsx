"use client";

import React, { useState, useEffect } from "react";
import { Shell } from "@/components/shell";
import { Dropzone } from "@/components/dropzone";
import { PreviewCard } from "@/components/preview";
import { ContextNote } from "@/components/context-note";
import { CheckButton } from "@/components/check-button";
import { HoldOffCard } from "@/components/hold-off-card";
import { AlternativeCard } from "@/components/alternative-card";
import { HotlineButton } from "@/components/hotline-button";
import { CountryPicker, detectBrowserCountry } from "@/components/country-picker";
import { ProgressCard } from "@/components/progress-card";
import { OfflineBanner } from "@/components/offline-banner";
import { FailureCard } from "@/components/failure-card";
import { QuotaCard } from "@/components/quota-card";
import { Button } from "@/components/button";
import { downscaleImage } from "@/lib/downscale";
import {
  requestPresignedUrl,
  uploadImageToS3,
  submitIntake,
  pollIncidentResult,
  type IncidentResult,
  ApiError,
} from "@/lib/api";
import { extractTextInBrowser } from "@/lib/ocr";
import { RotateCcw } from "lucide-react";

type FlowState = "idle" | "uploading" | "extracting" | "analyzing" | "result" | "quota" | "error";

export default function SeniorPage() {
  const [flowState, setFlowState] = useState<FlowState>("idle");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [contextText, setContextText] = useState("");
  const [selectedCountry, setSelectedCountry] = useState("US");
  const [result, setResult] = useState<IncidentResult | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    setSelectedCountry(detectBrowserCountry());
  }, []);

  const handleFileSelect = (file: File) => {
    setErrorMessage(null);
    const url = URL.createObjectURL(file);
    setSelectedFile(file);
    setPreviewUrl(url);
  };

  const handleClearFile = () => {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    setSelectedFile(null);
    setPreviewUrl(null);
    setErrorMessage(null);
  };

  const handleReset = () => {
    handleClearFile();
    setContextText("");
    setResult(null);
    setErrorMessage(null);
    setFlowState("idle");
  };

  const handleCheck = async () => {
    if (!selectedFile && !contextText.trim()) return;

    setErrorMessage(null);

    try {
      let objectKey: string | undefined;
      let preparedBlob: Blob | undefined;

      if (selectedFile) {
        setFlowState("uploading");
        const { blob } = await downscaleImage(selectedFile);
        preparedBlob = blob;
        const presign = await requestPresignedUrl(blob.type || "image/jpeg");
        await uploadImageToS3(presign.upload_url, blob);
        objectKey = presign.object_key;
      }

      setFlowState("analyzing");

      let intakeRes;
      try {
        intakeRes = await submitIntake({
          object_key: objectKey,
          context_text: contextText.trim() || undefined,
          country: selectedCountry,
        });
      } catch (err: unknown) {
        if (err instanceof ApiError && (err.status === 422 || err.code === "TEXTRACT_UNAVAILABLE") && preparedBlob) {
          // Fall back to client side browser ocr
          setFlowState("extracting");
          const extractedText = await extractTextInBrowser(preparedBlob);
          setFlowState("analyzing");
          intakeRes = await submitIntake({
            object_key: objectKey,
            ocr_text: extractedText,
            context_text: contextText.trim() || undefined,
            country: selectedCountry,
          });
        } else if (err instanceof ApiError && err.status === 429) {
          setFlowState("quota");
          return;
        } else {
          throw err;
        }
      }

      const completed = await pollIncidentResult(intakeRes.incident_id);
      setResult(completed);
      setFlowState("result");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Unable to complete safety check";
      setErrorMessage(msg);
      setFlowState("error");
    }
  };

  const hasInput = Boolean(selectedFile || contextText.trim());

  return (
    <Shell activeSection="senior">
      <div className="senior-flow">
        <OfflineBanner />

        <header className="senior-header">
          <div className="senior-header-top">
            <h1 className="senior-title">Check a message or screenshot</h1>
            <CountryPicker
              selectedCountry={selectedCountry}
              onCountryChange={setSelectedCountry}
            />
          </div>
          <p className="senior-subtitle">
            If something feels off, drop an image or paste text below. We will help you figure out what to do.
          </p>
        </header>

        {flowState === "idle" && (
          <div className="senior-workspace">
            {!selectedFile ? (
              <Dropzone
                onFileSelect={handleFileSelect}
                onTextPaste={(text) => setContextText((prev) => (prev ? `${prev}\n${text}` : text))}
              />
            ) : (
              previewUrl && (
                <PreviewCard
                  file={selectedFile}
                  previewUrl={previewUrl}
                  onClear={handleClearFile}
                />
              )
            )}

            <ContextNote
              value={contextText}
              onChange={setContextText}
            />

            <CheckButton
              onClick={handleCheck}
              hasInput={hasInput}
            />
          </div>
        )}

        {(flowState === "uploading" || flowState === "extracting" || flowState === "analyzing") && (
          <ProgressCard phase={flowState === "uploading" ? "uploading" : "analyzing"} />
        )}

        {flowState === "quota" && <QuotaCard country={selectedCountry} />}

        {flowState === "error" && (
          <FailureCard
            errorMessage={errorMessage || undefined}
            onRetry={handleCheck}
          />
        )}

        {flowState === "result" && result && (
          <div className="result-container">
            <HoldOffCard advice={result.hold_off} />
            <AlternativeCard action={result.safe_alternative} />
            <HotlineButton
              country={result.country || selectedCountry}
              hotlineName={result.hotline_name}
              phoneNumber={result.hotline_phone}
              officialUrl={result.hotline_url}
            />

            <div className="result-actions">
              <Button variant="secondary" size="large" onClick={handleReset} fullWidth>
                <RotateCcw size={20} aria-hidden="true" />
                <span>Check another screenshot</span>
              </Button>
            </div>
          </div>
        )}
      </div>

      <style jsx>{`
        .senior-flow {
          display: flex;
          flex-direction: column;
          gap: 28px;
        }

        .senior-header {
          display: flex;
          flex-direction: column;
          gap: 12px;
          padding: 8px 0;
        }

        .senior-header-top {
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 16px;
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
          line-height: 1.5;
        }

        .senior-workspace {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .result-container {
          display: flex;
          flex-direction: column;
          gap: 24px;
        }

        .result-actions {
          margin-top: 8px;
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
