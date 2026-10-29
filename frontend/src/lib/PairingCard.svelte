<script lang="ts">
  import { pairDevice } from '$lib/api';
  import { getDeviceId, setPaired } from '$lib/device';

  let { oncomplete }: { oncomplete: () => void } = $props();

  let code = $state('');
  let busy = $state(false);
  let errorMessage = $state('');

  function handleInput(e: Event) {
    const input = e.target as HTMLInputElement;
    code = input.value.toUpperCase().replace(/[^A-Z0-9]/g, '').slice(0, 6);
  }

  async function handlePair() {
    if (code.length < 6) {
      errorMessage = 'Enter full six character code';
      return;
    }

    busy = true;
    errorMessage = '';

    try {
      const deviceId = getDeviceId();
      const res = await pairDevice(code, deviceId);
      setPaired(res.senior_name);
      oncomplete();
    } catch (err: unknown) {
      errorMessage = err instanceof Error ? err.message : 'Connection failed';
    } finally {
      busy = false;
    }
  }

  function handleSkip() {
    oncomplete();
  }
</script>

<div class="pairing-container" role="region" aria-label="Caregiver Connection">
  <div class="pairing-card">
    <h2 class="pairing-title">Connect Caregiver</h2>
    <p class="pairing-desc">Enter the invite code from your caregiver</p>

    <div class="input-wrap">
      <input
        type="text"
        class="code-input"
        value={code}
        oninput={handleInput}
        placeholder="ABCDEF"
        maxlength="6"
        autocomplete="off"
        autocapitalize="characters"
        spellcheck="false"
        aria-label="Six character invite code"
        disabled={busy}
      />
    </div>

    {#if errorMessage}
      <div class="error-banner" role="alert">
        {errorMessage}
      </div>
    {/if}

    <button
      type="button"
      class="connect-btn"
      disabled={code.length < 6 || busy}
      onclick={handlePair}
    >
      {busy ? 'Connecting...' : 'Connect'}
    </button>

    <button
      type="button"
      class="skip-btn"
      onclick={handleSkip}
    >
      Skip for now
    </button>
  </div>
</div>

<style>
  .pairing-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 24px 0;
  }

  .pairing-card {
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 480px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-lg);
    padding: 32px 24px;
    gap: 20px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
  }

  .pairing-title {
    font-size: var(--font-xl);
    font-weight: 800;
    color: var(--text-primary);
    margin: 0;
    text-align: center;
    letter-spacing: -0.01em;
  }

  .pairing-desc {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-secondary);
    margin: 0;
    text-align: center;
  }

  .input-wrap {
    display: flex;
    justify-content: center;
  }

  .code-input {
    width: 100%;
    max-width: 320px;
    min-height: var(--touch-min);
    font-size: 32px;
    font-weight: 800;
    letter-spacing: 0.25em;
    text-align: center;
    text-transform: uppercase;
    color: var(--text-primary);
    background-color: var(--bg-canvas);
    border: 3px solid var(--border-strong);
    border-radius: var(--radius-md);
    padding: 8px 16px;
    transition: border-color 150ms ease;
  }

  .code-input:focus-visible {
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
    border-color: var(--text-primary);
  }

  .error-banner {
    padding: 12px 16px;
    background-color: var(--rose-50);
    border: 2px solid var(--rose-600);
    color: var(--rose-800);
    font-size: var(--font-sm);
    font-weight: 700;
    border-radius: var(--radius-sm);
    text-align: center;
  }

  .connect-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    width: 100%;
    padding: 16px 24px;
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--bg-surface);
    background-color: var(--text-primary);
    border: 2px solid var(--text-primary);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: opacity 150ms ease;
  }

  .connect-btn:hover:not(:disabled),
  .connect-btn:focus-visible:not(:disabled) {
    opacity: 0.9;
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .connect-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .skip-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    width: 100%;
    padding: 12px 16px;
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-secondary);
    background: transparent;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }

  .skip-btn:hover,
  .skip-btn:focus-visible {
    color: var(--text-primary);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }
</style>
