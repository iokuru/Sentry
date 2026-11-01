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

function getCognitoEndpoint(region?: string): string {
  const reg = region || "ap-south-1";
  return `https://cognito-idp.${reg}.amazonaws.com/`;
}

export async function signInCaregiver(
  email: string,
  pass: string
): Promise<{ success: boolean; session?: AuthSession; error?: string }> {
  const config = await loadAppConfig();

  try {
    const res = await fetch(getCognitoEndpoint(config.region), {
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

export async function signUpCaregiver(
  email: string,
  pass: string
): Promise<{ success: boolean; userSub?: string; needsConfirmation?: boolean; error?: string }> {
  const config = await loadAppConfig();

  try {
    const res = await fetch(getCognitoEndpoint(config.region), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-amz-json-1.1",
        "X-Amz-Target": "AWSCognitoIdentityProviderService.SignUp",
      },
      body: JSON.stringify({
        ClientId: config.userPoolClientId,
        Username: email.trim(),
        Password: pass,
        UserAttributes: [{ Name: "email", Value: email.trim() }],
      }),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      return {
        success: false,
        error: data.message || data.Message || "Sign up failed. Password must be at least 8 characters with numbers and special symbols.",
      };
    }

    return {
      success: true,
      userSub: data.UserSub,
      needsConfirmation: !data.UserConfirmed,
    };
  } catch {
    return { success: false, error: "Could not reach the registration service. Check your connection and try again." };
  }
}

export async function confirmSignUpCaregiver(
  email: string,
  code: string
): Promise<{ success: boolean; error?: string }> {
  const config = await loadAppConfig();

  try {
    const res = await fetch(getCognitoEndpoint(config.region), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-amz-json-1.1",
        "X-Amz-Target": "AWSCognitoIdentityProviderService.ConfirmSignUp",
      },
      body: JSON.stringify({
        ClientId: config.userPoolClientId,
        Username: email.trim(),
        ConfirmationCode: code.trim(),
      }),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      return { success: false, error: data.message || data.Message || "Verification failed. Check the verification code sent to your email." };
    }

    return { success: true };
  } catch {
    return { success: false, error: "Could not reach the verification service. Check your connection and try again." };
  }
}

export async function resendConfirmationCode(
  email: string
): Promise<{ success: boolean; error?: string }> {
  const config = await loadAppConfig();

  try {
    const res = await fetch(getCognitoEndpoint(config.region), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-amz-json-1.1",
        "X-Amz-Target": "AWSCognitoIdentityProviderService.ResendConfirmationCode",
      },
      body: JSON.stringify({
        ClientId: config.userPoolClientId,
        Username: email.trim(),
      }),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      return { success: false, error: data.message || data.Message || "Could not resend verification code." };
    }

    return { success: true };
  } catch {
    return { success: false, error: "Could not reach the authentication service. Check your connection and try again." };
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
