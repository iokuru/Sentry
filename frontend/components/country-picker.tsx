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
        <Globe size={15} aria-hidden="true" />
        <span className="sr-only">Country:</span>
      </label>

      <select
        id="country-select"
        value={selectedCountry}
        onChange={(e) => onCountryChange(e.target.value)}
        className="country-select"
        aria-label="Select your country for local fraud hotline"
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
          gap: 6px;
          background: var(--surface-raised);
          border: 1px solid var(--border-light);
          min-height: 38px;
          padding: 4px 12px;
          border-radius: var(--radius-full);
          transition: border-color 140ms ease;
        }

        .country-wrapper:hover {
          border-color: var(--text-muted);
        }

        .country-label {
          display: inline-flex;
          align-items: center;
          color: var(--text-muted);
        }

        .country-select {
          background: transparent;
          border: none;
          font-size: 0.86rem;
          font-weight: 650;
          color: var(--text-main);
          cursor: pointer;
          outline: none;
          color-scheme: inherit;
        }

        .country-select:focus-visible {
          outline: none;
        }

        .country-select option {
          background-color: var(--bg-card);
          color: var(--text-main);
          font-weight: 500;
          padding: 8px 12px;
        }
      `}</style>
    </div>
  );
}
