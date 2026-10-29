import { describe, it, expect, vi, beforeEach } from 'vitest';
import { submitIntake, TextractBlockedError } from './api';

describe('intake ocr handoff and text-only path', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('throws TextractBlockedError when server returns 422 TEXTRACT_UNAVAILABLE', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 422,
      ok: false,
      json: async () => ({
        code: 'TEXTRACT_UNAVAILABLE',
        error: 'textract unavailable: use browser ocr'
      })
    } as any);

    await expect(
      submitIntake({
        image_key: 'screenshots/scam.png',
        country: 'US',
        senior_id: 'dev-1'
      })
    ).rejects.toThrow(TextractBlockedError);
  });

  it('succeeds for text-only submission', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 202,
      ok: true,
      json: async () => ({
        incident_id: 'inc-text-only-1',
        status: 'pending'
      })
    } as any);

    const res = await submitIntake({
      user_context: 'Someone claiming to be from my bank asked for my PIN',
      country: 'US',
      senior_id: 'dev-1'
    });

    expect(res.incident_id).toBe('inc-text-only-1');
    expect(res.status).toBe('pending');
  });

  it('succeeds when client provides browser extracted text with image thumbnail', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      status: 202,
      ok: true,
      json: async () => ({
        incident_id: 'inc-browser-ocr-1',
        status: 'pending'
      })
    } as any);

    const res = await submitIntake({
      image_key: 'screenshots/scam.png',
      user_context: 'Urgent security notice: please verify password',
      country: 'US',
      senior_id: 'dev-1'
    });

    expect(res.incident_id).toBe('inc-browser-ocr-1');
    expect(res.status).toBe('pending');
  });
});
