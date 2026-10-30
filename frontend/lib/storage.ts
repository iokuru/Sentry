import type { IncidentResult } from "./api";

export interface HistoryEntry {
  id: string;
  timestamp: string;
  severity: "HIGH" | "MEDIUM" | "LOW";
  summary: string;
}

const LAST_RESULT_KEY = "sentry_last_result";
const HISTORY_KEY = "sentry_check_history";

export function saveLastResult(result: IncidentResult): void {
  if (typeof window === "undefined") return;

  try {
    localStorage.setItem(LAST_RESULT_KEY, JSON.stringify(result));

    const history = getCheckHistory();
    const newEntry: HistoryEntry = {
      id: result.incident_id,
      timestamp: new Date().toISOString(),
      severity: result.severity || "MEDIUM",
      summary: result.hold_off || "Check completed",
    };

    const updated = [newEntry, ...history.filter((h) => h.id !== result.incident_id)].slice(0, 10);
    localStorage.setItem(HISTORY_KEY, JSON.stringify(updated));
  } catch {
    // quota or private browsing mode
  }
}

export function getLastResult(): IncidentResult | null {
  if (typeof window === "undefined") return null;

  try {
    const raw = localStorage.getItem(LAST_RESULT_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function getCheckHistory(): HistoryEntry[] {
  if (typeof window === "undefined") return [];

  try {
    const raw = localStorage.getItem(HISTORY_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function clearCheckHistory(): void {
  if (typeof window === "undefined") return;

  try {
    localStorage.removeItem(LAST_RESULT_KEY);
    localStorage.removeItem(HISTORY_KEY);
  } catch {
    // ignore
  }
}
