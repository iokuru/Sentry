import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { detectCountry } from './locale';

describe('detectCountry', () => {
  const originalNavigator = globalThis.navigator;

  afterEach(() => {
    Object.defineProperty(globalThis, 'navigator', {
      value: originalNavigator,
      configurable: true,
      writable: true
    });
  });

  it('detects UK locale', () => {
    Object.defineProperty(globalThis, 'navigator', {
      value: { language: 'en-GB' },
      configurable: true,
      writable: true
    });
    expect(detectCountry()).toBe('GB');
  });

  it('detects Australia locale', () => {
    Object.defineProperty(globalThis, 'navigator', {
      value: { language: 'en-AU' },
      configurable: true,
      writable: true
    });
    expect(detectCountry()).toBe('AU');
  });

  it('detects India locale', () => {
    Object.defineProperty(globalThis, 'navigator', {
      value: { language: 'hi-IN' },
      configurable: true,
      writable: true
    });
    expect(detectCountry()).toBe('IN');
  });

  it('falls back to US for generic locale', () => {
    Object.defineProperty(globalThis, 'navigator', {
      value: { language: 'en-US' },
      configurable: true,
      writable: true
    });
    expect(detectCountry()).toBe('US');
  });
});
