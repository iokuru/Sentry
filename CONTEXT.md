# sentry context

sentry helps older adults check suspicious screenshots or messages for fraud.

## repository structure
- cmd/: entrypoints for lambdas, builder, devserver, smoke, and deploy tools
  - cmd/build: pure go arm64 bootstrap compiler and zip packager
  - cmd/devserver: local net/http server with in-memory zero aws fakes
  - cmd/smoke: automated verification loop for local and live environments
  - cmd/deploy: build, plan, apply, config sync, and cdn invalidation pipeline
  - cmd/costguard: enforces zero dollar aws free tier resource constraints
  - cmd/seed: generates realistic png screenshots and verifies category scoring
  - cmd/{presign,intake,analysis,caregiver_list,caregiver_status,digest,telegram_hook}
- internal/: core domain logic and platform adapters
  - internal/risk: deterministic weighted scam scoring and category classification
  - internal/hotline: official national fraud hotline numbers (us, uk, in, au)
  - internal/store: dynamodb persistence for incidents and relationships
  - internal/quota: daily request throttle for textract and comprehend
  - internal/telegram: net/http client, alert formatting, and webhook validator
  - internal/costguard: resource inspection against non-free infrastructure
- frontend/: sveltekit static client with wcag aaa accessibility
- infra/: terraform modules (storage, api, data, auth, notify)
- fixtures/: generated scam screenshot images across all categories
- openapi.yaml: typed api contract

## commands
- make dev: run local devserver and frontend with zero aws dependencies
- make test: run go test suite and frontend vitest suite
- make build: package arm64 lambda archives and build frontend bundle
- make smoke: run local end to end verification loop
- make smoke-live: run live verification loop against deployed aws stack
- make seed: generate fixture images and verify scoring engine
- make costguard: scan terraform files for non-free resource violations
- make deploy: run full deployment pipeline to aws
- make telegram-hook: register telegram webhook with deployed api url

## conventions
- go standard library first, no web frameworks, no orms
- aws sdk go v2 only, errors wrapped with fmt.Errorf
- pure functions for scoring signals with individual table tests
- sveltekit hand-written css with custom properties, minimum 60px touch targets
- zero dollar target: serverless pay-per-request and free tier only
- comments explain why, single line, lowercase, no full stop
