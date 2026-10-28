# sentry

sentry helps an older adult check a suspicious screenshot or message. they drop an image, optionally type what happened, and get exactly three things:
1. what to hold off on right now
2. a safe alternative action
3. a one tap button that dials the official fraud hotline for their country (us ftc, uk action fraud, india 1930, australia scamwatch)

when risk is high, a caregiver receives an immediate alert via email and optional telegram. caregivers sign in to a portal that lists incidents as cards with severity, time, and status, with buttons to acknowledge or resolve.

## design decisions

- **why no bedrock:** bedrock charges per input/output token with zero free tier allowance. llms can also hallucinate telephone numbers, panic older adults, or output inconsistent advice. sentry uses a deterministic scoring engine with weighted regex signals, ensuring free, deterministic, and instant classification.
- **why no react:** react packages bring large client bundles and component library bloat. sentry uses sveltekit with static adapter and hand-written css custom properties, producing minimal vanilla javascript, high contrast wcag aaa accessibility, and minimum 60px touch targets.
- **why no orm:** orms add reflection overhead, memory overhead, and cold start latency to go lambda binaries. sentry uses the official aws sdk go v2 dynamodb attributevalue marshaler directly.

## exact free tier allowances

every service runs at zero usd on aws free tier or always-free allowances:
- **s3:** 5 gb standard storage and 20,000 get requests. private uploads expire after 1 day via s3 lifecycle rule.
- **cloudfront:** 1 tb data transfer out per month and 10,000,000 requests (always free).
- **api gateway:** http api has 1,000,000 free requests per month.
- **lambda:** 1,000,000 free requests and 3,200,000 seconds of compute time per month (always free). compiled as single arm64 bootstrap binaries on provided.al2023.
- **dynamodb:** on-demand pay-per-request billing mode with zero provisioned capacity. 25 gb free storage.
- **textract:** 1,000 pages per month free allowance for 3 months. a hard daily cap of 100 requests is enforced in code, falling back to rule-based analysis when reached.
- **comprehend:** 50,000 text units per month free allowance for 12 months. capped in code at 100 requests per day, falling back to local linguistic rules.
- **sqs:** 1,000,000 requests per month (always free). async analysis queue with dead-letter queue.
- **sns:** 1,000,000 publishes and 100,000 email deliveries per month (always free). paid sms is strictly prohibited.
- **cognito:** 50,000 monthly active users (always free). simple email and password sign in without hosted ui.
- **eventbridge scheduler:** 14,000,000 invocations per month (always free). runs daily digest email.
- **cloudwatch:** 7 day log retention across all functions and one alarm on dlq depth.
- **ssm parameter store:** standard tier parameters only (always free).
- **aws budgets:** 2 free budgets. a 1 usd budget alarm is configured to alert on unexpected usage.

## local development

start the backend and frontend with zero aws credentials:

```bash
make dev
```

this launches:
- `cmd/devserver` on http://localhost:8080 with in-memory stores, fake ocr, fake comprehend, and fake telegram dispatcher.
- sveltekit frontend on http://localhost:5173 with proxy to backend.

run the test suite:

```bash
make test
```

run the end to end smoke test loop:

```bash
make smoke
```

seed sample incidents across all eight categories:

```bash
make seed
```

## deployment

deploy to aws in one command:

```bash
make deploy
```

this compiles arm64 lambda packages, builds the frontend bundle, verifies cost guard compliance, applies terraform, writes the runtime `/config.json`, uploads assets to s3, and invalidates the cloudfront cache.

## telegram alerts setup

telegram alerts are optional and secondary. when unconfigured, telegram is silently disabled and email alerts still fire.

1. message `@BotFather` on telegram to create a bot: `/newbot`. copy the bot token.
2. start a chat with your bot and send a message. get your chat id using `@userinfobot` or `https://api.telegram.org/bot<TOKEN>/getUpdates`.
3. set ssm parameters:
   ```bash
   aws ssm put-parameter --name "/sentry/telegram_token" --value "<YOUR_BOT_TOKEN>" --type "SecureString" --overwrite
   aws ssm put-parameter --name "/sentry/telegram_chat_id" --value "<YOUR_CHAT_ID>" --type "SecureString" --overwrite
   aws ssm put-parameter --name "/sentry/telegram_webhook_secret" --value "<RANDOM_SECRET_HEX>" --type "SecureString" --overwrite
   ```
4. register the webhook with telegram:
   ```bash
   make telegram-hook
   ```
