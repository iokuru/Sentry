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
import { BeforeYouReplyBadge, ReassuranceStrip } from "@/components/badges";
import { PracticeScenarios } from "@/components/practice-scenarios";
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
import {
  RotateCcw,
  Image as ImageIcon,
  AlignLeft,
  Sparkles,
  ShieldCheck,
  EyeOff,
  BellRing,
  PhoneForwarded,
} from "lucide-react";

type FlowState = "idle" | "uploading" | "extracting" | "analyzing" | "result" | "quota" | "error";

export default function SeniorPage() {
  const [flowState, setFlowState] = useState<FlowState>("idle");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [contextText, setContextText] = useState("");
  const [pastedMessage, setPastedMessage] = useState("");
  const [inputMode, setInputMode] = useState<"image" | "text" | "practice">("image");
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
    setPastedMessage("");
    setInputMode("image");
    setResult(null);
    setErrorMessage(null);
    setFlowState("idle");
  };

  const handleTextPaste = (text: string) => {
    setPastedMessage((prev) => (prev ? `${prev}\n${text}` : text));
    setInputMode("text");
  };

  const handleSelectScenario = (messageText: string) => {
    handleClearFile();
    setPastedMessage(messageText);
    setInputMode("text");
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleCheck = async () => {
    const activeFile = inputMode === "image" ? selectedFile : null;
    const activeText = inputMode === "text" ? pastedMessage.trim() : contextText.trim();
    if (!activeFile && !activeText) return;

    setErrorMessage(null);

    try {
      let objectKey: string | undefined;
      let preparedBlob: Blob | undefined;

      if (activeFile) {
        setFlowState("uploading");
        const { blob } = await downscaleImage(activeFile);
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
          context_text: activeText || undefined,
          country: selectedCountry,
        });
      } catch (err: unknown) {
        if (
          err instanceof ApiError &&
          (err.status === 422 || err.code === "TEXTRACT_UNAVAILABLE") &&
          preparedBlob
        ) {
          setFlowState("extracting");
          const extractedText = await extractTextInBrowser(preparedBlob);
          setFlowState("analyzing");
          intakeRes = await submitIntake({
            object_key: objectKey,
            ocr_text: extractedText,
            context_text: activeText || undefined,
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
      window.scrollTo({ top: 0, behavior: "smooth" });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Unable to complete safety check";
      setErrorMessage(msg);
      setFlowState("error");
    }
  };

  const hasInput =
    inputMode === "image"
      ? Boolean(selectedFile || contextText.trim())
      : Boolean(pastedMessage.trim());

  return (
    <Shell activeSection="senior">
      <div className={`page-wrapper ${flowState !== "idle" ? "page-active" : ""}`}>
        <OfflineBanner />

        {flowState === "idle" && (
          <div className="landing-layout">
            {/* Hero Section */}
            <header className="hero-section">
              <div className="hero-badge-wrap">
                <BeforeYouReplyBadge />
              </div>
              <h1 className="hero-title">A little help before you reply.</h1>
              <p className="hero-subtitle">
                Received a suspicious message, text, or urgent notice? We’ll verify if it’s safe, point out what to pause on, and connect you with official help.
              </p>
              <div className="hero-reassurance">
                <ReassuranceStrip />
              </div>
            </header>

            {/* Interactive Workspace Card */}
            <div className="workspace-card">
              <div className="workspace-header">
                <div className="workspace-tabs" role="tablist" aria-label="Input options">
                  <button
                    type="button"
                    role="tab"
                    aria-selected={inputMode === "image"}
                    className={`tab-btn ${inputMode === "image" ? "tab-btn-active" : ""}`}
                    onClick={() => setInputMode("image")}
                  >
                    <ImageIcon size={18} />
                    <span>Screenshot or photo</span>
                  </button>
                  <button
                    type="button"
                    role="tab"
                    aria-selected={inputMode === "text"}
                    className={`tab-btn ${inputMode === "text" ? "tab-btn-active" : ""}`}
                    onClick={() => setInputMode("text")}
                  >
                    <AlignLeft size={18} />
                    <span>Paste message text</span>
                  </button>
                  <button
                    type="button"
                    role="tab"
                    aria-selected={inputMode === "practice"}
                    className={`tab-btn ${inputMode === "practice" ? "tab-btn-active" : ""}`}
                    onClick={() => setInputMode("practice")}
                  >
                    <Sparkles size={18} />
                    <span>Sample cases</span>
                  </button>
                </div>

                <div className="country-holder">
                  <CountryPicker
                    selectedCountry={selectedCountry}
                    onCountryChange={setSelectedCountry}
                  />
                </div>
              </div>

              <div className="workspace-body">
                {inputMode === "image" && (
                  <div className="input-pane">
                    {!selectedFile ? (
                      <Dropzone
                        onFileSelect={handleFileSelect}
                        onTextPaste={handleTextPaste}
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
                      mode="context"
                      value={contextText}
                      onChange={setContextText}
                    />
                  </div>
                )}

                {inputMode === "text" && (
                  <div className="input-pane">
                    <ContextNote
                      mode="paste"
                      value={pastedMessage}
                      onChange={setPastedMessage}
                    />
                  </div>
                )}

                {inputMode === "practice" && (
                  <div className="input-pane">
                    <PracticeScenarios onSelectScenario={handleSelectScenario} />
                  </div>
                )}

                {inputMode !== "practice" && (
                  <div className="workspace-submit">
                    <CheckButton
                      onClick={handleCheck}
                      hasInput={hasInput}
                      loading={false}
                    />
                  </div>
                )}
              </div>
            </div>

            {/* Product Bento Grid: What Sentry Delivers */}
            <section className="bento-section" aria-labelledby="bento-heading">
              <div className="bento-header">
                <span className="bento-eyebrow">How Sentry Protects You</span>
                <h2 id="bento-heading" className="bento-title">
                  Three clear answers. Zero guesswork.
                </h2>
              </div>

              <div className="bento-grid">
                <div className="bento-card bento-highlight">
                  <div className="bento-icon-pill" style={{ color: "var(--color-accent-gold)" }}>
                    <ShieldCheck size={22} />
                  </div>
                  <h3>1. What to hold off on</h3>
                  <p>
                    Clear, calm advice on what never to click, sign, or send. We identify fake login links, remote access tools, and urgency traps.
                  </p>
                </div>

                <div className="bento-card">
                  <div className="bento-icon-pill" style={{ color: "var(--color-primary)" }}>
                    <PhoneForwarded size={22} />
                  </div>
                  <h3>2. Safe alternative action</h3>
                  <p>
                    A calm, verifiable step you can take right now, such as turning over your card to call the number on the back.
                  </p>
                </div>

                <div className="bento-card">
                  <div className="bento-icon-pill" style={{ color: "#3b82f6" }}>
                    <BellRing size={22} />
                  </div>
                  <h3>3. One-tap fraud hotline</h3>
                  <p>
                    Direct connection to official federal and national fraud authorities (US FTC, UK Action Fraud, India 1930, Australia Scamwatch).
                  </p>
                </div>

                <div className="bento-card bento-privacy">
                  <div className="bento-icon-pill" style={{ color: "#10b981" }}>
                    <EyeOff size={22} />
                  </div>
                  <h3>Private by design</h3>
                  <p>
                    Zero account signup required. Screenshots are downscaled in your browser and automatically expired after your check.
                  </p>
                </div>
              </div>
            </section>
          </div>
        )}

        {(flowState === "uploading" ||
          flowState === "extracting" ||
          flowState === "analyzing") && (
          <ProgressCard
            phase={flowState === "uploading" ? "uploading" : "analyzing"}
          />
        )}

        {flowState === "quota" && <QuotaCard country={selectedCountry} />}

        {flowState === "error" && (
          <FailureCard
            errorMessage={errorMessage || undefined}
            onRetry={handleCheck}
          />
        )}

        {flowState === "result" && result && (
          <div className="result-layout">
            <header className="result-header">
              <span className="result-eyebrow">Safety Assessment</span>
              <h1 className="result-title">Here’s what to do next.</h1>
              <p className="result-subtitle">
                Take a breath. Review these steps before responding or giving any personal information.
              </p>
            </header>

            <div className="result-cards">
              <HoldOffCard advice={result.hold_off} />
              <AlternativeCard action={result.safe_alternative} />
              <HotlineButton
                country={result.country || selectedCountry}
                hotlineName={result.hotline_name}
                phoneNumber={result.hotline_phone}
                officialUrl={result.hotline_url}
              />
            </div>

            <div className="result-reset">
              <Button
                variant="secondary"
                size="large"
                onClick={handleReset}
                fullWidth
              >
                <RotateCcw size={18} aria-hidden="true" />
                <span>Check another message</span>
              </Button>
            </div>
          </div>
        )}
      </div>

      <style jsx>{`
        .page-wrapper {
          color: var(--text-main);
          width: 100%;
        }

        .landing-layout {
          display: flex;
          flex-direction: column;
          gap: 48px;
        }

        /* Hero */
        .hero-section {
          text-align: center;
          max-width: 820px;
          margin: 12px auto 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 16px;
        }

        .hero-badge-wrap {
          margin-bottom: 4px;
        }

        .hero-title {
          font-size: clamp(2.4rem, 4.5vw, 3.8rem);
          font-weight: 850;
          line-height: 1.08;
          letter-spacing: -0.04em;
          color: var(--text-main);
        }

        .hero-subtitle {
          font-size: clamp(1.05rem, 2vw, 1.25rem);
          line-height: 1.55;
          color: var(--text-muted);
          max-width: 680px;
        }

        .hero-reassurance {
          margin-top: 8px;
        }

        /* Workspace Card */
        .workspace-card {
          background: var(--bg-card);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-xl);
          box-shadow: var(--shadow-lg);
          padding: 24px;
          max-width: 840px;
          margin: 0 auto;
          width: 100%;
          position: relative;
          backdrop-filter: blur(12px);
          transition: border-color 200ms ease;
        }

        .workspace-card:hover {
          border-color: var(--border-medium);
        }

        .workspace-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-wrap: wrap;
          gap: 16px;
          padding-bottom: 20px;
          border-bottom: 1px solid var(--border-light);
        }

        .workspace-tabs {
          display: inline-flex;
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-full);
          padding: 4px;
          gap: 4px;
        }

        .tab-btn {
          display: inline-flex;
          align-items: center;
          gap: 7px;
          padding: 8px 14px;
          border-radius: var(--radius-full);
          background: transparent;
          border: none;
          color: var(--text-muted);
          font-size: 0.9rem;
          font-weight: 600;
          cursor: pointer;
          transition: color 140ms ease, background-color 140ms ease;
        }

        .tab-btn:hover {
          color: var(--text-main);
        }

        .tab-btn-active {
          background: var(--bg-card);
          color: var(--text-main);
          box-shadow: var(--shadow-sm);
        }

        .country-holder {
          display: flex;
          align-items: center;
        }

        .workspace-body {
          padding-top: 24px;
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .input-pane {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .workspace-submit {
          margin-top: 6px;
        }

        /* Bento Grid */
        .bento-section {
          margin-top: 48px;
          padding-top: 40px;
          border-top: 1px solid var(--border-light);
        }

        .bento-header {
          text-align: center;
          margin-bottom: 32px;
        }

        .bento-eyebrow {
          font-size: 0.82rem;
          font-weight: 750;
          text-transform: uppercase;
          letter-spacing: 0.12em;
          color: var(--color-primary);
        }

        .bento-title {
          font-size: clamp(1.8rem, 3vw, 2.4rem);
          font-weight: 800;
          letter-spacing: -0.03em;
          color: var(--text-main);
          margin-top: 6px;
        }

        .bento-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
          gap: 20px;
        }

        .bento-card {
          background: var(--bg-card);
          border: 1px solid var(--border-light);
          border-radius: var(--radius-lg);
          padding: 28px 24px;
          display: flex;
          flex-direction: column;
          gap: 12px;
          transition: transform 160ms ease, border-color 160ms ease, box-shadow 160ms ease;
        }

        .bento-card:hover {
          transform: translateY(-2px);
          border-color: var(--border-medium);
          box-shadow: var(--shadow-md);
        }

        .bento-icon-pill {
          width: 44px;
          height: 44px;
          border-radius: var(--radius-md);
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          display: inline-flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 4px;
        }

        .bento-card h3 {
          font-size: 1.15rem;
          font-weight: 750;
          letter-spacing: -0.02em;
          color: var(--text-main);
        }

        .bento-card p {
          font-size: 0.94rem;
          line-height: 1.55;
          color: var(--text-muted);
        }

        /* Result View */
        .result-layout {
          max-width: 820px;
          margin: 0 auto;
          display: flex;
          flex-direction: column;
          gap: 24px;
        }

        .result-header {
          text-align: left;
        }

        .result-eyebrow {
          font-size: 0.8rem;
          font-weight: 750;
          text-transform: uppercase;
          letter-spacing: 0.12em;
          color: var(--color-primary);
        }

        .result-title {
          font-size: clamp(2.2rem, 4vw, 3.2rem);
          font-weight: 850;
          letter-spacing: -0.04em;
          color: var(--text-main);
          margin-top: 4px;
        }

        .result-subtitle {
          font-size: 1.1rem;
          color: var(--text-muted);
          margin-top: 6px;
        }

        .result-cards {
          display: flex;
          flex-direction: column;
          gap: 20px;
        }

        .result-reset {
          margin-top: 12px;
        }

        .page-active {
          min-height: 70vh;
        }

        @media (max-width: 768px) {
          .workspace-header {
            flex-direction: column;
            align-items: stretch;
          }
          .workspace-tabs {
            width: 100%;
            justify-content: space-between;
          }
          .tab-btn {
            flex: 1;
            justify-content: center;
            font-size: 0.82rem;
            padding: 8px 6px;
          }
          .country-holder {
            width: 100%;
            justify-content: flex-end;
          }
          .bento-grid {
            grid-template-columns: 1fr;
          }
        }
      `}</style>
    </Shell>
  );
}
