export interface AppConfig {
  apiUrl: string;
  region: string;
  userPoolId: string;
  userPoolClientId: string;
  ocrMode: string;
}

const defaultConfig: AppConfig = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL || "https://ni3bvlk951.execute-api.ap-south-1.amazonaws.com",
  region: "ap-south-1",
  userPoolId: "ap-south-1_NRUTetHDH",
  userPoolClientId: "4oo2js9m64ql6qbumihklgg4kv",
  ocrMode: "browser",
};

let cachedConfig: AppConfig | null = null;

export async function loadAppConfig(): Promise<AppConfig> {
  if (cachedConfig) {
    return cachedConfig;
  }

  try {
    const res = await fetch("/config.json", { cache: "no-store" });
    if (res.ok) {
      const data = await res.json() as Partial<AppConfig> & {
        api_url?: string;
        user_pool_id?: string;
        user_pool_client_id?: string;
        ocr_mode?: string;
      };
      cachedConfig = {
        apiUrl: data.apiUrl || data.api_url || defaultConfig.apiUrl,
        region: data.region || defaultConfig.region,
        userPoolId: data.userPoolId || data.user_pool_id || defaultConfig.userPoolId,
        userPoolClientId: data.userPoolClientId || data.user_pool_client_id || defaultConfig.userPoolClientId,
        ocrMode: data.ocrMode || data.ocr_mode || defaultConfig.ocrMode,
      };
      return cachedConfig!;
    }
  } catch {
    // fallback to defaults when running in local dev or offline
  }

  cachedConfig = defaultConfig;
  return cachedConfig;
}
