import { describe, it, expect } from "vitest";
import { metadata } from "../app/layout";

describe("landing page metadata and branding", () => {
  it("uses sentry brand in site title", () => {
    const title = String(metadata.title || "");
    expect(title).toContain("Sentry");
    expect(title).not.toContain("Take a moment");
  });

  it("describes the core value proposition", () => {
    const desc = String(metadata.description || "");
    expect(desc).toContain("Upload a screenshot or paste the text");
  });
});
