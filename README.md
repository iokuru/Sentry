# sentry

sentry helps older adults check suspicious screenshots or messages. users drop an image, optionally type context, and receive three clear outputs:
1. what to hold off on right now
2. a safe alternative action
3. a one tap dial button for their national fraud hotline (us ftc, uk action fraud, india 1930, australia scamwatch)

when a high risk scam is detected, caregivers receive an email notification and can review, acknowledge, or resolve the incident in the caregiver portal.

## local development

start the backend and frontend locally using stubbed aws services:

```bash
make dev
```

in another terminal, run tests:

```bash
make test
```

seed sample bank, delivery, and utility scams:

```bash
make seed
```

## deployment

deploy to aws using terraform:

```bash
make build
make plan
make apply
```

to tear down all resources:

```bash
make destroy
```

## cost notes

sentry is designed to run at zero cost on aws free tier and always-free allowances:
- s3, dynamodb (on demand pay-per-request), cloudfront, sqs, sns, cognito, api gateway http api, and eventbridge scheduler all operate within standard free tier quotas.
- no bedrock, no paid sms (alerts use sns email), no nat gateways, and no provisioned databases.
- textract and comprehend are the only services with limited free tier allowances (1,000 document pages and 50,000 text units per month).
- a strict daily request cap of 100 requests per service is enforced in code. when the cap is reached, analysis falls back to local rule-based inspection.
- a monthly aws budget alarm is set at 1 usd to catch unexpected usage.
