<script lang="ts">
  let { onsignin } = $props<{
    onsignin: (email: string, pass: string) => void;
  }>();

  let email = $state('');
  let password = $state('');
  let busy = $state(false);

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!email || !password) return;
    busy = true;
    onsignin(email, password);
  }
</script>

<div class="auth-card">
  <h1 class="title">Caregiver Sign In</h1>
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
      Sign In
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
  }

  .submit-btn {
    min-height: var(--touch-min);
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-inverse);
    background-color: var(--text-primary);
    border: 2px solid var(--text-primary);
    border-radius: var(--radius-sm);
    padding: 14px 20px;
    cursor: pointer;
  }

  .submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
