<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import CaregiverSignIn from '$lib/CaregiverSignIn.svelte';
  import SeverityChip from '$lib/SeverityChip.svelte';
  import { signIn, signOut, isAuthenticated } from '$lib/auth';
  import { listCaregiverIncidents, updateCaregiverStatus, type Incident } from '$lib/api';

  let signedIn = $state(false);
  let signInBusy = $state(false);
  let signInError = $state('');
  let incidents = $state<Incident[]>([]);
  let loading = $state(false);
  let feedError = $state('');
  let actionBusy = $state<Record<string, boolean>>({});
  let pollInterval: ReturnType<typeof setInterval> | null = null;

  onMount(() => {
    signedIn = isAuthenticated();
    if (signedIn) {
      loadIncidents();
      startPolling();
    }
    if (typeof document !== 'undefined') {
      document.addEventListener('visibilitychange', handleVisibilityChange);
    }
  });

  onDestroy(() => {
    stopPolling();
    if (typeof document !== 'undefined') {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    }
  });

  function handleVisibilityChange() {
    if (document.visibilityState === 'hidden') {
      stopPolling();
    } else if (document.visibilityState === 'visible' && signedIn) {
      loadIncidents();
      startPolling();
    }
  }

  async function handleSignIn(email: string, pass: string) {
    signInBusy = true;
    signInError = '';
    try {
      await signIn(email, pass);
      signedIn = true;
      await loadIncidents();
      startPolling();
    } catch (err: unknown) {
      signInError = err instanceof Error ? err.message : 'Sign in failed';
    } finally {
      signInBusy = false;
    }
  }

  function handleSignOut() {
    stopPolling();
    signOut();
    signedIn = false;
    incidents = [];
  }

  function startPolling() {
    stopPolling();
    // poll every fifteen seconds for fresh caregiver reports
    pollInterval = setInterval(loadIncidents, 15000);
  }

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
  }

  async function loadIncidents() {
    if (incidents.length === 0) {
      loading = true;
    }
    try {
      incidents = await listCaregiverIncidents();
      feedError = '';
    } catch (err: unknown) {
      if (err instanceof Error && err.message === 'Unauthorized') {
        handleSignOut();
      } else {
        feedError = 'Failed to refresh incidents';
      }
    } finally {
      loading = false;
    }
  }

  async function updateStatus(id: string, newStatus: 'acknowledged' | 'resolved') {
    const idx = incidents.findIndex((item) => item.id === id);
    if (idx === -1) return;

    const previousStatus = incidents[idx].status;
    // optimistic update
    incidents[idx].status = newStatus;
    actionBusy[id] = true;
    feedError = '';

    try {
      await updateCaregiverStatus(id, newStatus);
    } catch {
      // rollback status on failure
      incidents[idx].status = previousStatus;
      feedError = 'Could not update status, please retry';
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
    <CaregiverSignIn
      busy={signInBusy}
      errorMessage={signInError}
      onsignin={handleSignIn}
    />
  {:else}
    {#if feedError}
      <div class="feed-error" role="alert">
        {feedError}
      </div>
    {/if}

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
    max-width: 720px;
    margin: 0 auto;
    padding: 20px 24px 60px;
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
    text-decoration: none;
    min-height: var(--touch-min);
    display: inline-flex;
    align-items: center;
  }

  .back-link:hover,
  .back-link:focus-visible {
    color: var(--text-primary);
  }

  .portal-title {
    font-size: var(--font-xl);
    font-weight: 800;
    color: var(--text-primary);
  }

  .signout-btn {
    min-height: var(--touch-min);
    padding: 0 16px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-secondary);
    background: transparent;
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .signout-btn:hover,
  .signout-btn:focus-visible {
    color: var(--text-primary);
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .feed-error {
    padding: 12px 16px;
    background-color: var(--rose-50);
    border: 2px solid var(--rose-600);
    color: var(--rose-800);
    font-size: var(--font-base);
    font-weight: 700;
    border-radius: var(--radius-sm);
  }

  .incident-feed {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .status-msg {
    font-size: var(--font-lg);
    color: var(--text-secondary);
    padding: 32px 0;
    text-align: center;
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

  .high-border {
    border-left: 6px solid var(--rose-600);
  }

  .card-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .incident-time {
    font-size: var(--font-sm);
    color: var(--text-muted);
    font-weight: 600;
  }

  .incident-text {
    font-size: var(--font-base);
    color: var(--text-primary);
    line-height: 1.5;
    word-break: break-word;
  }

  .card-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    padding-top: 8px;
  }

  .action-btn {
    min-height: var(--touch-min);
    min-width: 140px;
    padding: 12px 20px;
    font-size: var(--font-base);
    font-weight: 700;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .ack-btn {
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    color: var(--text-primary);
  }

  .ack-btn:hover:not(:disabled),
  .ack-btn:focus-visible:not(:disabled) {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .resolve-btn {
    background-color: var(--stone-900);
    border: 2px solid var(--stone-900);
    color: var(--text-inverse);
  }

  .resolve-btn:hover:not(:disabled),
  .resolve-btn:focus-visible:not(:disabled) {
    background-color: var(--stone-700);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .resolved-label {
    min-height: var(--touch-min);
    display: inline-flex;
    align-items: center;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-muted);
  }
</style>
