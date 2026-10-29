<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import CaregiverSignIn from '$lib/CaregiverSignIn.svelte';
  import SeverityChip from '$lib/SeverityChip.svelte';
  import { signIn, signOut, isAuthenticated } from '$lib/auth';
  import {
    listCaregiverIncidents,
    updateCaregiverStatus,
    createCaregiverInvite,
    fetchCaregiverLinks,
    removeCaregiverLink,
    type Incident,
    type CaregiverLinkInfo
  } from '$lib/api';

  let signedIn = $state(false);
  let signInBusy = $state(false);
  let signInError = $state('');
  let incidents = $state<Incident[]>([]);
  let loading = $state(false);
  let feedError = $state('');
  let actionBusy = $state<Record<string, boolean>>({});
  let pollInterval: ReturnType<typeof setInterval> | null = null;
  let links = $state<CaregiverLinkInfo[]>([]);
  let inviteCode = $state('');
  let inviteSeniorName = $state('');
  let inviteBusy = $state(false);
  let showInviteForm = $state(false);
  let removeBusy = $state<Record<string, boolean>>({});
  let activeFilter = $state<'open' | 'acknowledged' | 'resolved'>('open');

  const seniorMap = $derived.by(() => {
    const map: Record<string, string> = {};
    for (const l of links) {
      if (l.senior_id) {
        map[l.senior_id] = l.senior_name || 'Senior';
      }
    }
    return map;
  });

  const filteredIncidents = $derived.by(() => {
    return incidents.filter((inc) => {
      if (activeFilter === 'open') {
        return inc.status !== 'acknowledged' && inc.status !== 'resolved';
      }
      return inc.status === activeFilter;
    });
  });

  const groupedIncidents = $derived.by(() => {
    const groups: { seniorId: string; seniorName: string; items: Incident[] }[] = [];
    const groupMap = new Map<string, Incident[]>();

    for (const inc of filteredIncidents) {
      const sid = inc.senior_id || 'unknown';
      if (!groupMap.has(sid)) {
        groupMap.set(sid, []);
      }
      groupMap.get(sid)!.push(inc);
    }

    for (const [sid, items] of groupMap.entries()) {
      groups.push({
        seniorId: sid,
        seniorName: seniorMap[sid] || (sid === 'unknown' ? 'Senior' : 'Senior'),
        items
      });
    }

    return groups;
  });

  onMount(() => {
    signedIn = isAuthenticated();
    if (signedIn) {
      loadIncidents();
      loadLinks();
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
      await Promise.all([loadIncidents(), loadLinks()]);
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
    links = [];
    inviteCode = '';
  }

  async function loadLinks() {
    try {
      links = await fetchCaregiverLinks();
    } catch {
      // quiet fallback
    }
  }

  async function handleCreateInvite() {
    inviteBusy = true;
    feedError = '';
    try {
      const res = await createCaregiverInvite(inviteSeniorName.trim());
      inviteCode = res.code;
    } catch (err: unknown) {
      feedError = err instanceof Error ? err.message : 'Could not create invite';
    } finally {
      inviteBusy = false;
    }
  }

  async function handleRemoveLink(seniorId: string) {
    removeBusy[seniorId] = true;
    feedError = '';
    try {
      await removeCaregiverLink(seniorId);
      links = links.filter((l) => l.senior_id !== seniorId);
    } catch (err: unknown) {
      feedError = err instanceof Error ? err.message : 'Could not remove link';
    } finally {
      removeBusy[seniorId] = false;
    }
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

    <section class="links-panel">
      <div class="links-header">
        <h2 class="section-title">Linked Seniors</h2>
        <button
          type="button"
          class="create-invite-btn"
          onclick={() => { showInviteForm = !showInviteForm; inviteCode = ''; }}
        >
          {showInviteForm ? 'Close' : 'Invite Senior'}
        </button>
      </div>

      {#if showInviteForm}
        <div class="invite-box">
          <div class="invite-inputs">
            <input
              type="text"
              class="senior-name-input"
              bind:value={inviteSeniorName}
              placeholder="Senior name (e.g. Mom)"
              maxlength="40"
              aria-label="Senior name"
            />
            <button
              type="button"
              class="generate-btn"
              disabled={inviteBusy}
              onclick={handleCreateInvite}
            >
              {inviteBusy ? 'Generating...' : 'Create Code'}
            </button>
          </div>

          {#if inviteCode}
            <div class="invite-result">
              <span class="invite-label">Code (valid 24h):</span>
              <strong class="invite-code">{inviteCode}</strong>
            </div>
          {/if}
        </div>
      {/if}

      {#if links.length > 0}
        <ul class="links-list">
          {#each links as link (link.senior_id)}
            <li class="link-item">
              <span class="link-name">{link.senior_name || 'Senior'}</span>
              <button
                type="button"
                class="remove-link-btn"
                disabled={removeBusy[link.senior_id]}
                onclick={() => handleRemoveLink(link.senior_id)}
              >
                {removeBusy[link.senior_id] ? 'Removing...' : 'Remove'}
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section class="incident-feed">
      <div class="filters-bar" role="tablist" aria-label="Incident status filters">
        <button
          type="button"
          class="filter-tab"
          class:active={activeFilter === 'open'}
          role="tab"
          aria-selected={activeFilter === 'open'}
          onclick={() => activeFilter = 'open'}
        >
          Open
        </button>
        <button
          type="button"
          class="filter-tab"
          class:active={activeFilter === 'acknowledged'}
          role="tab"
          aria-selected={activeFilter === 'acknowledged'}
          onclick={() => activeFilter = 'acknowledged'}
        >
          Acknowledged
        </button>
        <button
          type="button"
          class="filter-tab"
          class:active={activeFilter === 'resolved'}
          role="tab"
          aria-selected={activeFilter === 'resolved'}
          onclick={() => activeFilter = 'resolved'}
        >
          Resolved
        </button>
      </div>

      {#if loading && incidents.length === 0}
        <p class="status-msg">Loading incidents...</p>
      {:else if groupedIncidents.length === 0}
        <div class="empty-state">
          <p class="status-msg">No {activeFilter} incidents</p>
        </div>
      {:else}
        {#each groupedIncidents as group (group.seniorId)}
          <div class="senior-group">
            <h3 class="senior-group-header">{group.seniorName}</h3>
            {#each group.items as inc (inc.id)}
              <article class="incident-card" class:high-border={inc.risk_level === 'high'}>
                <div class="card-top">
                  <div class="card-badges">
                    <SeverityChip level={inc.risk_level} score={inc.score} />
                    <span class="category-chip">{inc.category}</span>
                  </div>
                  <time class="incident-time" datetime={inc.created_at}>
                    {formatTime(inc.created_at)}
                  </time>
                </div>

                {#if inc.signals && inc.signals.length > 0}
                  <div class="signals-wrap">
                    {#each inc.signals as sig}
                      <span class="signal-chip">{sig.name}</span>
                    {/each}
                  </div>
                {/if}

                {#if inc.thumbnail_url}
                  <div class="thumbnail-container">
                    <img
                      src={inc.thumbnail_url}
                      alt="Screenshot thumbnail"
                      class="screenshot-thumbnail"
                      loading="lazy"
                    />
                  </div>
                {/if}

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
          </div>
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

  .links-panel {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 20px;
    background-color: var(--bg-surface);
    border: 2px solid var(--border-subtle);
    border-radius: var(--radius-md);
  }

  .links-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--text-primary);
    margin: 0;
  }

  .create-invite-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 16px;
    font-size: var(--font-sm);
    font-weight: 700;
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .create-invite-btn:hover,
  .create-invite-btn:focus-visible {
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .invite-box {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background-color: var(--bg-canvas);
    border: 2px dashed var(--border-subtle);
    border-radius: var(--radius-sm);
  }

  .invite-inputs {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }

  .senior-name-input {
    flex: 1;
    min-width: 200px;
    min-height: var(--touch-min);
    padding: 0 14px;
    font-size: var(--font-base);
    color: var(--text-primary);
    background-color: var(--bg-surface);
    border: 2px solid var(--border-strong);
    border-radius: var(--radius-sm);
  }

  .senior-name-input:focus-visible {
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .generate-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 18px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--bg-surface);
    background-color: var(--text-primary);
    border: 2px solid var(--text-primary);
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .generate-btn:hover:not(:disabled),
  .generate-btn:focus-visible:not(:disabled) {
    opacity: 0.9;
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .generate-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .invite-result {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background-color: var(--emerald-50, #ecfdf5);
    border: 2px solid var(--emerald-600, #059669);
    border-radius: var(--radius-sm);
  }

  .invite-label {
    font-size: var(--font-sm);
    font-weight: 700;
    color: var(--emerald-800, #065f46);
  }

  .invite-code {
    font-size: var(--font-xl);
    font-weight: 800;
    letter-spacing: 0.15em;
    color: var(--emerald-900, #064e3b);
  }

  .links-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .link-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 14px;
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
  }

  .link-name {
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-primary);
  }

  .remove-link-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 12px;
    font-size: var(--font-sm);
    font-weight: 700;
    color: var(--rose-800);
    background-color: var(--rose-50);
    border: 1px solid var(--rose-600);
    border-radius: var(--radius-sm);
    cursor: pointer;
  }

  .remove-link-btn:hover:not(:disabled),
  .remove-link-btn:focus-visible:not(:disabled) {
    background-color: var(--rose-100, #ffe4e6);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .remove-link-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .incident-feed {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .filters-bar {
    display: flex;
    gap: 8px;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--border-subtle);
  }

  .filter-tab {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: var(--touch-min);
    padding: 0 18px;
    font-size: var(--font-base);
    font-weight: 700;
    color: var(--text-secondary);
    background: transparent;
    border: 2px solid transparent;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 150ms ease;
  }

  .filter-tab:hover,
  .filter-tab:focus-visible {
    color: var(--text-primary);
    background-color: var(--bg-hover);
    outline: 3px solid var(--ring-focus);
    outline-offset: 2px;
  }

  .filter-tab.active {
    color: var(--bg-surface);
    background-color: var(--text-primary);
    border-color: var(--text-primary);
  }

  .senior-group {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .senior-group-header {
    font-size: var(--font-lg);
    font-weight: 800;
    color: var(--text-primary);
    margin: 8px 0 0;
    padding-bottom: 4px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .empty-state {
    display: flex;
    justify-content: center;
    padding: 32px 0;
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
    flex-wrap: wrap;
    gap: 8px;
  }

  .card-badges {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .category-chip {
    display: inline-flex;
    align-items: center;
    padding: 4px 10px;
    font-size: var(--font-xs);
    font-weight: 700;
    text-transform: capitalize;
    color: var(--text-primary);
    background-color: var(--bg-hover);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
  }

  .signals-wrap {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .signal-chip {
    display: inline-flex;
    align-items: center;
    padding: 3px 8px;
    font-size: var(--font-xs);
    font-weight: 600;
    color: var(--text-secondary);
    background-color: var(--bg-canvas);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
  }

  .thumbnail-container {
    max-width: 240px;
    max-height: 160px;
    overflow: hidden;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-subtle);
  }

  .screenshot-thumbnail {
    width: 100%;
    height: auto;
    max-height: 160px;
    object-fit: cover;
    display: block;
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
