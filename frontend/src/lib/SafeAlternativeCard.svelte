<script lang="ts">
  let { action } = $props<{ action: string }>();
</script>

<div class="result-card safe-alternative" role="region" aria-label="Safe alternative action">
  <div class="card-header">
    <svg class="header-icon" viewBox="0 0 24 24" width="28" height="28" aria-hidden="true">
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
    gap: 12px;
    padding: 24px;
    background-color: var(--amber-bg);
    border: 3px solid var(--amber-border);
    border-radius: var(--radius-md);
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
  }

  .content {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--amber-text);
    line-height: 1.4;
  }
</style>
