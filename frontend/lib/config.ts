export interface AppConfig {
  apiUrl: string;
  region: string;
  userPoolId: string;
  userPoolClientId: string;
  ocrMode: string;
}

const defaultConfig: AppConfig = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  region: "us-east-1",
  userPoolId: "",
  userPoolClientId: "",
  ocrMode: "auto",
};

let cachedConfig: AppConfig | null = null;

export async function loadAppConfig(): Promise<AppConfig> {
  if (cachedConfig) {
    return cachedConfig;
  }

  try {
    const res = await fetch("/config.json", { cache: "no-store" });
    if (res.ok) {
      const data = await res.json();
      cachedConfig = {
        ...defaultConfig,
        ...data,
      };
      return cachedConfig!;
    }
  } catch {
    // fallback to defaults when running in local dev or offline
  }

  cachedConfig = defaultConfig;
  return cachedConfig;
}
