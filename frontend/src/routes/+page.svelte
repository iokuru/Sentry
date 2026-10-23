<script lang="ts">
  let country = $state('US');
</script>

<main class="page">
  <header class="top-bar">
    <span class="brand">Sentry</span>
    <nav class="nav">
      <a href="/caregiver" class="portal-link">Caregiver</a>
    </nav>
  </header>

  <section class="intake-section">
    <!-- intake layout -->
  </section>
</main>

<style>
  .page {
    display: flex;
    flex-direction: column;
    gap: 32px;
  }

  .top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 16px;
    border-bottom: 2px solid var(--border-subtle);
  }

  .brand {
    font-size: var(--font-xl);
    font-weight: 700;
    letter-spacing: -0.02em;
    color: var(--text-primary);
  }

  .portal-link {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-muted);
    padding: 8px 12px;
    border-radius: var(--radius-sm);
  }

  .portal-link:hover {
    color: var(--text-primary);
  }

  .intake-section {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }
</style>
