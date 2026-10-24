<script lang="ts">
  import { onDestroy } from 'svelte';
  import CaregiverSignIn from '$lib/CaregiverSignIn.svelte';
  import SeverityChip from '$lib/SeverityChip.svelte';

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
  let actionBusy = $state<Record<string, boolean>>({});
  let pollInterval: ReturnType<typeof setInterval> | null = null;

  function handleSignIn(email: string, pass: string) {
    signedIn = true;
    loadIncidents();
    startPolling();
  }

  function handleSignOut() {
    stopPolling();
    signedIn = false;
  }

  function startPolling() {
    stopPolling();
    // poll every five seconds for fresh reports
    pollInterval = setInterval(loadIncidents, 5000);
  }

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
  }

  onDestroy(() => {
    stopPolling();
  });

  async function loadIncidents() {
    if (incidents.length === 0) {
      loading = true;
    }
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

  async function updateStatus(id: string, status: 'acknowledged' | 'resolved') {
    actionBusy[id] = true;
    try {
      const res = await fetch('/api/caregiver/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ incident_id: id, status })
      });

      if (res.ok) {
        const idx = incidents.findIndex((item) => item.id === id);
        if (idx !== -1) {
          incidents[idx].status = status;
        }
      }
    } catch {
      // ignore
    } finally {
      actionBusy[id] = false;
    }
  }

  function formatTime(iso: string): string {
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

<main class="caregiver-page">
  <header class="caregiver-header">
    <div class="brand-group">
      <a href="/" class="back-link">&larr; Sentry</a>
      <h1 class="portal-title">Caregiver Portal</h1>
    </div>
    {#if signedIn}
      <button type="button" class="signout-btn" onclick={handleSignOut}>
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
          <article class="incident-card" class:high-border={inc.risk_level === 'high'}>
            <div class="card-top">
              <SeverityChip level={inc.risk_level} score={inc.score} />
              <time class="incident-time" datetime={inc.created_at}>
                {formatTime(inc.created_at)}
              </time>
            </div>

            <p class="incident-text">
              {inc.user_context || inc.extracted_text || `Suspicious ${inc.category} message`}
            </p>

            <div class="card-actions">
              {#if inc.status !== 'resolved'}
                {#if inc.status !== 'acknowledged'}
                  <button
                    type="button"
                    class="action-btn ack-btn"
                    disabled={actionBusy[inc.id]}
                    onclick={() => updateStatus(inc.id, 'acknowledged')}
                  >
                    Acknowledge
                  </button>
                {/if}
                <button
                  type="button"
                  class="action-btn resolve-btn"
                  disabled={actionBusy[inc.id]}
                  onclick={() => updateStatus(inc.id, 'resolved')}
                >
                  Resolve
                </button>
              {:else}
                <span class="resolved-label">Resolved</span>
              {/if}
            </div>
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
    gap: 16px;
    padding: 24px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .incident-card.high-border {
    border-color: var(--rose-border);
  }

  .card-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .incident-time {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text-muted);
  }

  .incident-text {
    font-size: var(--font-base);
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.4;
  }

  .card-actions {
    display: flex;
    gap: 12px;
    margin-top: 4px;
  }

  .action-btn {
    flex: 1;
    min-height: var(--touch-min);
    font-size: var(--font-base);
    font-weight: 700;
    border-radius: var(--radius-sm);
    padding: 12px 20px;
    cursor: pointer;
    text-align: center;
    transition: background-color 0.15s ease;
  }

  .ack-btn {
    color: var(--amber-text);
    background-color: var(--amber-bg);
    border: 2px solid var(--amber-border);
  }

  .ack-btn:hover:not(:disabled) {
    background-color: #fef3c7;
  }

  .resolve-btn {
    color: var(--text-inverse);
    background-color: var(--text-primary);
    border: 2px solid var(--text-primary);
  }

  .resolve-btn:hover:not(:disabled) {
    background-color: #292524;
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .resolved-label {
    font-size: var(--font-base);
    font-weight: 700;
    color: #065f46;
    background-color: #ecfdf5;
    padding: 8px 16px;
    border-radius: var(--radius-sm);
  }
</style>
