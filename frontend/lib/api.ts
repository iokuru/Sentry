import { loadAppConfig } from "./config";

export interface PresignResponse {
  upload_url: string;
  object_key: string;
}

export interface IntakePayload {
  object_key?: string;
  context_text?: string;
  ocr_text?: string;
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
  seniorId?: string
): Promise<PresignResponse> {
  const config = await loadAppConfig();
  const url = new URL("/api/presign", config.apiUrl);

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      content_type: contentType,
      senior_id: seniorId,
    }),
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new ApiError(errorText || "Failed to create upload url", res.status);
  }

  return res.json();
}

export async function uploadImageToS3(uploadUrl: string, blob: Blob): Promise<void> {
  const res = await fetch(uploadUrl, {
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

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
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
  const url = new URL(`/api/incidents/${incidentId}`, config.apiUrl);

  const res = await fetch(url.toString(), {
    headers: { Accept: "application/json" },
    cache: "no-store",
  });

  if (!res.ok) {
    throw new ApiError("Failed to fetch incident result", res.status);
  }

  return res.json();
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
