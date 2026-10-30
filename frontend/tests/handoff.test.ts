import { describe, it, expect, vi, beforeEach } from "vitest";
import { submitIntake, ApiError } from "../lib/api";

describe("intake ocr handoff and text-only path", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("throws ApiError when server returns 422 TEXTRACT_UNAVAILABLE", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 422,
      ok: false,
      json: async () => ({
        code: "TEXTRACT_UNAVAILABLE",
        message: "textract unavailable: use browser ocr",
      }),
      text: async () => JSON.stringify({
        code: "TEXTRACT_UNAVAILABLE",
        message: "textract unavailable: use browser ocr",
      }),
    } as any);

    await expect(
      submitIntake({
        object_key: "screenshots/scam.png",
        country: "US",
        senior_id: "dev-1",
      })
    ).rejects.toThrow(ApiError);
  });

  it("succeeds for text-only submission", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 202,
      ok: true,
      json: async () => ({
        incident_id: "inc-text-only-1",
      }),
    } as any);

    const res = await submitIntake({
      context_text: "Someone claiming to be from my bank asked for my PIN",
      country: "US",
      senior_id: "dev-1",
    });

    expect(res.incident_id).toBe("inc-text-only-1");
  });

  it("succeeds when client provides browser extracted text", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 202,
      ok: true,
      json: async () => ({
        incident_id: "inc-browser-ocr-1",
      }),
    } as any);

    const res = await submitIntake({
      object_key: "screenshots/scam.png",
      ocr_text: "Urgent security notice: please verify password",
      country: "US",
      senior_id: "dev-1",
    });

    expect(res.incident_id).toBe("inc-browser-ocr-1");
  });
});
