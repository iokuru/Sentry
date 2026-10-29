<script lang="ts">
  interface Props {
    onfile: (file: File) => void;
    ontext?: (text: string) => void;
    hasFile?: boolean;
  }
  let { onfile, ontext, hasFile = false }: Props = $props();

  let isDragging = $state(false);
  let fileInput: HTMLInputElement;
  let cameraInput: HTMLInputElement;

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
    if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
      onfile(e.dataTransfer.files[0]);
    } else if (e.dataTransfer) {
      const text = e.dataTransfer.getData('text/plain');
      if (text && ontext) {
        ontext(text);
      }
    }
  }

  function handleFileChange(e: Event) {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      onfile(target.files[0]);
    }
  }

  function triggerSelect() {
    fileInput?.click();
  }

  function triggerCamera() {
    cameraInput?.click();
  }
</script>

<div
  class="dropzone"
  class:dragging={isDragging}
  class:has-file={hasFile}
  ondragover={(e) => { e.preventDefault(); isDragging = true; }}
  ondragleave={() => { isDragging = false; }}
  ondrop={handleDrop}
  onclick={triggerSelect}
  onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); triggerSelect(); } }}
  role="button"
  tabindex="0"
  aria-label="Upload screenshot or photo of suspicious message"
>
  <input
    bind:this={fileInput}
    type="file"
    accept="image/png,image/jpeg,image/webp,image/heic,image/heif,.heic,.heif"
    class="file-input"
    onchange={handleFileChange}
  />
  <input
    bind:this={cameraInput}
    type="file"
    accept="image/*"
    capture="environment"
    class="file-input"
    onchange={handleFileChange}
  />

  <svg class="icon" viewBox="0 0 24 24" width="48" height="48" aria-hidden="true">
    <path
      fill="currentColor"
      d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"
    />
  </svg>

  <span class="label">
    {hasFile ? 'Replace image' : 'Drop message or screenshot here'}
  </span>

  <div class="drop-actions">
    <button
      type="button"
      class="drop-action-btn"
      onclick={(e) => { e.stopPropagation(); triggerSelect(); }}
    >
      Upload image
    </button>
    <button
      type="button"
      class="drop-action-btn"
      onclick={(e) => { e.stopPropagation(); triggerCamera(); }}
    >
      Take photo
    </button>
  </div>
</div>

<style>
  .dropzone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    min-height: 180px;
    padding: 32px 24px;
    background-color: var(--bg-surface);
    border: 3px dashed var(--border-strong);
    border-radius: var(--radius-md);
    cursor: pointer;
    text-align: center;
    transition: background-color 0.15s ease, border-color 0.15s ease;
  }

  .dropzone:hover,
  .dropzone.dragging {
    background-color: var(--stone-bg);
    border-color: var(--text-primary);
  }

  .dropzone.has-file {
    min-height: 80px;
    padding: 16px 20px;
    border-style: solid;
  }

  .file-input {
    display: none;
  }

  .icon {
    color: var(--text-primary);
  }

  .label {
    font-size: var(--font-lg);
    font-weight: 700;
    color: var(--text-primary);
    line-height: 1.3;
  }

  .drop-actions {
    display: flex;
    gap: 12px;
    margin-top: 8px;
    flex-wrap: wrap;
    justify-content: center;
  }

  .drop-action-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 16px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .drop-action-btn:hover,
  .drop-action-btn:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }
</style>
