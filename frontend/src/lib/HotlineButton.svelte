<script lang="ts">
  interface Props {
    hotline: {
      country: string;
      name: string;
      number: string;
      tel_uri: string;
      url?: string;
    };
  }
  let { hotline }: Props = $props();
</script>

<a
  href={hotline.tel_uri}
  class="hotline-btn"
  aria-label={`Dial ${hotline.name} directly on ${hotline.number}`}
>
  <svg class="phone-icon" viewBox="0 0 24 24" width="28" height="28" aria-hidden="true">
    <path
      fill="currentColor"
      d="M6.62 10.79a15.05 15.05 0 006.59 6.59l2.2-2.2a1 1 0 011.02-.24c1.12.37 2.33.57 3.57.57a1 1 0 011 1V20a1 1 0 01-1 1A17 17 0 013 4a1 1 0 011-1h3.5a1 1 0 011 1c0 1.25.2 2.45.57 3.57a1 1 0 01-.25 1.02l-2.2 2.2z"
    />
  </svg>
  <span class="btn-text">
    Call {hotline.name} ({hotline.number})
  </span>
</a>

<style>
  .hotline-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 14px;
    width: 100%;
    min-height: var(--touch-min);
    padding: 18px 24px;
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--text-inverse);
    background-color: var(--text-primary);
    border: 3px solid var(--text-primary);
    border-radius: var(--radius-sm);
    text-align: center;
    cursor: pointer;
    transition: background-color 0.15s ease;
  }

  .hotline-btn:hover {
    background-color: #292524;
  }

  .phone-icon {
    flex-shrink: 0;
  }

  .btn-text {
    line-height: 1.3;
  }
</style>
