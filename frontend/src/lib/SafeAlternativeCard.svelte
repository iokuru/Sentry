<script lang="ts">
  let { action } = $props<{ action: string }>();
</script>

<div class="result-card safe-alternative" role="region" aria-label="Safe alternative action" tabindex="0">
  <div class="card-header">
    <svg class="header-icon" viewBox="0 0 24 24" width="32" height="32" aria-hidden="true">
      <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
    </svg>
    <h2 class="title">Safe next step</h2>
  </div>
  <p class="content">{action}</p>
</div>

<style>
  .result-card {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 24px;
    background-color: var(--amber-bg);
    border: 3px solid var(--amber-border);
    border-radius: var(--radius-md);
  }

  .result-card:focus-visible {
    outline: 4px solid var(--amber-text);
    outline-offset: 3px;
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .header-icon {
    color: var(--amber-text);
    flex-shrink: 0;
  }

  .title {
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--amber-text);
    letter-spacing: -0.01em;
  }

  .content {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--amber-text);
    line-height: 1.45;
  }
</style>
