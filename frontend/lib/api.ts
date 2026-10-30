import { loadAppConfig } from "./config";
import { getOrCreateDeviceId } from "./device";

export interface PresignResponse {
  upload_url: string;
  image_key: string;
}

export interface IntakePayload {
  image_key?: string;
  user_context?: string;
  locale?: string;
  senior_id?: string;
  country?: string;
}

export interface IntakeResponse {
  incident_id: string;
}

export interface IncidentResult {
  incident_id: string;
  status: "NEW" | "ACKNOWLEDGED" | "RESOLVED" | "ANALYZING";
  severity: "HIGH" | "MEDIUM" | "LOW";
  signals: string[];
  hold_off: string;
  safe_alternative: string;
  hotline: string;
  hotline_name?: string;
  hotline_phone?: string;
  hotline_url?: string;
  country: string;
  created_at: string;
  error?: string;
}

export function normalizeIncident(raw: any, fallbackId = ""): IncidentResult {
  const output = raw.output || {};
  const hotline = output.hotline || {};
  const rawStatus = String(raw.status || "ANALYZING").toUpperCase();
  const risk = String(raw.severity || raw.risk_level || "LOW").toUpperCase();
  const signals = Array.isArray(raw.signals)
    ? raw.signals.map((signal: string | { description?: string; name?: string; id?: string }) =>
        typeof signal === "string" ? signal : signal.description || signal.name || signal.id || ""
      ).filter(Boolean)
    : [];
  const backendHoldOff = raw.hold_off || output.hold_off || "";
  const holdOff = signals.length > 0 && /nothing urgent found/i.test(backendHoldOff)
    ? `A warning sign was found: ${signals[0]}. Pause before clicking or paying, and verify the request through the organization’s official website or phone number.`
    : backendHoldOff;

  return {
    incident_id: raw.incident_id || raw.id || fallbackId,
    status: rawStatus === "ANALYZED" || rawStatus === "PENDING" || rawStatus === "NEW"
      ? "NEW"
      : rawStatus === "ACKNOWLEDGED" || rawStatus === "RESOLVED"
        ? rawStatus
        : "ANALYZING",
    severity: risk === "HIGH" || risk === "MEDIUM" ? risk : "LOW",
    signals,
    hold_off: holdOff,
    safe_alternative: raw.safe_alternative || output.safe_alternative || "",
    hotline: raw.hotline || hotline.number || "",
    hotline_name: raw.hotline_name || hotline.name,
    hotline_phone: raw.hotline_phone || hotline.number,
    hotline_url: raw.hotline_url || hotline.url,
    country: raw.country || hotline.country || "",
    created_at: raw.created_at || "",
    error: raw.error,
  };
}

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(message: string, status: number, code?: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

export async function requestPresignedUrl(
  contentType: string,
  filename?: string
): Promise<PresignResponse> {
  const config = await loadAppConfig();
  const url = new URL("/api/presign", config.apiUrl);

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      content_type: contentType,
      filename,
    }),
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new ApiError(errorText || "Failed to create upload url", res.status);
  }

  return res.json();
}

export async function uploadImageToS3(uploadUrl: string, blob: Blob): Promise<void> {
  const config = await loadAppConfig();
  const resolvedUploadUrl = new URL(uploadUrl, config.apiUrl).toString();
  const res = await fetch(resolvedUploadUrl, {
    method: "PUT",
    headers: { "Content-Type": blob.type || "image/jpeg" },
    body: blob,
  });

  if (!res.ok) {
    throw new ApiError("Failed to upload image data to storage", res.status);
  }
}

export async function submitIntake(payload: IntakePayload): Promise<IntakeResponse> {
  const config = await loadAppConfig();
  const url = new URL("/api/intake", config.apiUrl);

  const finalPayload: IntakePayload = {
    ...payload,
    senior_id: payload.senior_id || getOrCreateDeviceId(),
  };

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(finalPayload),
  });

  if (!res.ok) {
    let errorCode: string | undefined;
    let errorMsg = "Failed to submit message for review";
    try {
      const data = await res.json();
      errorCode = data.code || data.error;
      errorMsg = data.message || errorMsg;
    } catch {
      // not json
    }
    throw new ApiError(errorMsg, res.status, errorCode);
  }

  return res.json();
}

export async function fetchIncident(incidentId: string): Promise<IncidentResult> {
  const config = await loadAppConfig();
  const url = new URL("/api/incidents", config.apiUrl);
  url.searchParams.set("id", incidentId);

  const res = await fetch(url.toString(), {
    headers: { Accept: "application/json" },
    cache: "no-store",
  });

  if (!res.ok) {
    throw new ApiError("Failed to fetch incident result", res.status);
  }

  return normalizeIncident(await res.json(), incidentId);
}

export async function pollIncidentResult(
  incidentId: string,
  maxAttempts = 15,
  intervalMs = 1000
): Promise<IncidentResult> {
  for (let attempt = 0; attempt < maxAttempts; attempt++) {
    try {
      const result = await fetchIncident(incidentId);
      if (result.status !== "ANALYZING" && (result.hold_off || result.safe_alternative)) {
        return result;
      }
    } catch {
      // transient network blip, continue polling
    }
    await new Promise((resolve) => setTimeout(resolve, intervalMs));
  }

  throw new Error("Analysis took longer than usual. Please check your connection or try again.");
}
