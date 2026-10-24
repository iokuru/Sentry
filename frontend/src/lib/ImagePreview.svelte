<script lang="ts">
  interface Props {
    src: string;
    onremove: () => void;
  }
  let { src, onremove }: Props = $props();
</script>

<div class="preview-card">
  <img {src} alt="Uploaded message preview" class="preview-img" />
  <button type="button" class="remove-btn" onclick={onremove} aria-label="Remove uploaded image">
    Remove image
  </button>
</div>

<style>
  .preview-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .preview-img {
    max-height: 280px;
    width: 100%;
    object-fit: contain;
    border-radius: var(--radius-sm);
    background-color: var(--stone-bg);
  }

  .remove-btn {
    min-height: var(--touch-min);
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--rose-text);
    background-color: var(--rose-bg);
    border: 2px solid var(--rose-border);
    border-radius: var(--radius-sm);
    padding: 12px 20px;
    cursor: pointer;
    text-align: center;
  }

  .remove-btn:hover {
    background-color: #fee2e2;
  }
</style>
