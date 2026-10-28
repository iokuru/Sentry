<script lang="ts">
  interface Props {
    busy?: boolean;
    errorMessage?: string;
    onsignin: (email: string, pass: string) => void;
  }
  let { busy = false, errorMessage = '', onsignin }: Props = $props();

  let email = $state('');
  let password = $state('');

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!email || !password || busy) return;
    onsignin(email, password);
  }
</script>

<div class="auth-card">
  <h1 class="title">Caregiver Sign In</h1>

  {#if errorMessage}
    <div class="error-banner" role="alert">
      {errorMessage}
    </div>
  {/if}

  <form onsubmit={handleSubmit} class="form">
    <div class="field">
      <label for="caregiver-email" class="label">Email address</label>
      <input
        id="caregiver-email"
        type="email"
        bind:value={email}
        required
        autocomplete="email"
        class="input"
      />
    </div>

    <div class="field">
      <label for="caregiver-password" class="label">Password</label>
      <input
        id="caregiver-password"
        type="password"
        bind:value={password}
        required
        autocomplete="current-password"
        class="input"
      />
    </div>

    <button type="submit" class="submit-btn" disabled={busy || !email || !password}>
      {busy ? 'Signing in...' : 'Sign In'}
    </button>
  </form>
</div>

<style>
  .auth-card {
    display: flex;
    flex-direction: column;
    gap: 24px;
    padding: 32px 24px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .title {
    font-size: var(--font-xl);
    font-weight: 800;
    color: var(--text-primary);
  }

  .error-banner {
    padding: 12px 16px;
    background-color: var(--rose-50);
    border: 2px solid var(--rose-600);
    color: var(--rose-800);
    font-size: var(--font-base);
    font-weight: 700;
    border-radius: var(--radius-sm);
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .label {
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
  }

  .input {
    min-height: var(--touch-min);
    padding: 12px 16px;
    font-size: var(--font-base);
    background-color: var(--bg-app);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
  }

  .input:focus-visible {
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .submit-btn {
    min-height: var(--touch-min);
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-inverse);
    background-color: var(--stone-900);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .submit-btn:hover:not(:disabled),
  .submit-btn:focus-visible:not(:disabled) {
    background-color: var(--stone-700);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
