import { loadAppConfig } from "./config";

const DEVICE_KEY = "sentry_device_id";

export function getOrCreateDeviceId(): string {
  if (typeof window === "undefined") return "server-id";

  let id = localStorage.getItem(DEVICE_KEY);
  if (!id) {
    id = "dev-" + Math.random().toString(36).substring(2, 10) + Math.random().toString(36).substring(2, 10);
    try {
      localStorage.setItem(DEVICE_KEY, id);
    } catch {
      // private browsing
    }
  }
  return id;
}

export async function redeemPairingCode(code: string): Promise<{ success: boolean; message?: string }> {
  const config = await loadAppConfig();
  const deviceId = getOrCreateDeviceId();

  const url = new URL("/api/senior/pair", config.apiUrl);
  const res = await fetch(url.toString(), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      code: code.trim().toUpperCase(),
      senior_id: deviceId,
    }),
  });

  if (!res.ok) {
    const errorText = await res.text();
    return { success: false, message: errorText || "Invalid or expired invite code." };
  }

  return { success: true };
}
