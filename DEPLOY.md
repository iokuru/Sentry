# sentry deployment guide

step-by-step instructions to deploy sentry to a cold aws account.

## prerequisites

ensure the following tools are installed on your machine:
- go 1.23 or newer
- node.js 18 or newer (with npm)
- terraform 1.5 or newer
- aws cli configured with an administrator or power user profile

verify credentials:
```bash
aws sts get-caller-identity
```

## step 1: configure environment variables

copy the example variables file:
```bash
cp infra/terraform.tfvars.example infra/terraform.tfvars
```

edit `infra/terraform.tfvars` with your email and region:
```hcl
aws_region         = "us-east-1"
environment        = "prod"
project_name       = "sentry"
budget_alert_email = "your-email@example.com"
```

## step 2: verify cost guard compliance

verify that all infrastructure adheres strictly to free tier and zero-cost constraints:
```bash
make costguard
```

the command confirms zero paid resources (no nat gateway, no rds, no provisioned capacity, no bedrock).

## step 3: build and deploy infrastructure

run the deployment pipeline:
```bash
make deploy
```

this automated command:
1. compiles all arm64 lambda binaries in pure go
2. builds the sveltekit static frontend
3. runs terraform plan and applies the infrastructure
4. generates the runtime `/config.json` with deployed endpoints
5. syncs the frontend assets to the private s3 bucket
6. invalidates the cloudfront edge cache

note the terraform outputs printed at the end:
- `cloudfront_url`: public website address for seniors and caregivers
- `api_url`: http api gateway base url
- `user_pool_id`: cognito user pool id
- `user_pool_client_id`: cognito app client id

## step 4: confirm sns caregiver email subscription

open your email inbox for the address specified in `budget_alert_email`.
find the email with subject:
`AWS Notification - Subscription Confirmation`

click the confirmation link inside the email to activate caregiver notifications.

## step 5: create initial caregiver account

create your caregiver account in the cognito user pool:

```bash
# retrieve user pool id from terraform
USER_POOL_ID=$(terraform -chdir=infra output -raw user_pool_id)

# create caregiver user
aws cognito-idp admin-create-user \
  --user-pool-id "$USER_POOL_ID" \
  --username "caregiver@example.com" \
  --user-attributes Name=email,Value="caregiver@example.com" Name=email_verified,Value=true \
  --message-action SUPPRESS

# set permanent password
aws cognito-idp admin-set-user-password \
  --user-pool-id "$USER_POOL_ID" \
  --username "caregiver@example.com" \
  --password "YourSecurePassword123!" \
  --permanent
```

## step 6: configure telegram alerts (optional)

telegram alerts run as a secondary channel alongside email.

1. open telegram and search for `@BotFather`.
2. send `/newbot` and follow prompts to choose a name and username.
3. copy the bot token provided (format: `123456789:ABCDefGhI...`).
4. open your new bot on telegram and send `/start`.
5. message `@userinfobot` to look up your numeric telegram user id.
6. store credentials in aws ssm parameter store:
   ```bash
   aws ssm put-parameter \
     --name "/sentry/telegram_token" \
     --value "<YOUR_BOT_TOKEN>" \
     --type "SecureString" \
     --overwrite

   aws ssm put-parameter \
     --name "/sentry/telegram_chat_id" \
     --value "<YOUR_NUMERIC_CHAT_ID>" \
     --type "SecureString" \
     --overwrite

   aws ssm put-parameter \
     --name "/sentry/telegram_webhook_secret" \
     --value "sentry_secret_$(openssl rand -hex 16)" \
     --type "SecureString" \
     --overwrite
   ```
7. register the webhook endpoint with telegram:
   ```bash
   make telegram-hook
   ```

## step 7: verify end to end system

run the live smoke test against your deployed aws infrastructure:

```bash
make smoke-live
```

this executes the complete sequence against live aws services:
1. requests a presigned s3 upload url from api gateway
2. puts a sample scam screenshot to s3
3. posts the intake request
4. polls dynamodb until async textract and risk scoring complete
5. signs in to cognito and verifies the incident appears in caregiver feed
6. tests acknowledging and resolving the incident

open the `cloudfront_url` in your browser to verify the senior intake screen and caregiver portal.
