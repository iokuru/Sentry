# sentry

sentry helps an older adult check a suspicious screenshot or message for fraud. they drop an image, optionally type what happened, and get exactly three things:
1. what to hold off on right now
2. a safe alternative action
3. a one tap button that dials the official fraud hotline for their country (us ftc, uk action fraud, india 1930, australia scamwatch)

when risk is high, linked caregivers receive an immediate alert by email and optional telegram. caregivers sign in to a portal that lists incidents as cards with severity, time, and status, with buttons to acknowledge or resolve.

## design decisions

- **why no bedrock:** bedrock charges per input and output token with zero free tier allowance. llms can also hallucinate phone numbers or panic older adults. sentry uses a deterministic scoring engine with weighted regex signals, ensuring zero cost, reproducible results, and sub-second classification.
- **why no orm:** orms add reflection overhead, memory overhead, and cold start latency to lambda binaries. sentry uses the official aws sdk go v2 dynamodb attributevalue marshaler directly.

interface design is adapted from asksafe home (https://github.com/v0-community/asksafe-home).

## exact free tier allowances

everything runs at zero usd on aws free tier or always-free allowances. textract and comprehend are the only services with limited free allowances, with a daily request cap of 100 requests per day enforced in code (falling back to rule-based classification when exhausted). an aws budget alarm is configured at 1 usd to alert on any unexpected spend.

### operational modes

sentry supports three operational modes and two frontend hosting options in `inputs.local.json`:

- **frontend host:** `frontend_host: "netlify"` publishes directly to netlify with exact origin cors restrictions on api gateway and uploads storage. `frontend_host: "cloudfront"` deploys to s3 behind cloudfront.

1. **auto mode (default):** `use_comprehend: "auto"`, `ocr_mode: "auto"`. probes textract and comprehend. uses cloud services when available; automatically falls back to deterministic rules if comprehend throws subscription errors, and automatically hands off to client-side browser ocr worker if textract is blocked. services used: s3, cloudfront, api gateway, lambda, dynamodb, sqs, sns, cognito, eventbridge scheduler, cloudwatch, ssm, plus textract and comprehend when accessible.
2. **cloud mode:** `use_comprehend: "on"`, `ocr_mode: "textract"`. enforces cloud textract ocr and comprehend analysis. preflight checks fail if your aws account lacks active subscriptions. services used: all thirteen services including textract and comprehend.
3. **zero-subscription free mode:** `use_comprehend: "off"`, `ocr_mode: "browser"`. completely disables textract and comprehend. terraform creates no textract or comprehend iam policies. text extraction runs entirely inside the user's browser using a self-hosted web worker with bundled language models, and classification runs using local deterministic rules. services used: eleven always-free or free tier services (s3, cloudfront, api gateway, lambda, dynamodb, sqs, sns, cognito, eventbridge scheduler, cloudwatch, ssm). neither textract nor comprehend is touched.

### services breakdown

- **s3:** 5 gb standard storage and 20,000 get requests. private uploads expire after 1 day via s3 lifecycle rule and are deleted immediately after analysis completes.
- **cloudfront:** 1 tb data transfer out per month and 10,000,000 requests (always free). serves static frontend with origin access control and security headers.
- **api gateway:** http api has 1,000,000 free requests per month. jwt authorizer on caregiver routes, throttled intake route for seniors.
- **lambda:** 1,000,000 free requests and 3,200,000 seconds of compute time per month (always free). compiled as single arm64 bootstrap binaries on provided.al2023.
- **dynamodb:** on-demand pay-per-request billing mode with zero provisioned capacity. 25 gb free storage. ttl of 90 days on incident records.
- **textract:** 1,000 pages per month free allowance for 3 months. a hard daily cap of 100 requests is enforced in code, falling back to rule-based analysis when reached.
- **comprehend:** 50,000 text units per month free allowance for 12 months. capped in code at 100 requests per day, falling back to local linguistic rules.
- **sqs:** 1,000,000 requests per month (always free). async analysis queue with dead-letter queue.
- **sns:** 1,000,000 publishes and 100,000 email deliveries per month (always free). paid sms is strictly prohibited.
- **cognito:** 50,000 monthly active users (always free). simple email and password sign in without hosted ui.
- **eventbridge scheduler:** 14,000,000 invocations per month (always free). runs daily digest email for unresolved incidents.
- **cloudwatch:** 7 day log retention across all functions and one alarm on dlq depth.
- **ssm parameter store:** standard tier parameters only (always free). stores telegram credentials and optional hotline registry overrides.
- **aws budgets:** 2 free budgets. a 1 usd budget alarm is configured to alert on unexpected usage.

## one-touch deployment

deploying sentry to a cold aws account requires only filling one inputs file and running two commands:

1. configure inputs:
   ```bash
   cp inputs.example.json inputs.local.json
   ```
   fill in your `aws_region`, `budget_alert_email`, `caregiver_email`, and optional telegram bot credentials.

2. verify preflight constraints:
   ```bash
   make preflight
   ```
   *(or `.\scripts\dev.ps1 preflight` on windows)*

3. execute automated go-live:
   ```bash
   make go-live
   ```
   *(or `.\scripts\dev.ps1 go-live` on windows)*

   this command compiles arm64 lambdas, builds the static frontend, applies terraform, stores ssm parameters, registers the telegram webhook, uploads assets to s3, invalidates cloudfront, and runs the live smoke test.

4. check deployment health anytime:
   ```bash
   make status
   ```
   *(or `.\scripts\dev.ps1 status` on windows)*

5. clean up all resources:
   ```bash
   make teardown
   ```
   *(or `.\scripts\dev.ps1 teardown` on windows)*

## local development

start the backend and frontend with zero aws credentials:

```bash
make dev
```
*(or `.\scripts\dev.ps1 dev` on windows)*

this launches:
- `cmd/devserver` on http://localhost:8080 with in-memory stores, fake ocr, fake comprehend, and fake telegram dispatcher.
- sveltekit frontend on http://localhost:5173 with proxy to backend.

run the test suite:

```bash
make test
```
*(or `.\scripts\dev.ps1 test` on windows)*

run the local verification loop:

```bash
make smoke
```

generate fixture images across all categories:

```bash
make seed
```
