# Sentry

Sentry is an open-source, privacy-first anti-fraud triage tool built for older adults and families. When an unexpected text message or suspicious screenshot arrives, Sentry analyzes the content and delivers immediate, calm, and actionable guidance:

1. **What to hold off on right now**: Immediate containment advice to prevent irreversible loss.
2. **Recommended next steps**: Clear, numbered actions to verify or isolate the threat safely.
3. **A safe alternative action**: Official channels to verify the communication independently.
4. **Official national fraud hotline**: One-tap phone dialer and verified government reporting link.

When high-risk fraud is identified, linked family members and caregivers receive real-time notifications via email and Telegram, with access to an incident dashboard to acknowledge and resolve issues.

Live demo: [sentry-security-krisnakk.netlify.app](https://sentry-security-krisnakk.netlify.app)

## How It Works

Visitors can evaluate suspicious messages through three intuitive paths:

- **Screenshot Upload**: Drop or attach a screenshot of an SMS, WhatsApp message, email, or browser popup. Text is read directly inside your browser via an in-memory OCR worker, keeping your personal screenshots private.
- **Pasted Text**: Copy and paste any suspicious message, tracking link, or voicemail transcript directly into the workbench.
- **Sample Scenarios**: Test real-world examples across delivery alerts, banking suspensions, fake lottery prizes, and urgent family claims.

### Immediate Triage Output

Each check produces a tailored safety summary:
- **Alert Level**: Color-coded risk status (High Risk Alert, Caution Advised, Standard Message).
- **Category Badge**: Automatic classification (Delivery Alert, Banking Alert, Family Impersonation, Legal / Tax Threat, Utility Alert, Tech Support Alert, Online Contact, Prize Claim).
- **Detected Warning Signs**: Specific signals identified in the message (lookalike domains, urgency wording, payment redirection, QR codes).
- **Official Helplines**: Country-specific fraud assistance numbers (US FTC 1-877-382-4357, UK Action Fraud 0300 123 2040, India Cybercrime 1930, Australia Scamwatch).

### Caregiver Portal

Caregivers can register and link to family members:
- Secure authentication via Amazon Cognito.
- Real-time incident feed showing severity, timestamps, detected signals, and review status.
- One-click incident acknowledgment and resolution workflow.

## Architecture and Design Decisions

- **Deterministic Rule Engine**: Rather than relying on large language models that introduce latency, hallucinations, and recurring token costs, Sentry uses a weighted heuristic classification engine in Go. Responses are instantaneous, consistent, and zero-cost.
- **Client-Side OCR Privacy**: Screenshot text extraction executes in the browser using a bundled Tesseract worker. Image payloads are discarded immediately after analysis.
- **Zero-Maintenance Serverless**: The backend runs on AWS Lambda (ARM64 custom runtime), API Gateway HTTP API, DynamoDB on-demand, SQS async queues, and Amazon SNS notifications.
- **Accessible Frontend**: Built with Next.js and TypeScript, designed with high-contrast type meeting WCAG AAA requirements, clear tap targets, and responsive layouts across mobile, tablet, and desktop screens.

## Scam Categories Covered

- **Delivery Fee Redirection**: Fake postal hold alerts demanding redelivery charges or customs duty.
- **Banking Phishing**: False card suspension alerts, unblock links, and one-time password theft.
- **Family Impersonation**: Emergency claims from unrecognized numbers ("Mom, I broke my phone").
- **Tech Support Remote Access**: Deceptive virus warnings directing users to call call-center software traps.
- **Utility Disconnection**: Sudden electricity or water shutoff warnings demanding immediate wire or gift card payment.
- **Prize and Advance-Fee Traps**: Lottery winnings and gift cards requiring upfront processing fees.
- **Digital Arrest and Authority Threats**: Intimidation tactics impersonating law enforcement, police, or revenue agencies.
- **Online Romance Requests**: Wire transfer or gift voucher demands from online contacts.
- **Task Deposit Schemes**: Work-from-home schemes demanding upfront deposits to release earned funds.

## Local Development

### Prerequisites

- Go 1.22 or newer
- Node.js 18 or newer
- npm

### Running the Local Environment

Start the local backend devserver and mock services:

```bash
make dev
```

Or run the Go server directly:

```bash
go run ./cmd/devserver
```

In a separate terminal, launch the Next.js frontend:

```bash
cd frontend
npm install
npm run dev
```

The application will be accessible at `http://localhost:3000` with the local mock API at `http://localhost:8080`.

### Running Tests

Run the full Go test suite:

```bash
make test
```

Run the case evaluation suite across 100+ scam and benign test scenarios:

```bash
make cases
```

Run the frontend test suite:

```bash
cd frontend
npm test
```

## Deployment

Sentry deploys to AWS via Terraform and compiled Go binaries:

1. Copy the example configuration:
   ```bash
   cp inputs.example.json inputs.local.json
   ```
2. Set your AWS region, caregiver email, and optional Telegram bot token.
3. Validate preflight checks:
   ```bash
   make preflight
   ```
4. Deploy the infrastructure and services:
   ```bash
   make go-live
   ```

## Contributing and Suggestions

Contributions, feature suggestions, and new scam test patterns are welcome.

- **Add New Test Cases**: Add representative scam messages to `testdata/cases/` with expected classifications and signals.
- **Suggest Detection Rules**: Enhance detection logic in `internal/risk/` to protect against emerging threats.
- **Improve Accessibility**: Ensure all user flows remain simple, clear, and reassuring for non-technical users.

Please review `CONTRIBUTING.md` for guidelines and workflow.

## License

This project is licensed under the MIT License.
