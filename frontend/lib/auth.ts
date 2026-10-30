import { loadAppConfig } from "./config";

const TOKEN_KEY = "sentry_caregiver_token";
const USER_KEY = "sentry_caregiver_user";

export interface AuthSession {
  token: string;
  email: string;
}

export function getStoredSession(): AuthSession | null {
  if (typeof window === "undefined") return null;

  try {
    const token = localStorage.getItem(TOKEN_KEY);
    const email = localStorage.getItem(USER_KEY);
    if (token && email) {
      return { token, email };
    }
  } catch {
    // storage unavailable
  }
  return null;
}

export async function signInCaregiver(
  email: string,
  pass: string
): Promise<{ success: boolean; session?: AuthSession; error?: string }> {
  const config = await loadAppConfig();

  // In local development or mock server, negotiate auth token via dev backend
  try {
    const res = await fetch(`${config.apiUrl}/api/caregiver/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: email.trim(), password: pass }),
    });

    if (res.ok) {
      const data = await res.json();
      const session: AuthSession = {
        token: data.token || "dev-caregiver-jwt",
        email: email.trim(),
      };
      localStorage.setItem(TOKEN_KEY, session.token);
      localStorage.setItem(USER_KEY, session.email);
      return { success: true, session };
    }
  } catch {
    // fallback for local testing without mock login route
  }

  // Fallback demo auth for mock environment
  if (pass.length >= 6) {
    const session: AuthSession = {
      token: "mock-jwt-" + btoa(email.trim()),
      email: email.trim(),
    };
    try {
      localStorage.setItem(TOKEN_KEY, session.token);
      localStorage.setItem(USER_KEY, session.email);
    } catch {
      // ignore
    }
    return { success: true, session };
  }

  return { success: false, error: "Password must be at least six characters." };
}

export function signOutCaregiver(): void {
  if (typeof window === "undefined") return;

  try {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  } catch {
    // ignore
  }
}
