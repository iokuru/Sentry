import { loadConfig } from './config';
import { getDeviceId } from './device';
import { getAuthToken } from './auth';

export interface PresignRequest {
  content_type: string;
  filename: string;
}

export class QuotaExceededError extends Error {
  constructor(message = 'Daily limit reached') {
    super(message);
    this.name = 'QuotaExceededError';
  }
}

export interface PresignResponse {
  upload_url: string;
  image_key: string;
}

export interface IntakeRequest {
  image_key?: string;
  user_context?: string;
  locale?: string;
  country?: string;
  senior_id?: string;
}

export interface IntakeResponse {
  incident_id: string;
  status: string;
}

export interface HotlineInfo {
  country: string;
  name: string;
  number: string;
  tel_uri: string;
}

export interface SeniorOutput {
  hold_off: string;
  safe_alternative: string;
  hotline: HotlineInfo;
}

export interface DetectedSignal {
  name: string;
  weight: number;
  description: string;
}

export interface Incident {
  id: string;
  senior_id?: string;
  image_key?: string;
  extracted_text?: string;
  user_context?: string;
  score: number;
  risk_level: 'low' | 'medium' | 'high';
  category: string;
  signals?: DetectedSignal[];
  output?: SeniorOutput;
  status: 'pending' | 'analyzed' | 'acknowledged' | 'resolved';
  thumbnail_url?: string;
  country?: string;
  created_at: string;
  updated_at: string;
}

async function getBaseUrl(): Promise<string> {
  const cfg = await loadConfig();
  return cfg.api_url ? cfg.api_url.replace(/\/$/, '') : '';
}

export async function presignUpload(contentType: string, filename: string): Promise<PresignResponse> {
  const base = await getBaseUrl();
  let res: Response;
  try {
    res = await fetch(`${base}/api/presign`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content_type: contentType, filename })
    });
  } catch {
    throw new Error('Connection failed, please check your network and retry');
  }

  if (res.status === 429) {
    throw new QuotaExceededError('Daily limit reached');
  }
  if (!res.ok) {
    throw new Error('Could not prepare upload');
  }

  return res.json();
}

export async function uploadScreenshot(uploadUrl: string, blob: Blob, contentType: string): Promise<void> {
  if (blob.size > 5 * 1024 * 1024) {
    throw new Error('Screenshot exceeds five megabyte limit');
  }

  const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp'];
  if (contentType && !validTypes.includes(contentType)) {
    throw new Error('Only PNG, JPG, and WebP images are supported');
  }

  let res: Response;
  try {
    res = await fetch(uploadUrl, {
      method: 'PUT',
      headers: { 'Content-Type': contentType },
      body: blob
    });
  } catch {
    throw new Error('Connection failed, please check your network and retry');
  }

  if (!res.ok) {
    throw new Error('Image upload failed, please retry');
  }
}

export async function submitIntake(req: IntakeRequest): Promise<IntakeResponse> {
  const base = await getBaseUrl();
  if (!req.senior_id) {
    req.senior_id = getDeviceId();
  }
  let res: Response;
  try {
    res = await fetch(`${base}/api/intake`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req)
    });
  } catch {
    throw new Error('Connection failed, please check your network and retry');
  }

  if (res.status === 429) {
    throw new QuotaExceededError('Daily limit reached');
  }
  if (!res.ok) {
    throw new Error('Analysis request failed, please retry');
  }

  return res.json();
}

export async function fetchIncident(id: string): Promise<Incident | null> {
  const base = await getBaseUrl();
  let res: Response;
  try {
    res = await fetch(`${base}/api/incidents?id=${encodeURIComponent(id)}`);
  } catch {
    return null;
  }

  if (!res.ok) {
    return null;
  }
  return res.json();
}

export async function pollIncident(id: string, maxAttempts = 30, intervalMs = 1500): Promise<Incident> {
  for (let i = 0; i < maxAttempts; i++) {
    const inc = await fetchIncident(id);
    if (inc && inc.output && inc.output.hold_off) {
      return inc;
    }
    await new Promise((r) => setTimeout(r, intervalMs));
  }
  throw new Error('Analysis timed out, please retry');
}

export async function listCaregiverIncidents(): Promise<Incident[]> {
  const base = await getBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}/api/caregiver/incidents`, {
    headers
  });

  if (res.status === 401) {
    throw new Error('Unauthorized');
  }
  if (!res.ok) {
    throw new Error('Failed to load incidents');
  }

  return res.json();
}

export async function updateCaregiverStatus(id: string, status: 'acknowledged' | 'resolved' | 'analyzed'): Promise<void> {
  const base = await getBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}/api/caregiver/status`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ incident_id: id, status })
  });

  if (res.status === 401) {
    throw new Error('Unauthorized');
  }
  if (!res.ok) {
    throw new Error('Status update failed');
  }
}

export async function pairDevice(code: string, deviceId: string): Promise<{ senior_id: string; senior_name: string }> {
  const base = await getBaseUrl();
  let res: Response;
  try {
    res = await fetch(`${base}/api/senior/pair`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code, device_id: deviceId })
    });
  } catch {
    throw new Error('Connection failed, please check your network and retry');
  }

  if (res.status === 404) {
    throw new Error('Invite code not found');
  }
  if (res.status === 410) {
    throw new Error('Invite code has expired');
  }
  if (res.status === 409) {
    throw new Error('Invite code has already been used');
  }
  if (!res.ok) {
    throw new Error('Pairing failed, please check your code and retry');
  }

  return res.json();
}

export async function createCaregiverInvite(seniorName?: string): Promise<{ code: string; expires_at: number }> {
  const base = await getBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}/api/caregiver/invites`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ senior_name: seniorName || '' })
  });

  if (!res.ok) {
    throw new Error('Failed to create invite code');
  }
  return res.json();
}

export interface CaregiverLinkInfo {
  caregiver_id: string;
  senior_id: string;
  senior_name: string;
  telegram_chat_id?: string;
  linked_at: string;
}

export async function fetchCaregiverLinks(): Promise<CaregiverLinkInfo[]> {
  const base = await getBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}/api/caregiver/links`, { headers });
  if (!res.ok) {
    return [];
  }
  return res.json();
}

export async function removeCaregiverLink(seniorId: string): Promise<void> {
  const base = await getBaseUrl();
  const token = await getAuthToken();
  const headers: Record<string, string> = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}/api/caregiver/links/${encodeURIComponent(seniorId)}`, {
    method: 'DELETE',
    headers
  });
  if (!res.ok) {
    throw new Error('Failed to remove link');
  }
}
