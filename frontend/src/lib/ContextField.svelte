<script lang="ts">
  let { value = $bindable('') } = $props<{ value: string }>();
</script>

<div class="field-group">
  <label for="context-input" class="label">What happened</label>
  <textarea
    id="context-input"
    class="context-area"
    bind:value={value}
    rows="3"
    spellcheck="true"
  ></textarea>
</div>

<style>
  .field-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .label {
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-primary);
  }

  .context-area {
    width: 100%;
    min-height: 100px;
    padding: 16px;
    font-size: var(--font-base);
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    resize: vertical;
    line-height: 1.4;
  }

  .context-area:focus {
    outline: 3px solid var(--text-primary);
    outline-offset: 2px;
  }
</style>
