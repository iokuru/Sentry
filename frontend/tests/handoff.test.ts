import { describe, it, expect, vi, beforeEach } from "vitest";
import { submitIntake, ApiError, normalizeIncident, requestPresignedUrl, uploadImageToS3, fetchIncident } from "../lib/api";

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
        image_key: "screenshots/scam.png",
        country: "US",
        senior_id: "dev-1",
      })
    ).rejects.toThrow(ApiError);
    const requestBody = vi.mocked(global.fetch).mock.calls.at(-1)?.[1]?.body as string;
    expect(JSON.parse(requestBody)).toMatchObject({ image_key: "screenshots/scam.png", country: "US" });
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
      user_context: "Someone claiming to be from my bank asked for my PIN",
      country: "US",
      senior_id: "dev-1",
    });

    expect(res.incident_id).toBe("inc-text-only-1");
    const requestBody = vi.mocked(global.fetch).mock.calls.at(-1)?.[1]?.body as string;
    expect(JSON.parse(requestBody)).toMatchObject({ user_context: "Someone claiming to be from my bank asked for my PIN", country: "US" });
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
      image_key: "screenshots/scam.png",
      user_context: "Urgent security notice: please verify password",
      country: "US",
      senior_id: "dev-1",
    });

    expect(res.incident_id).toBe("inc-browser-ocr-1");
  });

  it("normalizes the backend incident shape for the result screen", () => {
    const incident = normalizeIncident({
      id: "inc-42",
      status: "analyzed",
      risk_level: "high",
      country: "US",
      signals: [{ description: "Unexpected payment request" }],
      output: {
        hold_off: "Pause before paying.",
        safe_alternative: "Open the carrier app directly.",
        hotline: { country: "US", name: "FTC", number: "1-877-382-4357", url: "https://reportfraud.ftc.gov" },
      },
    });

    expect(incident).toMatchObject({
      incident_id: "inc-42",
      status: "NEW",
      severity: "HIGH",
      signals: ["Unexpected payment request"],
      hold_off: "Pause before paying.",
      safe_alternative: "Open the carrier app directly.",
      hotline_name: "FTC",
      hotline_phone: "1-877-382-4357",
      hotline_url: "https://reportfraud.ftc.gov",
    });
  });

  it("does not suppress a detected signal behind a low-risk default message", () => {
    const incident = normalizeIncident({
      id: "inc-low-with-signal",
      status: "analyzed",
      risk_level: "low",
      signals: [{ description: "uses a hidden shortened link to obscure destination" }],
      output: { hold_off: "Nothing urgent found. Take your time.", safe_alternative: "Verify the sender." },
    });

    expect(incident.hold_off).toContain("A warning sign was found");
    expect(incident.hold_off).toContain("shortened link");
  });

  it("uses the backend incident query route and normalizes its response", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        id: "inc-42",
        status: "analyzed",
        risk_level: "medium",
        output: { hold_off: "Pause before paying.", safe_alternative: "Use the official app.", hotline: { country: "US" } },
      }),
    } as any);

    const result = await fetchIncident("inc-42");

    expect(vi.mocked(global.fetch).mock.calls[0][0]).toBe("http://localhost:8080/api/incidents?id=inc-42");
    expect(result).toMatchObject({ incident_id: "inc-42", status: "NEW", severity: "MEDIUM", hold_off: "Pause before paying." });
  });

  it("uses the image_key presign contract and resolves local upload paths", async () => {
    const mockFetch = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => ({ upload_url: "/upload/test.jpg", image_key: "screenshots/test.jpg" }) } as any)
      .mockResolvedValueOnce({ ok: true } as any);
    global.fetch = mockFetch;

    const presign = await requestPresignedUrl("image/jpeg", "test.jpg");
    await uploadImageToS3(presign.upload_url, new Blob(["image"], { type: "image/jpeg" }));

    expect(mockFetch.mock.calls[0][1]?.body).toContain('"filename":"test.jpg"');
    expect(mockFetch.mock.calls[1][0]).toBe("http://localhost:8080/upload/test.jpg");
  });
});
