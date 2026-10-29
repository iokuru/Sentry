# sentry deployment guide

step-by-step instructions to deploy sentry to a cold aws account.

## prerequisites

ensure the following tools are installed:
- go 1.23 or newer
- node.js 18 or newer (with npm)
- terraform 1.5 or newer
- aws cli configured with an administrator or power user profile

verify credentials:
```bash
aws sts get-caller-identity
```

## step 1: configure inputs

copy the example inputs file:
```bash
cp inputs.example.json inputs.local.json
```

edit `inputs.local.json` with your configuration:
```json
{
  "aws_profile": "default",
  "aws_region": "us-east-1",
  "project_name": "sentry",
  "environment": "prod",
  "caregiver_email": "caregiver@example.com",
  "budget_alert_email": "billing@example.com",
  "telegram_bot_token": "",
  "telegram_chat_id": ""
}
```

telegram fields are optional. leave them empty if you only want email alerts.

## step 2: run preflight check

verify configuration, aws connectivity, terraform tooling, and free tier compliance before deploying:
```bash
make preflight
```
on windows powershell:
```powershell
.\scripts\dev.ps1 preflight
```

this checks:
- presence and syntax of `inputs.local.json`
- aws sts caller identity using your configured profile and region
- terraform binary presence and formatting
- strict zero-cost infrastructure allowlist across all terraform files
- telegram bot token validation (if provided)

## step 3: run go-live

deploy everything in a single command:
```bash
make go-live
```
on windows powershell:
```powershell
.\scripts\dev.ps1 go-live
```

this automated command executes:
1. preflight verification
2. arm64 lambda compilation and sveltekit frontend build
3. terraform apply with zero provisioned resources
4. parameter population in aws ssm parameter store
5. telegram webhook registration (if token is provided)
6. runtime `/config.json` creation, s3 sync, and cloudfront cache invalidation
7. live smoke test against the deployed endpoints

note the outputs printed upon completion:
- `Frontend URL`: public cloudfront address for seniors and caregivers
- `Caregiver Sign In`: login path (`/caregiver`)
- `API Gateway Endpoint`: base api url
- `Cognito User Pool ID` and `Cognito App Client ID`

## step 4: confirm sns email subscription

open your email inbox for the address specified in `budget_alert_email`.
find the email from aws notifications and click the subscription confirmation link.

## step 5: create caregiver account

open the frontend url in your browser and navigate to `/caregiver`.
sign up with your email and password, or create an account via the aws cli:

```bash
USER_POOL_ID=$(terraform -chdir=infra output -raw user_pool_id)

aws cognito-idp admin-create-user \
  --user-pool-id "$USER_POOL_ID" \
  --username "caregiver@example.com" \
  --user-attributes Name=email,Value="caregiver@example.com" Name=email_verified,Value=true \
  --message-action SUPPRESS

aws cognito-idp admin-set-user-password \
  --user-pool-id "$USER_POOL_ID" \
  --username "caregiver@example.com" \
  --password "YourSecurePassword123!" \
  --permanent
```

once logged in to the portal:
1. click `pair senior` to generate a 6-character single-use invite code.
2. on the senior device, navigate to the frontend url and enter the code to link devices.

## step 6: monitor system health

check live deployment status, api availability, dlq depth, and daily quotas:
```bash
make status
```
on windows powershell:
```powershell
.\scripts\dev.ps1 status
```

## teardown

to clean up and destroy all cloud resources safely:
```bash
make teardown
```
on windows powershell:
```powershell
.\scripts\dev.ps1 teardown
```

this empties the private uploads and frontend s3 buckets, removes ssm parameters, and destroys terraform resources without leaving orphan state.
