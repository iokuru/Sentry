export interface AppConfig {
  api_url: string;
  user_pool_id: string;
  user_pool_client_id: string;
  region: string;
  ocr_mode?: string;
}

let cachedConfig: AppConfig | null = null;

export async function loadConfig(): Promise<AppConfig> {
  if (cachedConfig) {
    return cachedConfig;
  }

  try {
    const res = await fetch('/config.json');
    if (res.ok) {
      cachedConfig = await res.json();
      return cachedConfig!;
    }
  } catch {
    // fallback to defaults when config.json is not reachable
  }

  cachedConfig = {
    api_url: '',
    user_pool_id: 'local',
    user_pool_client_id: 'local',
    region: 'local'
  };
  return cachedConfig;
}

export function getConfig(): AppConfig | null {
  return cachedConfig;
}
