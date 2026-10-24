<script lang="ts">
  interface Props {
    level: string;
    score?: number;
  }
  let { level, score }: Props = $props();

  const normalized = $derived((level || 'low').toLowerCase());
</script>

<span
  class="chip {normalized}"
  role="status"
  aria-label={`Risk level ${normalized}`}
>
  {#if normalized === 'high'}
    High Risk {score !== undefined ? `(${score})` : ''}
  {:else if normalized === 'medium'}
    Caution {score !== undefined ? `(${score})` : ''}
  {:else}
    Low Risk {score !== undefined ? `(${score})` : ''}
  {/if}
</span>

<style>
  .chip {
    display: inline-flex;
    align-items: center;
    padding: 4px 12px;
    font-size: 0.875rem;
    font-weight: 800;
    border-radius: 999px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    border: 2px solid currentColor;
  }

  .chip.high {
    background-color: var(--rose-bg);
    border-color: var(--rose-border);
    color: var(--rose-text);
  }

  .chip.medium {
    background-color: var(--amber-bg);
    border-color: var(--amber-border);
    color: var(--amber-text);
  }

  .chip.low {
    background-color: var(--stone-bg);
    border-color: var(--stone-border);
    color: var(--stone-text);
  }
</style>
