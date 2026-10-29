# sentry context

sentry helps older adults check suspicious screenshots or messages for fraud.

## repository structure
- cmd/: entrypoints for lambdas, builder, devserver, smoke, and operations
  - cmd/build: pure go arm64 bootstrap compiler and zip packager
  - cmd/devserver: local net/http server with in-memory zero aws fakes
  - cmd/preflight: validates inputs, aws credentials, and free-tier compliance
  - cmd/golive: automated end-to-end deployment, ssm config, and cdn invalidation
  - cmd/status: checks cloudfront, api health, dlq depth, and quota usage
  - cmd/teardown: cleans s3 buckets, ssm params, and destroys terraform stack
  - cmd/smoke: verification loop for local and live environments
  - cmd/costguard: enforces zero dollar aws free tier resource constraints
  - cmd/seed: generates realistic png screenshots and verifies category scoring
  - cmd/{presign,intake,analysis,caregiver_list,caregiver_status,digest,telegram_hook}
  - cmd/{caregiver_invite,senior_pair}
- internal/: core domain logic and platform adapters
  - internal/risk: deterministic weighted scam scoring and category classification
  - internal/hotline: official national fraud hotline numbers (us, uk, in, au)
  - internal/store: dynamodb persistence for incidents, pairings, and tokens
  - internal/quota: daily request throttle for textract and comprehend (100/day)
  - internal/telegram: net/http client, alert formatting, and webhook validator
  - internal/costguard: resource inspection against non-free infrastructure
  - internal/logging: structured json logging with request id extraction
  - internal/awsconfig: retry and timeout wrappers for aws clients
- frontend/: sveltekit static client with wcag aaa accessibility
- infra/: terraform modules (storage, api, data, auth, notify)
- fixtures/: generated scam screenshot images across all categories
- openapi.yaml: typed api contract

## operational commands
- make preflight: validate configuration, credentials, and constraints
- make go-live: complete one-touch build, terraform apply, and cdn deployment
- make status: check live health, dlq queue depth, and service quotas
- make teardown: empty storage and destroy all cloud resources safely
- make dev: run local devserver and frontend with zero aws dependencies
- make test: run go test suite and frontend vitest suite
- make smoke: run local end-to-end verification loop
- make smoke-live: run live verification loop against deployed aws stack
- make seed: generate fixture images and verify scoring engine
- make costguard: scan terraform files for non-free resource violations

## conventions
- go standard library first, no web frameworks, no orms
- aws sdk go v2 only, errors wrapped with fmt.Errorf
- pure functions for scoring signals with individual table tests
- sveltekit hand-written css with custom properties, minimum 60px touch targets
- zero dollar target: serverless pay-per-request and free tier only
- comments explain why, single line, lowercase, no full stop
