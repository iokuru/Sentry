import React from "react";
import { Globe } from "lucide-react";

export interface CountryPickerProps {
  selectedCountry: string;
  onCountryChange: (country: string) => void;
}

const supportedCountries = [
  { code: "US", label: "United States" },
  { code: "UK", label: "United Kingdom" },
  { code: "IN", label: "India" },
  { code: "AU", label: "Australia" },
];

export function detectBrowserCountry(): string {
  if (typeof navigator === "undefined") return "US";

  const languages = navigator.languages || [navigator.language || ""];
  for (const lang of languages) {
    if (!lang) continue;
    const parts = lang.split("-");
    if (parts.length > 1) {
      const region = parts[parts.length - 1].toUpperCase();
      if (supportedCountries.some((c) => c.code === region)) {
        return region;
      }
      if (region === "GB") return "UK";
    }
  }

  return "US";
}

export function CountryPicker({
  selectedCountry,
  onCountryChange,
}: CountryPickerProps) {
  return (
    <div className="country-wrapper">
      <label htmlFor="country-select" className="country-label">
        <Globe size={18} aria-hidden="true" />
        <span>Country:</span>
      </label>

      <select
        id="country-select"
        value={selectedCountry}
        onChange={(e) => onCountryChange(e.target.value)}
        className="country-select"
        aria-label="Select your country for local hotline"
      >
        {supportedCountries.map((c) => (
          <option key={c.code} value={c.code}>
            {c.label}
          </option>
        ))}
      </select>

      <style jsx>{`
        .country-wrapper {
          display: inline-flex;
          align-items: center;
          gap: 8px;
          background-color: var(--bg-card);
          border: 1px solid var(--border-light);
          padding: 6px 12px;
          border-radius: var(--radius-md);
        }

        .country-label {
          display: inline-flex;
          align-items: center;
          gap: 6px;
          font-size: 0.95rem;
          font-weight: 600;
          color: var(--text-muted);
        }

        .country-select {
          background: transparent;
          border: none;
          font-size: 0.95rem;
          font-weight: 700;
          color: var(--text-main);
          cursor: pointer;
          outline: none;
        }

        .country-select:focus-visible {
          outline: 2px solid var(--color-primary);
        }
      `}</style>
    </div>
  );
}
