<script lang="ts">
  import type { Incident } from '$lib/api';

  interface Props {
    items: Incident[];
    onselect: (item: Incident) => void;
    onclear: () => void;
  }

  let { items, onselect, onclear }: Props = $props();
  let expanded = $state(false);

  function formatDate(iso: string): string {
    if (!iso) return '';
    try {
      const d = new Date(iso);
      return d.toLocaleDateString(undefined, {
        month: 'short',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit'
      });
    } catch {
      return iso;
    }
  }
</script>

{#if items.length > 0}
  <div class="history-container">
    <div class="history-header">
      <button
        type="button"
        class="toggle-btn"
        onclick={() => expanded = !expanded}
        aria-expanded={expanded}
      >
        <span>Recent checks ({items.length})</span>
        <span class="chevron" aria-hidden="true">{expanded ? '▲' : '▼'}</span>
      </button>

      {#if expanded}
        <button type="button" class="clear-btn" onclick={onclear}>
          Clear history
        </button>
      {/if}
    </div>

    {#if expanded}
      <ul class="history-list">
        {#each items as item (item.id)}
          <li>
            <button
              type="button"
              class="history-item-btn"
              onclick={() => onselect(item)}
            >
              <div class="item-meta">
                <span class="category-name">{item.category}</span>
                <span class="history-time">{formatDate(item.created_at)}</span>
              </div>
              <span class="risk-badge" class:high={item.risk_level === 'high'}>
                {item.risk_level} risk
              </span>
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
{/if}

<style>
  .history-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 2px solid var(--border-subtle);
  }

  .history-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .toggle-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    min-height: var(--touch-min);
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 8px 0;
  }

  .toggle-btn:hover,
  .toggle-btn:focus-visible {
    color: var(--text-secondary);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .chevron {
    font-size: var(--font-sm);
  }

  .clear-btn {
    min-height: var(--touch-min);
    padding: 0 12px;
    font-size: var(--font-sm);
    font-weight: 600;
    color: var(--text-muted);
    background: transparent;
    border: none;
    cursor: pointer;
    text-decoration: underline;
  }

  .clear-btn:hover,
  .clear-btn:focus-visible {
    color: var(--rose-700);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .history-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .history-item-btn {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    min-height: var(--touch-min);
    padding: 12px 16px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    cursor: pointer;
    text-align: left;
    transition: background-color 150ms ease;
  }

  .history-item-btn:hover,
  .history-item-btn:focus-visible {
    background-color: var(--bg-hover);
    border-color: var(--border-strong);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .item-meta {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .category-name {
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    text-transform: capitalize;
  }

  .history-time {
    font-size: var(--font-sm);
    color: var(--text-muted);
  }

  .risk-badge {
    padding: 4px 10px;
    font-size: var(--font-xs);
    font-weight: 700;
    text-transform: uppercase;
    border-radius: var(--radius-sm);
    background-color: var(--stone-100, #f5f5f4);
    color: var(--stone-800, #292524);
    border: 1px solid var(--stone-300, #d6d3d1);
  }

  .risk-badge.high {
    background-color: var(--rose-50);
    color: var(--rose-800);
    border-color: var(--rose-600);
  }
</style>
