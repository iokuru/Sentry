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
  import PairingCard from '$lib/PairingCard.svelte';
  import { detectCountry } from '$lib/locale';
  import { downscaleImage } from '$lib/downscale';
  import { isPaired, getSeniorName } from '$lib/device';
  import {
    presignUpload,
    uploadScreenshot,
    submitIntake,
    pollIncident,
    QuotaExceededError,
    type Incident
  } from '$lib/api';

  let country = $state('US');
  let selectedFile = $state<File | null>(null);
  let processedBlob = $state<Blob | null>(null);
  let previewUrl = $state('');
  let contextText = $state('');
  let busy = $state(false);
  let isUploading = $state(false);
  let isAnalyzing = $state(false);
  let isOffline = $state(false);
  let isFailed = $state(false);
  let isQuotaReached = $state(false);
  let errorMessage = $state('');
  let result = $state<Incident | null>(null);
  let paired = $state(false);
  let seniorName = $state('');
  let showPairing = $state(false);

  onMount(() => {
    country = detectCountry();
    paired = isPaired();
    seniorName = getSeniorName();
    const skipped = localStorage.getItem('sentry_pair_skipped') === 'true';
    if (!paired && !skipped) {
      showPairing = true;
    }

    const savedResult = localStorage.getItem('sentry_last_result');
    if (savedResult) {
      try {
        result = JSON.parse(savedResult);
      } catch {
        localStorage.removeItem('sentry_last_result');
      }
    }

    const updateOnline = () => {
      isOffline = typeof navigator !== 'undefined' && !navigator.onLine;
    };
    window.addEventListener('online', updateOnline);
    window.addEventListener('offline', updateOnline);
    updateOnline();

    const handlePaste = (e: ClipboardEvent) => {
      if (!e.clipboardData) return;
      const files = e.clipboardData.files;
      if (files && files.length > 0) {
        for (let i = 0; i < files.length; i++) {
          if (files[i].type.startsWith('image/')) {
            e.preventDefault();
            handleFileSelected(files[i]);
            return;
          }
        }
      }
      const items = e.clipboardData.items;
      if (items) {
        for (let i = 0; i < items.length; i++) {
          if (items[i].type.startsWith('image/')) {
            const blob = items[i].getAsFile();
            if (blob) {
              e.preventDefault();
              handleFileSelected(blob);
              return;
            }
          }
        }
      }
    };

    window.addEventListener('paste', handlePaste);
    return () => {
      window.removeEventListener('paste', handlePaste);
      window.removeEventListener('online', updateOnline);
      window.removeEventListener('offline', updateOnline);
    };
  });

  function handlePairComplete() {
    paired = isPaired();
    seniorName = getSeniorName();
    showPairing = false;
    localStorage.setItem('sentry_pair_skipped', 'true');
  }

  async function handleFileSelected(file: File) {
    errorMessage = '';

    const isHeic = file.type === 'image/heic' || file.type === 'image/heif' || file.name.toLowerCase().endsWith('.heic') || file.name.toLowerCase().endsWith('.heif');
    if (isHeic) {
      errorMessage = 'Please take a screenshot or use a standard photo format';
      return;
    }

    const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      errorMessage = 'Only PNG, JPG, and WebP images are supported';
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      errorMessage = 'Screenshot exceeds five megabyte limit';
      return;
    }

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

    if (isOffline) {
      errorMessage = 'Connect to the internet to check message';
      return;
    }

    busy = true;
    isFailed = false;
    errorMessage = '';

    try {
      let imageKey = '';

      if (processedBlob) {
        isUploading = true;
        const contentType = processedBlob.type || 'image/png';
        const filename = selectedFile?.name || 'screenshot.png';

        const presignData = await presignUpload(contentType, filename);
        imageKey = presignData.image_key;

        await uploadScreenshot(presignData.upload_url, processedBlob, contentType);
        isUploading = false;
      }

      isAnalyzing = true;
      const intakeData = await submitIntake({
        image_key: imageKey,
        user_context: contextText.trim(),
        country: country
      });

      result = await pollIncident(intakeData.incident_id);
      if (result) {
        try {
          localStorage.setItem('sentry_last_result', JSON.stringify(result));
        } catch {
          // ignore storage quota error
        }
      }
    } catch (err: unknown) {
      if (err instanceof QuotaExceededError || (err instanceof Error && err.name === 'QuotaExceededError')) {
        isQuotaReached = true;
      } else {
        isFailed = true;
        errorMessage = err instanceof Error ? err.message : 'Analysis failed';
      }
    } finally {
      busy = false;
      isUploading = false;
      isAnalyzing = false;
    }
  }

  function handleReset() {
    handleRemoveFile();
    contextText = '';
    result = null;
    errorMessage = '';
    isFailed = false;
    isQuotaReached = false;
    try {
      localStorage.removeItem('sentry_last_result');
    } catch {
      // quiet fail
    }
  }
