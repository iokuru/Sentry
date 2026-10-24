<script lang="ts">
  interface Props {
    busy?: boolean;
    disabled?: boolean;
    onclick: () => void;
  }
  let { busy = false, disabled = false, onclick }: Props = $props();
</script>

<button
  type="button"
  class="check-btn"
  class:busy
  disabled={busy || disabled}
  onclick={onclick}
  aria-busy={busy}
>
  {#if busy}
    <span class="spinner" aria-hidden="true"></span>
    <span>Checking...</span>
  {:else}
    <span>Check message</span>
  {/if}
</button>

<style>
  .check-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    width: 100%;
    min-height: var(--touch-min);
    padding: 16px 24px;
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-inverse);
    background-color: var(--text-primary);
    border: 3px solid var(--text-primary);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.15s ease, opacity 0.15s ease;
  }

  .check-btn:hover:not(:disabled) {
    background-color: #292524;
  }

  .check-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 3px solid rgba(255, 255, 255, 0.3);
    border-top-color: #ffffff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
