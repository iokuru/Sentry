# sentry context

sentry helps older adults verify suspicious messages and screenshots.

## repository structure
- cmd/: entrypoints for lambda binaries, local dev server, and seed script
- internal/: core business logic, risk engine, hotlines, store, quota
  - internal/risk: deterministic scam scoring engine
  - internal/hotline: official fraud hotlines (us, uk, in, au)
  - internal/store: dynamodb incident and relationship persistence
  - internal/quota: daily textract and comprehend call limiter
- frontend/: sveltekit static client
- infra/: terraform root and modules (storage, api, data, auth, notify)
- openapi.yaml: typed http api contract

## commands
- make dev: run local backend server with stubbed aws services
- make test: run go tests and frontend unit tests
- make build: compile arm64 lambda binaries and frontend bundle
- make seed: populate database with sample incidents
- make plan: run terraform plan
- make apply: deploy cloud infrastructure

## conventions
- go standard library first, no web frameworks, no orm
- aws sdk go v2 only
- errors wrapped with context via fmt.Errorf
- pure functions for scoring signals with individual unit tests
- sveltekit hand-written css using custom properties, no tailwind
- zero cost target: serverless pay-per-request and free tier only
