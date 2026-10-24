<script lang="ts">
  import CaregiverSignIn from '$lib/CaregiverSignIn.svelte';

  interface Incident {
    id: string;
    score: number;
    risk_level: string;
    category: string;
    status: string;
    user_context?: string;
    extracted_text?: string;
    created_at: string;
  }

  let signedIn = $state(false);
  let incidents = $state<Incident[]>([]);
  let loading = $state(false);

  function handleSignIn(email: string, pass: string) {
    signedIn = true;
    loadIncidents();
  }

  async function loadIncidents() {
    loading = true;
    try {
      const res = await fetch('/api/caregiver/incidents');
      if (res.ok) {
        incidents = await res.json();
      }
    } catch {
      // ignore
    } finally {
      loading = false;
    }
  }
</script>

<main class="caregiver-page">
  <header class="caregiver-header">
    <div class="brand-group">
      <a href="/" class="back-link">&larr; Sentry</a>
      <h1 class="portal-title">Caregiver Portal</h1>
    </div>
    {#if signedIn}
      <button type="button" class="signout-btn" onclick={() => { signedIn = false; }}>
        Sign out
      </button>
    {/if}
  </header>

  {#if !signedIn}
    <CaregiverSignIn onsignin={handleSignIn} />
  {:else}
    <section class="incident-feed">
      {#if loading && incidents.length === 0}
        <p class="status-msg">Loading incidents...</p>
      {:else if incidents.length === 0}
        <p class="status-msg">No incidents reported yet.</p>
      {:else}
        {#each incidents as inc (inc.id)}
          <article class="incident-card">
            <div class="card-meta">
              <span class="incident-id">{inc.id}</span>
              <span class="incident-time">{inc.created_at}</span>
            </div>
            <p class="incident-text">
              {inc.user_context || inc.extracted_text || inc.category}
            </p>
          </article>
        {/each}
      {/if}
    </section>
  {/if}
</main>

<style>
  .caregiver-page {
    display: flex;
    flex-direction: column;
    gap: 28px;
  }

  .caregiver-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 16px;
    border-bottom: 2px solid var(--border-subtle);
  }

  .brand-group {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .back-link {
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-secondary);
  }

  .portal-title {
    font-size: var(--font-xl);
    font-weight: 800;
  }

  .signout-btn {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-muted);
    padding: 6px 12px;
  }

  .incident-feed {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .status-msg {
    font-size: var(--font-lg);
    color: var(--text-muted);
    text-align: center;
    padding: 40px 0;
  }

  .incident-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .card-meta {
    display: flex;
    justify-content: space-between;
    font-size: 0.875rem;
    color: var(--text-muted);
  }

  .incident-id {
    font-weight: 700;
  }

  .incident-text {
    font-size: var(--font-base);
    color: var(--text-primary);
  }
</style>
