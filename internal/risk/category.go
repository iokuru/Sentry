package risk

import (
	"regexp"
	"strings"

	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/model"
)

var (
	bankRegex     = regexp.MustCompile(`(?i)\b(bank|account|checking|savings|overdraft|wire|transfer|credit\s*card|debit\s*card|chase|wells\s*fargo|bank\s*of\s*america|citi|barclays|commbank|anz|nab|westpac|sbi|hdfc|icici)\b`)
	deliveryRegex = regexp.MustCompile(`(?i)\b(parcel|package|delivery|redelivery|shipment|tracking|courier|postage|customs\s*fee|usps|ups|fedex|auspost|australia\s*post|royal\s*mail|dhl)\b`)
	utilityRegex  = regexp.MustCompile(`(?i)\b(electric|electricity|power|gas\s*bill|water\s*service|broadband|utility|shutoff|meter|disconnection)\b`)
	prizeRegex    = regexp.MustCompile(`(?i)\b(lottery|winner|won|prize|sweepstakes|reward\s*points|jackpot|claim\s*reward|lucky\s*draw|congratulations\s*you\s*won)\b`)
	techRegex     = regexp.MustCompile(`(?i)\b(microsoft|windows\s*defender|apple\s*support|geek\s*squad|computer\s*virus|trojan|malware|anydesk|teamviewer|remote\s*access|firewall\s*breached)\b`)
	familyRegex   = regexp.MustCompile(`(?i)\b(hi\s*mum|hi\s*mom|hey\s*dad|new\s*number|broke\s*my\s*phone|grandson|granddaughter|bail\s*money|family\s*emergency)\b`)
	romanceRegex  = regexp.MustCompile(`(?i)\b(soulmate|sweetheart|military\s*doctor|oil\s*rig|send\s*money\s*for\s*flight|my\s*love|darling|meet\s*you\s*in\s*person)\b`)
	taxRegex      = regexp.MustCompile(`(?i)\b(tax\s*refund|irs|ato|hmrc|unpaid\s*taxes|tax\s*audit|taxation\s*office|revenue\s*service)\b`)
)

func ClassifyCategory(in Input) model.Category {
	combined := strings.ToLower(in.Text + " " + in.Context)

	switch {
	case familyRegex.MatchString(combined):
		return model.CategoryRelativeInNeed
	case taxRegex.MatchString(combined):
		return model.CategoryTax
	case techRegex.MatchString(combined):
		return model.CategoryTechSupport
	case bankRegex.MatchString(combined):
		return model.CategoryBank
	case deliveryRegex.MatchString(combined):
		return model.CategoryDelivery
	case utilityRegex.MatchString(combined):
		return model.CategoryUtility
	case prizeRegex.MatchString(combined):
		return model.CategoryPrize
	case romanceRegex.MatchString(combined):
		return model.CategoryRomance
	default:
		return model.CategoryOther
	}
}

func BuildSeniorOutput(category model.Category, level model.RiskLevel, country string) model.SeniorOutput {
	hl := hotline.Lookup(country)

	var holdOff string
	var alternative string

	switch category {
	case model.CategoryBank:
		holdOff = "Do not tap message links, transfer money, or share any card PINs or passwords."
		alternative = "Call the trusted phone number printed on the back of your bank card."
	case model.CategoryDelivery:
		holdOff = "Do not click tracking links or pay unexpected redelivery fees."
		alternative = "Track the item directly on the carrier official website or app using your original order number."
	case model.CategoryUtility:
		holdOff = "Do not make immediate card payments over the phone or buy gift cards to stop shutoffs."
		alternative = "Check your utility account online or call the customer service number on your last paper bill."
	case model.CategoryPrize:
		holdOff = "Do not pay any processing fee, tax deposit, or shipping fee to claim winnings."
		alternative = "Delete the message and remember legitimate lotteries never require payment to release prizes."
	case model.CategoryTechSupport:
		holdOff = "Do not grant screen access, install remote software, or give callers your computer password."
		alternative = "Shut down your computer or ask a family member or local technician to inspect it in person."
	case model.CategoryRelativeInNeed:
		holdOff = "Do not send funds or purchase vouchers until you have spoken to them directly."
		alternative = "Call the relative on their original phone number that was already saved in your contacts."
	case model.CategoryRomance:
		holdOff = "Do not send wire transfers, cryptocurrency, or gift cards to someone you have only met online."
		alternative = "Talk through the situation with a trusted friend or family member before making any financial decisions."
	case model.CategoryTax:
		holdOff = "Do not transfer money or disclose your tax identification number under threat of arrest."
		alternative = "Contact the official national tax agency directly through their verified government website."
	default:
		holdOff = "Do not click links, send money, or share private personal information."
		alternative = "Verify the request through an established contact method you know and trust."
	}

	return model.SeniorOutput{
		HoldOff:         holdOff,
		SafeAlternative: alternative,
		Hotline:         hl,
	}
}
