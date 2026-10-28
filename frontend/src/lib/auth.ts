import { loadConfig } from './config';

interface AuthSession {
  accessToken: string;
  idToken: string;
  refreshToken?: string;
  expiresAt: number;
}

const STORAGE_KEY = 'sentry_caregiver_session';

let memorySession: AuthSession | null = null;

export function initAuthFromStorage(): void {
  if (typeof window === 'undefined') return;
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed: AuthSession = JSON.parse(raw);
      if (parsed && parsed.expiresAt > Date.now()) {
        memorySession = parsed;
      } else {
        sessionStorage.removeItem(STORAGE_KEY);
      }
    }
  } catch {
    // ignore storage error
  }
}

export function isAuthenticated(): boolean {
  if (!memorySession) {
    initAuthFromStorage();
  }
  return !!memorySession && memorySession.expiresAt > Date.now();
}

export async function signIn(email: string, pass: string): Promise<void> {
  const cfg = await loadConfig();

  // local mock auth for dev server without aws
  if (cfg.user_pool_id === 'local' || cfg.region === 'local') {
    let res: Response;
    try {
      res = await fetch('/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-amz-json-1.1',
          'X-Amz-Target': 'AWSCognitoIdentityProviderService.InitiateAuth'
        },
        body: JSON.stringify({ email, pass })
      });
    } catch {
      // fallback local token
      saveSession({
        accessToken: 'dev-access-token',
        idToken: 'dev-id-token',
        refreshToken: 'dev-refresh-token',
        expiresAt: Date.now() + 3600 * 1000
      });
      return;
    }

    if (res.ok) {
      const data = await res.json();
      const auth = data.AuthenticationResult;
      saveSession({
        accessToken: auth.AccessToken,
        idToken: auth.IdToken,
        refreshToken: auth.RefreshToken,
        expiresAt: Date.now() + (auth.ExpiresIn || 3600) * 1000
      });
      return;
    }
  }

  // cognito identity provider http api
  const endpoint = `https://cognito-idp.${cfg.region}.amazonaws.com/`;
  const body = {
    AuthFlow: 'USER_PASSWORD_AUTH',
    ClientId: cfg.user_pool_client_id,
    AuthParameters: {
      USERNAME: email,
      PASSWORD: pass
    }
  };

  let res: Response;
  try {
    res = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-amz-json-1.1',
        'X-Amz-Target': 'AWSCognitoIdentityProviderService.InitiateAuth'
      },
      body: JSON.stringify(body)
    });
  } catch {
    throw new Error('Connection failed, please check your network and retry');
  }

  if (!res.ok) {
    const errData = await res.json().catch(() => ({}));
    const message = errData.__type === 'NotAuthorizedException'
      ? 'Incorrect email address or password'
      : 'Sign in failed, please check credentials and retry';
    throw new Error(message);
  }

  const data = await res.json();
  const auth = data.AuthenticationResult;
  saveSession({
    accessToken: auth.AccessToken,
    idToken: auth.IdToken,
    refreshToken: auth.RefreshToken,
    expiresAt: Date.now() + (auth.ExpiresIn || 3600) * 1000
  });
}

export async function refreshAuthSession(): Promise<boolean> {
  if (!memorySession || !memorySession.refreshToken) {
    return false;
  }

  const cfg = await loadConfig();
  if (cfg.user_pool_id === 'local' || cfg.region === 'local') {
    memorySession.expiresAt = Date.now() + 3600 * 1000;
    saveSession(memorySession);
    return true;
  }

  const endpoint = `https://cognito-idp.${cfg.region}.amazonaws.com/`;
  const body = {
    AuthFlow: 'REFRESH_TOKEN_AUTH',
    ClientId: cfg.user_pool_client_id,
    AuthParameters: {
      REFRESH_TOKEN: memorySession.refreshToken
    }
  };

  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-amz-json-1.1',
        'X-Amz-Target': 'AWSCognitoIdentityProviderService.InitiateAuth'
      },
      body: JSON.stringify(body)
    });

    if (!res.ok) {
      signOut();
      return false;
    }

    const data = await res.json();
    const auth = data.AuthenticationResult;
    saveSession({
      accessToken: auth.AccessToken,
      idToken: auth.IdToken,
      refreshToken: auth.RefreshToken || memorySession.refreshToken,
      expiresAt: Date.now() + (auth.ExpiresIn || 3600) * 1000
    });
    return true;
  } catch {
    return false;
  }
}

export async function getAuthToken(): Promise<string | null> {
  if (!memorySession) {
    initAuthFromStorage();
  }
  if (!memorySession) {
    return null;
  }

  // refresh token if within 60 seconds of expiration
  if (memorySession.expiresAt - Date.now() < 60000) {
    await refreshAuthSession();
  }

  return memorySession ? memorySession.idToken : null;
}

export function signOut(): void {
  memorySession = null;
  if (typeof window !== 'undefined') {
    try {
      sessionStorage.removeItem(STORAGE_KEY);
    } catch {
      // ignore
    }
  }
}

function saveSession(session: AuthSession): void {
  memorySession = session;
  if (typeof window !== 'undefined') {
    try {
      sessionStorage.setItem(STORAGE_KEY, JSON.stringify(session));
    } catch {
      // ignore
    }
  }
}