</script>

<main class="page">
  <header class="top-bar">
    <span class="brand">Sentry</span>
    <div class="actions">
      {#if paired}
        <span class="pair-badge" title="Linked to caregiver">{seniorName || 'Linked'}</span>
      {:else if !showPairing}
        <button type="button" class="link-btn" onclick={() => showPairing = true}>
          Link
        </button>
      {/if}
      <CountrySwitch bind:country />
      <a href="/caregiver" class="portal-link">Caregiver</a>
    </div>
  </header>

  {#if showPairing}
    <PairingCard oncomplete={handlePairComplete} />
  {:else if !result || !result.output}
    <section class="intake-section">
      {#if !previewUrl}
        <Dropzone
          onfile={handleFileSelected}
          ontext={(txt) => { contextText = txt; }}
          hasFile={false}
        />
      {:else}
        <ImagePreview src={previewUrl} onremove={handleRemoveFile} />
      {/if}

      <ContextField bind:value={contextText} />

      {#if isQuotaReached}
        <div class="state-card quota-card" role="alert">
          <h2 class="state-title">Daily Limit Reached</h2>
          <p class="state-desc">The free analysis allowance has been reached for today. Resets at midnight UTC.</p>
          <button type="button" class="retry-btn" onclick={handleReset}>
            Check another message
          </button>
        </div>
      {/if}

      {#if isOffline}
        <div class="state-card offline-card" role="alert">
          <h2 class="state-title">Offline</h2>
          <p class="state-desc">Connect to the internet to check suspicious messages</p>
        </div>
      {/if}

      {#if isFailed && !busy}
        <div class="state-card failed-card" role="alert">
          <h2 class="state-title">Check Failed</h2>
          <p class="state-desc">{errorMessage || 'Analysis failed, please try again'}</p>
          <button type="button" class="retry-btn" onclick={handleCheck}>
            Try again
          </button>
        </div>
      {:else if errorMessage && !isFailed}
        <div class="error-banner" role="alert">
          {errorMessage}
        </div>
      {/if}

      <CheckButton
        {busy}
        busyLabel={isUploading ? 'Uploading screenshot...' : (isAnalyzing ? 'Analyzing message...' : 'Checking...')}
        disabled={!processedBlob && !contextText.trim() || isOffline}
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
    min-height: 100vh;
    padding: 16px 20px 48px;
    max-width: 680px;
    margin: 0 auto;
    gap: 24px;
  }

  .top-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--border-subtle);
  }

  .brand {
    font-size: var(--font-xl);
    font-weight: 800;
    color: var(--text-primary);
    letter-spacing: -0.02em;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .pair-badge {
    display: inline-flex;
    align-items: center;
    padding: 6px 12px;
    font-size: var(--font-sm);
    font-weight: 700;
    color: var(--emerald-800, #065f46);
    background-color: var(--emerald-50, #ecfdf5);
    border: 2px solid var(--emerald-600, #059669);
    border-radius: var(--radius-sm);
  }

  .link-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 14px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .link-btn:hover,
  .link-btn:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .portal-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 16px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    text-decoration: none;
    transition: background-color 150ms ease;
  }

  .portal-link:hover,
  .portal-link:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .intake-section,
  .result-section {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .error-banner {
    padding: 14px 18px;
    background-color: var(--rose-50);
    border: 2px solid var(--rose-600);
    color: var(--rose-800);
    font-size: var(--font-base);
    font-weight: 700;
    border-radius: var(--radius-sm);
  }

  .state-card {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 20px 24px;
    border-radius: var(--radius-md);
    border: 2px solid;
  }

  .offline-card {
    background-color: var(--amber-50, #fffbeb);
    border-color: var(--amber-600, #d97706);
    color: var(--amber-900, #78350f);
  }

  .failed-card {
    background-color: var(--rose-50);
    border-color: var(--rose-600);
    color: var(--rose-900, #881337);
  }

  .quota-card {
    background-color: var(--stone-100, #f5f5f4);
    border-color: var(--stone-500, #78716c);
    color: var(--stone-900, #1c1917);
  }

  .state-title {
    font-size: var(--font-lg);
    font-weight: 800;
    margin: 0;
  }

  .state-desc {
    font-size: var(--font-base);
    font-weight: 600;
    margin: 0;
  }

  .retry-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    align-self: flex-start;
    min-height: var(--touch-min);
    padding: 0 20px;
    margin-top: 8px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .retry-btn:hover,
  .retry-btn:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .reset-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    width: 100%;
    padding: 16px 24px;
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .reset-btn:hover,
  .reset-btn:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }
</style>
