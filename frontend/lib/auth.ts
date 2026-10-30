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

  try {
    const authUrl = new URL("/", config.apiUrl);
    const res = await fetch(authUrl.toString(), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-amz-json-1.1",
        "X-Amz-Target": "AWSCognitoIdentityProviderService.InitiateAuth",
      },
      body: JSON.stringify({
        AuthFlow: "USER_PASSWORD_AUTH",
        ClientId: config.userPoolClientId,
        AuthParameters: { USERNAME: email.trim(), PASSWORD: pass },
      }),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      return { success: false, error: data.message || data.Message || "Sign in failed. Check your email and password." };
    }

    const token = data.AuthenticationResult?.IdToken;
    if (!token) return { success: false, error: "The sign-in service did not return a valid session." };

    const session: AuthSession = {
      token,
      email: email.trim(),
    };
    try {
      localStorage.setItem(TOKEN_KEY, session.token);
      localStorage.setItem(USER_KEY, session.email);
    } catch {
      return { success: false, error: "This browser could not save your sign-in. Enable site storage and try again." };
    }
    return { success: true, session };
  } catch {
    return { success: false, error: "Could not reach the sign-in service. Check your connection and try again." };
  }
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
