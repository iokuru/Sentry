import { describe, it, expect, afterEach } from "vitest";
import { detectBrowserCountry } from "../components/country-picker";

describe("detectBrowserCountry", () => {
  const originalNavigator = globalThis.navigator;

  afterEach(() => {
    Object.defineProperty(globalThis, "navigator", {
      value: originalNavigator,
      configurable: true,
      writable: true,
    });
  });

  it("detects UK locale", () => {
    Object.defineProperty(globalThis, "navigator", {
      value: { languages: ["en-GB"] },
      configurable: true,
      writable: true,
    });
    expect(detectBrowserCountry()).toBe("UK");
  });

  it("detects Australia locale", () => {
    Object.defineProperty(globalThis, "navigator", {
      value: { languages: ["en-AU"] },
      configurable: true,
      writable: true,
    });
    expect(detectBrowserCountry()).toBe("AU");
  });

  it("detects India locale", () => {
    Object.defineProperty(globalThis, "navigator", {
      value: { languages: ["hi-IN"] },
      configurable: true,
      writable: true,
    });
    expect(detectBrowserCountry()).toBe("IN");
  });

  it("falls back to US for generic locale", () => {
    Object.defineProperty(globalThis, "navigator", {
      value: { languages: ["en-US"] },
      configurable: true,
      writable: true,
    });
    expect(detectBrowserCountry()).toBe("US");
  });
});
