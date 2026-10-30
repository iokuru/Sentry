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
import { PracticeScenarios, sampleScenarios } from "@/components/practice-scenarios";
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
  Upload,
  ClipboardPaste,
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
    const params = new URLSearchParams(window.location.search);
    const requestedMode = params.get("mode");
    if (requestedMode === "image" || requestedMode === "text" || requestedMode === "practice") {
      setInputMode(requestedMode);
    }
    const requestedScenario = params.get("scenario") || params.get("example");
    if (requestedScenario) {
      const match = sampleScenarios.find((s) => s.id === requestedScenario);
      const textToRun = match ? match.message : requestedScenario;
      setPastedMessage(textToRun);
      setInputMode("text");
      void executeCheck(textToRun);
    }
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

  const executeCheck = async (textOverride?: string) => {
    const activeFile = inputMode === "image" && !textOverride ? selectedFile : null;
    const activeText = textOverride !== undefined
      ? textOverride.trim()
      : (inputMode === "text" ? pastedMessage.trim() : contextText.trim());
    if (!activeFile && !activeText) return;

    setErrorMessage(null);

    try {
      let objectKey: string | undefined;
      let preparedBlob: Blob | undefined;

      if (activeFile) {
        setFlowState("uploading");
        const { blob } = await downscaleImage(activeFile);
        preparedBlob = blob;
        const presign = await requestPresignedUrl(blob.type || "image/jpeg", activeFile.name);
        await uploadImageToS3(presign.upload_url, blob);
        objectKey = presign.image_key;
      }

      setFlowState("analyzing");

      let intakeRes;
      try {
        intakeRes = await submitIntake({
          image_key: objectKey,
          user_context: activeText || undefined,
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
            image_key: objectKey,
            user_context: [extractedText, activeText].filter(Boolean).join("\n") || undefined,
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

  const handleCheck = () => {
    void executeCheck();
  };

  const handleSelectScenario = (messageText: string) => {
    handleClearFile();
    setPastedMessage(messageText);
    setInputMode("text");
    window.scrollTo({ top: 0, behavior: "smooth" });
    void executeCheck(messageText);
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
          <div className="workbench">
            <header className="page-intro">
              <div>
                <p className="page-kicker">MESSAGE REVIEW</p>
                <h1>Check before you reply</h1>
                <p className="page-description">Bring the message here. We’ll help you decide what to do next.</p>
              </div>
              <div className="region-control"><span>Hotline region</span><CountryPicker selectedCountry={selectedCountry} onCountryChange={setSelectedCountry} /></div>
            </header>

            <section className="gallery-section" aria-labelledby="start-heading">
              <div className="section-heading"><div><h2 id="start-heading">Start with a message</h2><p>Choose the format you have on hand.</p></div><span className="section-count">01 — 03</span></div>
              <div className="mode-gallery" role="tablist" aria-label="Choose a message format">
                <button type="button" role="tab" aria-selected={inputMode === "image"} className={`mode-card ${inputMode === "image" ? "mode-selected" : ""}`} onClick={() => setInputMode("image")}>
                  <span className="mode-card-head"><span className="mode-icon"><ImageIcon size={17} /></span><span className="mode-index">01</span></span>
                  <span className="mode-preview upload-preview"><span className="preview-file"><Upload size={18} /><span>Screenshot or photo</span></span><span className="preview-action">Choose image <span>↗</span></span></span>
                  <span className="mode-card-title">Screenshot</span><span className="mode-card-copy">Upload a photo or paste one from your clipboard.</span>
                </button>
                <button type="button" role="tab" aria-selected={inputMode === "text"} className={`mode-card ${inputMode === "text" ? "mode-selected" : ""}`} onClick={() => setInputMode("text")}>
                  <span className="mode-card-head"><span className="mode-icon"><AlignLeft size={17} /></span><span className="mode-index">02</span></span>
                  <span className="mode-preview text-preview"><span className="preview-lines"><i /><i /><i /></span><span className="preview-cursor" /></span>
                  <span className="mode-card-title">Message text</span><span className="mode-card-copy">Paste the words you received, links included.</span>
                </button>
                <button type="button" role="tab" aria-selected={inputMode === "practice"} className={`mode-card ${inputMode === "practice" ? "mode-selected" : ""}`} onClick={() => setInputMode("practice")}>
                  <span className="mode-card-head"><span className="mode-icon"><ClipboardPaste size={17} /></span><span className="mode-index">03</span></span>
                  <span className="mode-preview sample-preview"><span className="sample-bubble">Your parcel is on hold. Pay a small redelivery fee now.</span><span className="sample-time">11:42 AM</span></span>
                  <span className="mode-card-title">Try a sample</span><span className="mode-card-copy">Walk through a practice message first.</span>
                </button>
              </div>
            </section>

            <section className="active-workspace" aria-labelledby="workspace-title">
              <div className="active-workspace-heading"><div><span className="workspace-index">{inputMode === "image" ? "01" : inputMode === "text" ? "02" : "03"}</span><div><h2 id="workspace-title">{inputMode === "image" ? "Add a screenshot" : inputMode === "text" ? "Paste the message" : "Practice with a sample"}</h2><p>{inputMode === "image" ? "Your image is processed for this check only." : inputMode === "text" ? "Include the sender or any unusual details if useful." : "Select a message to load it into the checker."}</p></div></div></div>
              <div className="workspace-body">
                {inputMode === "image" && <div className="input-pane">{!selectedFile ? <Dropzone onFileSelect={handleFileSelect} onTextPaste={handleTextPaste} /> : previewUrl && <PreviewCard file={selectedFile} previewUrl={previewUrl} onClear={handleClearFile} />}<ContextNote mode="context" value={contextText} onChange={setContextText} /></div>}
                {inputMode === "text" && <div className="input-pane"><ContextNote mode="paste" value={pastedMessage} onChange={setPastedMessage} /></div>}
                {inputMode === "practice" && <div className="input-pane"><PracticeScenarios onSelectScenario={handleSelectScenario} /></div>}
                {inputMode !== "practice" && <div className="workspace-submit"><CheckButton onClick={handleCheck} hasInput={hasInput} loading={false} /></div>}
              </div>
            </section>
            <div className="privacy-line"><span className="privacy-dot" /> No account needed <span className="privacy-separator">·</span> Your screenshot is deleted after its check</div>
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

        .workbench { display:grid; gap:35px; max-width:1250px; margin:4px auto 0; }
        .page-intro { display:flex; justify-content:space-between; align-items:end; gap:24px; padding:0 1px 25px; border-bottom:1px solid var(--border-light); }
        .page-kicker { margin:0 0 10px; color:var(--text-subtle); font:600 12px/1 var(--font-mono); letter-spacing:.12em; }
        .page-intro h1 { margin:0; color:var(--text-main); font-size:clamp(30px,3.5vw,43px); line-height:1.12; letter-spacing:-.045em; font-weight:650; }
        .page-description { margin:10px 0 0; color:var(--text-muted); font-size:16px; }
        .region-control { display:grid; gap:8px; color:var(--text-subtle); font-size:13px; }
        .gallery-section { display:grid; gap:17px; }
        .section-heading { display:flex; align-items:end; justify-content:space-between; gap:18px; }
        .section-heading h2 { margin:0; color:var(--text-main); font-size:21px; font-weight:620; letter-spacing:-.025em; }
        .section-heading p { margin:5px 0 0; color:var(--text-muted); font-size:14px; }
        .section-count { color:var(--text-subtle); font:12px var(--font-mono); }
        .mode-gallery { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:16px; }
        .mode-card { display:grid; grid-template-rows:auto 132px auto auto; gap:13px; padding:15px; text-align:left; color:var(--text-main); background:var(--gallery-card); border:1px solid var(--border-light); border-radius:var(--radius-md); cursor:pointer; transition:border-color .16s,background .16s,transform .16s; }
        .mode-card:hover { border-color:var(--border-medium); transform:translateY(-2px); }
        .mode-selected { border-color:var(--color-primary); box-shadow:inset 0 0 0 1px var(--color-primary); }
        .mode-card-head { display:flex; align-items:center; justify-content:space-between; }
        .mode-icon { display:grid; place-items:center; width:30px; height:30px; border:1px solid var(--border-light); border-radius:7px; color:var(--text-muted); }
        .mode-selected .mode-icon { color:var(--color-primary); border-color:var(--color-primary); }
        .mode-index { color:var(--text-subtle); font:11px var(--font-mono); }
        .mode-preview { position:relative; overflow:hidden; display:flex; align-items:center; justify-content:center; min-width:0; border-radius:6px; background:var(--preview-bg); color:var(--preview-ink); }
        .upload-preview { justify-content:space-between; flex-direction:column; padding:14px; background:var(--preview-bg); border:1px dashed var(--preview-border); }
        .preview-file { display:flex; align-items:center; gap:8px; font-size:13px; }
        .preview-action { display:flex; align-items:center; justify-content:space-between; width:100%; padding:8px 10px; border-radius:5px; background:var(--preview-action-bg); color:var(--preview-action-ink); font-size:12px; }
        .text-preview { justify-content:flex-start; align-items:flex-start; padding:18px; border:1px solid var(--preview-border); }
        .preview-lines { display:grid; gap:9px; width:82%; }
        .preview-lines i { display:block; width:100%; height:5px; border-radius:5px; background:var(--preview-line); }
        .preview-lines i:nth-child(2) { width:72%; }.preview-lines i:nth-child(3) { width:48%; }
        .preview-cursor { position:absolute; left:calc(18px + 48%); top:56px; height:14px; width:1px; background:var(--color-primary); }
        .sample-preview { flex-direction:column; align-items:flex-end; justify-content:center; gap:8px; padding:16px; background:var(--preview-bg); }
        .sample-bubble { max-width:91%; padding:10px 12px; border-radius:10px 10px 2px 10px; background:var(--sample-bubble); color:var(--sample-ink); font-size:12px; line-height:1.45; }
        .sample-time { color:var(--text-subtle); font-size:10px; }
        .mode-card-title { font-size:16px; font-weight:620; letter-spacing:-.015em; }
        .mode-card-copy { color:var(--text-muted); font-size:13px; line-height:1.45; }
        .active-workspace { overflow:hidden; background:var(--work-card); border:1px solid var(--border-light); border-radius:var(--radius-md); }
        .active-workspace-heading { display:flex; align-items:center; justify-content:space-between; gap:18px; padding:19px 21px; border-bottom:1px solid var(--border-light); }
        .active-workspace-heading > div { display:flex; align-items:center; gap:13px; }
        .workspace-index { display:grid; place-items:center; width:34px; height:34px; border:1px solid var(--border-light); border-radius:7px; color:var(--color-primary); font:12px var(--font-mono); }
        .active-workspace-heading h2 { margin:0; font-size:17px; font-weight:620; }
        .active-workspace-heading p { margin:4px 0 0; color:var(--text-muted); font-size:13px; }
        
        .privacy-dot { width:7px; height:7px; border-radius:50%; background:#73b54a; }
        .workspace-body { display:flex; flex-direction:column; gap:20px; padding:21px; }
        .input-pane { display:flex; flex-direction:column; gap:17px; }
        .workspace-submit { margin-top:2px; }
        .privacy-line { display:flex; align-items:center; justify-content:center; gap:9px; color:var(--text-muted); font-size:13px; }
        .privacy-separator { color:var(--text-subtle); }

        /* Result View */
        .result-layout {
          max-width: 1000px;
          margin: 20px auto 0;
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

        @media (max-width: 900px) { .mode-gallery { gap:11px; }.mode-card { padding:12px; }.mode-card-copy { font-size:12px; } }
        @media (max-width: 800px) { .workbench { gap:27px; }.page-intro { align-items:flex-start; flex-direction:column; }.region-control { width:100%; grid-template-columns:1fr auto; align-items:center; }.mode-gallery { grid-template-columns:1fr; }.mode-card { grid-template-columns:35px minmax(100px,.7fr) 1.4fr; grid-template-rows:auto auto; align-items:center; gap:5px 12px; }.mode-card-head { grid-row:1 / 3; align-self:start; }.mode-preview { grid-row:1 / 3; grid-column:2; height:88px; }.mode-card-title,.mode-card-copy { grid-column:3; }.mode-card-title { align-self:end; }.mode-card-copy { align-self:start; }.active-workspace-heading { align-items:flex-start; }.workspace-body { padding:15px; }.privacy-line { flex-wrap:wrap; font-size:12px; } }
      `}</style>
    </Shell>
  );
}
