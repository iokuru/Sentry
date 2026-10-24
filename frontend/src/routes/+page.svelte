<script lang="ts">
  import { onMount } from 'svelte';
  import Dropzone from '$lib/Dropzone.svelte';
  import ImagePreview from '$lib/ImagePreview.svelte';
  import ContextField from '$lib/ContextField.svelte';
  import CheckButton from '$lib/CheckButton.svelte';
  import HoldOffCard from '$lib/HoldOffCard.svelte';
  import SafeAlternativeCard from '$lib/SafeAlternativeCard.svelte';
  import HotlineButton from '$lib/HotlineButton.svelte';
  import CountrySwitch from '$lib/CountrySwitch.svelte';
  import { detectCountry } from '$lib/locale';
  import { downscaleImage } from '$lib/downscale';

  let country = $state('US');
  let selectedFile = $state<File | null>(null);
  let processedBlob = $state<Blob | null>(null);
  let previewUrl = $state('');
  let contextText = $state('');
  let busy = $state(false);
  let errorMessage = $state('');

  interface IncidentResult {
    id: string;
    score: number;
    risk_level: string;
    category: string;
    output: {
      hold_off: string;
      safe_alternative: string;
      hotline: {
        country: string;
        name: string;
        number: string;
        tel_uri: string;
      };
    };
  }

  let result = $state<IncidentResult | null>(null);

  onMount(() => {
    country = detectCountry();
  });

  async function handleFileSelected(file: File) {
    errorMessage = '';
    selectedFile = file;
    try {
      const downscaled = await downscaleImage(file);
      processedBlob = downscaled.blob;
      previewUrl = downscaled.previewUrl;
    } catch {
      processedBlob = file;
      previewUrl = URL.createObjectURL(file);
    }
  }

  function handleRemoveFile() {
    selectedFile = null;
    processedBlob = null;
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
      previewUrl = '';
    }
  }

  async function handleCheck() {
    if (!processedBlob && !contextText.trim()) {
      return;
    }

    busy = true;
    errorMessage = '';

    try {
      let imageKey = '';

      if (processedBlob) {
        const presignRes = await fetch('/api/presign', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            content_type: processedBlob.type || 'image/png',
            filename: selectedFile?.name || 'screenshot.png'
          })
        });

        if (!presignRes.ok) {
          throw new Error('Could not prepare upload');
        }

        const presignData = await presignRes.json();
        imageKey = presignData.image_key;

        const uploadRes = await fetch(presignData.upload_url, {
          method: 'PUT',
          headers: { 'Content-Type': processedBlob.type || 'image/png' },
          body: processedBlob
        });

        if (!uploadRes.ok) {
          throw new Error('Could not upload image');
        }
      }

      const intakeRes = await fetch('/api/intake', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          image_key: imageKey,
          user_context: contextText.trim(),
          country: country
        })
      });

      if (!intakeRes.ok) {
        throw new Error('Analysis request failed');
      }

      const intakeData = await intakeRes.json();
      const incidentId = intakeData.incident_id;

      // poll for finished result
      for (let attempt = 0; attempt < 25; attempt++) {
        await new Promise((resolve) => setTimeout(resolve, 1500));
        const res = await fetch(`/api/incidents?id=${incidentId}`);
        if (res.ok) {
          const data: IncidentResult = await res.json();
          if (data && data.output && data.output.hold_off) {
            result = data;
            break;
          }
        }
      }

      if (!result) {
        throw new Error('Analysis timed out. Please try again.');
      }
    } catch (err: unknown) {
      errorMessage = err instanceof Error ? err.message : 'Something went wrong';
    } finally {
      busy = false;
    }
  }

  function handleReset() {
    handleRemoveFile();
    contextText = '';
    result = null;
    errorMessage = '';
  }
</script>

<main class="page">
  <header class="top-bar">
    <span class="brand">Sentry</span>
    <div class="actions">
      <CountrySwitch bind:country />
      <a href="/caregiver" class="portal-link">Caregiver</a>
    </div>
  </header>

  {#if !result}
    <section class="intake-section">
      {#if !previewUrl}
        <Dropzone onfile={handleFileSelected} hasFile={false} />
      {:else}
        <ImagePreview src={previewUrl} onremove={handleRemoveFile} />
      {/if}

      <ContextField bind:value={contextText} />

      {#if errorMessage}
        <div class="error-banner" role="alert">
          {errorMessage}
        </div>
      {/if}

      <CheckButton
        {busy}
        disabled={!processedBlob && !contextText.trim()}
        onclick={handleCheck}
      />
    </section>
  {:else}
    <section class="result-section">
      <HoldOffCard advice={result.output.hold_off} />
      <SafeAlternativeCard action={result.output.safe_alternative} />
      <HotlineButton hotline={result.output.hotline} />

      <button type="button" class="reset-btn" onclick={handleReset}>
        Check another message
      </button>
    </section>
  {/if}
</main>

<style>
  .page {
    display: flex;
    flex-direction: column;
    gap: 32px;
  }

  .top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 16px;
    border-bottom: 2px solid var(--border-subtle);
  }

  .brand {
    font-size: var(--font-xl);
    font-weight: 800;
    letter-spacing: -0.02em;
    color: var(--text-primary);
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .portal-link {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-muted);
    padding: 6px 12px;
    border-radius: var(--radius-sm);
  }

  .portal-link:hover {
    color: var(--text-primary);
  }

  .intake-section {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .result-section {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .error-banner {
    padding: 16px;
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--rose-text);
    background-color: var(--rose-bg);
    border: 2px solid var(--rose-border);
    border-radius: var(--radius-sm);
  }

  .reset-btn {
    min-height: var(--touch-min);
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-secondary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    padding: 14px 20px;
    cursor: pointer;
    text-align: center;
  }

  .reset-btn:hover {
    color: var(--text-primary);
    border-color: var(--border-strong);
  }
</style>
