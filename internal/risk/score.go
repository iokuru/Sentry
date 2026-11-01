package risk

import (
	"github.com/iokuru/sentry/internal/model"
)

type Rule struct {
	ID      string
	Name    string
	Weight  int
	Checker func(in Input) (bool, string)
}

var Rules = []Rule{
	{
		ID:      "damage_control_stated",
		Name:    "user states funds transferred or codes shared",
		Weight:  60,
		Checker: checkDamageControl,
	},
	{
		ID:      "digital_arrest_or_legal_threat",
		Name:    "threat of digital arrest or narcotics seizure",
		Weight:  60,
		Checker: checkDigitalArrestOrLegalThreat,
	},
	{
		ID:      "upi_collect_or_qr_receive",
		Name:    "qr code or upi pin to receive money",
		Weight:  60,
		Checker: checkUPICollectOrQRReceive,
	},
	{
		ID:      "family_impersonation",
		Name:    "new number family claim",
		Weight:  60,
		Checker: checkNewNumberFamilyClaims,
	},
	{
		ID:      "tech_support_remote_access",
		Name:    "tech support remote screen access prompt",
		Weight:  60,
		Checker: checkTechSupportRemoteAccess,
	},
	{
		ID:      "subscription_renewal_trap",
		Name:    "fake subscription auto-renewal invoice",
		Weight:  60,
		Checker: checkSubscriptionRenewalTrap,
	},
	{
		ID:      "refund_overpayment_claim",
		Name:    "accidental overpayment and refund demand",
		Weight:  60,
		Checker: checkRefundOverpaymentClaim,
	},
	{
		ID:      "prize_lottery_claim",
		Name:    "lottery win requiring upfront fee",
		Weight:  60,
		Checker: checkPrizeLotteryClaim,
	},
	{
		ID:      "utility_disconnection_threat",
		Name:    "threat of immediate utility disconnection",
		Weight:  60,
		Checker: checkUtilityDisconnectionThreat,
	},
	{
		ID:      "fake_job_or_task",
		Name:    "work-from-home task deposit scheme",
		Weight:  60,
		Checker: checkFakeJobOrTask,
	},
	{
		ID:      "investment_trading_scheme",
		Name:    "guaranteed crypto or trading returns",
		Weight:  60,
		Checker: checkInvestmentTradingScheme,
	},
	{
		ID:      "kyc_pan_deadline",
		Name:    "mandatory kyc or pan deactivation deadline",
		Weight:  60,
		Checker: checkKYCPANDeadline,
	},
	{
		ID:      "delivery_fee_or_redirection",
		Name:    "delivery parcel on hold or fee link",
		Weight:  60,
		Checker: checkDeliveryFeeOrRedirection,
	},
	{
		ID:      "otp_request",
		Name:    "one-time passcode or pin request",
		Weight:  30,
		Checker: checkOTPAndPIN,
	},
	{
		ID:      "bank_security_alert",
		Name:    "fake bank card or account suspension",
		Weight:  60,
		Checker: checkBankSecurityAlert,
	},
	{
		ID:      "untraceable_payment",
		Name:    "untraceable payment request",
		Weight:  30,
		Checker: checkUntraceablePayment,
	},
	{
		ID:      "authority_impersonation",
		Name:    "authority impersonation and threats",
		Weight:  25,
		Checker: checkAuthorityImpersonation,
	},
	{
		ID:      "suspicious_link",
		Name:    "shortened or lookalike link",
		Weight:  20,
		Checker: checkLinkShortenersAndDomains,
	},
	{
		ID:      "pii_harvesting",
		Name:    "sensitive personal identity requested",
		Weight:  20,
		Checker: checkPIIRequested,
	},
	{
		ID:      "urgency_wording",
		Name:    "urgency and time pressure",
		Weight:  15,
		Checker: checkUrgencyWording,
	},
}

func CalculateScore(in Input) (int, model.RiskLevel, []model.SignalMatch) {
	total := 0
	var matches []model.SignalMatch

	for _, rule := range Rules {
		matched, desc := rule.Checker(in)
		if matched {
			total += rule.Weight
			matches = append(matches, model.SignalMatch{
				ID:          rule.ID,
				Name:        rule.Name,
				Weight:      rule.Weight,
				Description: desc,
			})
		}
	}

	score := total
	if score > 100 {
		score = 100
	}

	var level model.RiskLevel
	switch {
	case score >= 60:
		level = model.RiskHigh
	case score >= 30:
		level = model.RiskMedium
	default:
		level = model.RiskLow
	}

	return score, level, matches
}
