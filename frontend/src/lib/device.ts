const DEVICE_ID_KEY = 'sentry_device_id';
const PAIRED_KEY = 'sentry_paired';
const SENIOR_NAME_KEY = 'sentry_senior_name';

export function getDeviceId(): string {
  if (typeof window === 'undefined' || !window.localStorage) {
    return 'browser-device';
  }

  let deviceId = localStorage.getItem(DEVICE_ID_KEY);
  if (!deviceId) {
    if (typeof crypto !== 'undefined' && crypto.randomUUID) {
      deviceId = crypto.randomUUID();
    } else {
      deviceId = 'dev-' + Math.random().toString(36).substring(2, 15);
    }
    localStorage.setItem(DEVICE_ID_KEY, deviceId);
  }
  return deviceId;
}

export function isPaired(): boolean {
  if (typeof window === 'undefined' || !window.localStorage) {
    return false;
  }
  return localStorage.getItem(PAIRED_KEY) === 'true';
}

export function setPaired(seniorName?: string): void {
  if (typeof window === 'undefined' || !window.localStorage) {
    return;
  }
  localStorage.setItem(PAIRED_KEY, 'true');
  if (seniorName) {
    localStorage.setItem(SENIOR_NAME_KEY, seniorName);
  }
}

export function getSeniorName(): string {
  if (typeof window === 'undefined' || !window.localStorage) {
    return '';
  }
  return localStorage.getItem(SENIOR_NAME_KEY) || '';
}

export function clearPairing(): void {
  if (typeof window === 'undefined' || !window.localStorage) {
    return;
  }
  localStorage.removeItem(PAIRED_KEY);
  localStorage.removeItem(SENIOR_NAME_KEY);
}
