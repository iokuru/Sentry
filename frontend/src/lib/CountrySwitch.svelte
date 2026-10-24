<script lang="ts">
  interface Props {
    country?: string;
  }
  let { country = $bindable('US') }: Props = $props();

  const countries = [
    { code: 'US', label: 'US' },
    { code: 'GB', label: 'UK' },
    { code: 'IN', label: 'IN' },
    { code: 'AU', label: 'AU' }
  ];

  function cycleCountry() {
    const idx = countries.findIndex((c) => c.code === country);
    const next = (idx + 1) % countries.length;
    country = countries[next].code;
  }
</script>

<button
  type="button"
  class="country-btn"
  onclick={cycleCountry}
  title="Switch country hotline"
  aria-label={`Current country is ${country}. Click to switch.`}
>
  <span class="flag-code">{country === 'GB' ? 'UK' : country}</span>
</button>

<style>
  .country-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 6px 10px;
    font-size: 0.875rem;
    font-weight: 700;
    color: var(--text-secondary);
    background-color: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    cursor: pointer;
    min-height: 36px;
    min-width: 44px;
  }

  .country-btn:hover {
    color: var(--text-primary);
    border-color: var(--border-strong);
  }

  .flag-code {
    letter-spacing: 0.05em;
  }
</style>
