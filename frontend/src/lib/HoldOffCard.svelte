<script lang="ts">
  interface Props {
    advice: string;
  }
  let { advice }: Props = $props();
</script>

<div class="result-card hold-off" role="region" aria-label="What to hold off on">
  <div class="card-header">
    <svg class="header-icon" viewBox="0 0 24 24" width="32" height="32" aria-hidden="true">
      <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
    </svg>
    <h2 class="title">Hold off on</h2>
  </div>
  <p class="content">{advice}</p>
</div>

<style>
  .result-card {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 24px;
    background-color: var(--rose-bg);
    border: 3px solid var(--rose-border);
    border-radius: var(--radius-md);
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .header-icon {
    color: var(--rose-text);
    flex-shrink: 0;
  }

  .title {
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--rose-text);
    letter-spacing: -0.01em;
  }

  .content {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--rose-text);
    line-height: 1.45;
  }
</style>
