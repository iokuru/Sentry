# Contributing to Sentry

Thank you for your interest in contributing to Sentry. Sentry is designed to protect older adults and families from unexpected fraud with fast, calm, deterministic triage.

## Ways to Contribute

### 1. Suggesting or Adding New Scam Patterns

Scammers constantly change phrasing, domains, and coercive tactics. If you encounter a new fraud technique:

1. Open an issue describing the scam format, common keywords, and typical urgency hooks.
2. Add a sample JSON case to `testdata/cases/` containing:
   - `id`: Unique identifier (e.g., `utility-threat-clean`)
   - `input_text`: The realistic text body
   - `context`: Source context (e.g., `SMS notification`)
   - `source`: `text` or `screenshot`
   - `expected_level`: `high`, `medium`, or `low`
   - `expected_category`: `bank`, `delivery`, `utility`, `prize`, `tech support`, `relative in need`, `romance`, `tax`, or `other`
   - `required_signals`: Array of signal rule IDs that must trigger
3. Run `make cases` to verify the evaluation suite passes.

### 2. Refining Detection Rules

Detection logic lives in `internal/risk/`:
- `signals.go`: Heuristic pattern matching and keyword extraction.
- `scams.go`: Specialized threat recognizers (digital arrest, parcel hold fees, UPI/QR collection, etc.).
- `links.go`: Shortener and lookalike domain detection.
- `score.go`: Signal weights and threshold calculation.
- `category.go`: Category assignment and senior-friendly action wording.

When adjusting rules, keep these guidelines in mind:
- **Zero Hallucination**: All classification must be deterministic and reproducible.
- **Calm, Actionable Guidance**: Action copy must never panic the user, blame the victim, or use technical jargon.
- **Wording Quality**: Output must conform to the wording rules verified in `internal/risk/lint.go`.

### 3. Improving the Frontend Experience

Frontend code lives in `frontend/`:
- **Next.js with TypeScript**: Responsive layouts supporting phones (360px width) through desktop screens.
- **Accessibility**: High-contrast typography meeting WCAG AAA requirements. Large touch targets (minimum 44px to 60px).
- **Client-Side Privacy**: OCR execution should continue to run inside the browser worker whenever possible.

### 4. Reporting False Alarms or Missed Scams

If a legitimate message is flagged as high risk, or a dangerous scam is categorized as low risk:
- Open an issue with the anonymized message content.
- Note any specific links, sender names, or currency formats.

## Development Workflow

1. Fork and clone the repository.
2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Run the test suite:
   ```bash
   make test
   make cases
   ```
4. Verify code format and cleanliness:
   ```bash
   go vet ./...
   cd frontend && npm run check
   ```
5. Commit your changes with clear, concise messages.
6. Push to your fork and submit a Pull Request.

## Code of Conduct

Sentry is built to serve vulnerable populations. Be respectful, constructive, and empathetic in all interactions across issues and pull requests.
