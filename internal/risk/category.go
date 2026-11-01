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
	case romanceRegex.MatchString(combined):
		return model.CategoryRomance
	case prizeRegex.MatchString(combined):
		return model.CategoryPrize
	case utilityRegex.MatchString(combined):
		return model.CategoryUtility
	case deliveryRegex.MatchString(combined):
		return model.CategoryDelivery
	case bankRegex.MatchString(combined):
		return model.CategoryBank
	default:
		return model.CategoryOther
	}
}

func BuildSeniorOutput(category model.Category, level model.RiskLevel, country string, damageControl ...bool) model.SeniorOutput {
	hl := hotline.Lookup(country)

	if len(damageControl) > 0 && damageControl[0] {
		return model.SeniorOutput{
			Headline:        "Take containment steps right away.",
			WhatThisIs:      "You shared information or sent money to an unverified sender. Immediate steps can help secure your accounts.",
			WhatToDoNext:    "1. Call your bank now to freeze your card or stop the transfer.\n2. Change passwords from another device.\n3. File a report at your national fraud hotline.",
			HoldOff:         "Call your bank now to freeze your accounts and stop any further transfers.",
			SafeAlternative: "Change your passwords from a separate trusted device and report the incident.",
			Hotline:         hl,
		}
	}

	if level == model.RiskLow {
		return model.SeniorOutput{
			Headline:        "This looks like a standard update.",
			WhatThisIs:      "This message matches everyday notifications like appointments, tracking updates, or requested sign-in codes.",
			WhatToDoNext:    "1. Review the details to ensure you were expecting them.\n2. Keep your codes private and never share them.\n3. Verify sender details if anything feels unfamiliar.",
			HoldOff:         "Nothing urgent found. Take your time and do not feel pressured to act immediately.",
			SafeAlternative: "Double-check the sender details through an official website or your existing contacts.",
			Hotline:         hl,
		}
	}

	if level == model.RiskMedium {
		return model.SeniorOutput{
			Headline:        "This may be an attempt to rush you.",
			WhatThisIs:      "The sender asks for money or personal details without a verified identity. Take a moment before responding.",
			WhatToDoNext:    "1. Pause and do not click links or send funds.\n2. Call the organization or person using a known number.\n3. Check with a family member before taking action.",
			HoldOff:         "Pause before clicking any link or sending money to this sender.",
			SafeAlternative: "Contact the organization through their official phone number or website.",
			Hotline:         hl,
		}
	}

	var headline, whatThisIs, whatToDoNext, holdOff, alternative string
	switch category {
	case model.CategoryBank:
		headline = "Do not reply or click the link."
		whatThisIs = "This message mimics a bank to gain access to your money. Legitimate banks do not ask for codes or passwords by SMS."
		whatToDoNext = "1. Do not tap links or reply to the message.\n2. Call the number printed on the back of your card.\n3. Review your recent transactions online."
		holdOff = "Do not tap message links, transfer money, or share any card PINs or passwords."
		alternative = "Call the trusted phone number printed on the back of your bank card."
	case model.CategoryDelivery:
		headline = "Do not pay redelivery fees from links."
		whatThisIs = "Carriers do not hold parcels for small surprise fees sent by text. These links lead to fake payment pages."
		whatToDoNext = "1. Do not click the tracking link.\n2. Track your parcel on the official carrier site or app.\n3. Delete the message if you have no pending delivery."
		holdOff = "Do not click tracking links or pay unexpected redelivery fees."
		alternative = "Track the item directly on the carrier official website or app using your original order number."
	case model.CategoryUtility:
		headline = "Do not pay shutoff fees over the phone."
		whatThisIs = "Utility providers send written notices and do not demand gift cards or immediate wire payment to avoid shutoff."
		whatToDoNext = "1. Do not make immediate payments through the link.\n2. Check your balance on your regular monthly bill.\n3. Call customer service using the number on your bill."
		holdOff = "Do not make immediate card payments over the phone or buy gift cards to stop shutoffs."
		alternative = "Check your utility account online or call the customer service number on your last paper bill."
	case model.CategoryPrize:
		headline = "Do not send money to claim any prize."
		whatThisIs = "Legitimate prizes never require an upfront fee or tax deposit to release winnings. This is a payment trap."
		whatToDoNext = "1. Do not send money, gift cards, or crypto.\n2. Delete the message and block the sender.\n3. Remember you cannot win a contest you never entered."
		holdOff = "Do not pay any processing fee, tax deposit, or shipping fee to claim winnings."
		alternative = "Delete the message and remember legitimate lotteries never require payment to release prizes."
	case model.CategoryTechSupport:
		headline = "Do not give screen access or call the number."
		whatThisIs = "Tech companies do not send popups or texts asking you to call for computer repair. They want remote access to your device."
		whatToDoNext = "1. Close the browser window or restart your device.\n2. Do not grant remote screen access to anyone.\n3. Have a trusted local professional inspect your device."
		holdOff = "Do not grant screen access, install remote software, or give callers your computer password."
		alternative = "Shut down your computer or ask a family member or local technician to inspect it in person."
	case model.CategoryRelativeInNeed:
		headline = "Call your family member on their old number."
		whatThisIs = "Scammers pretend to be children or grandchildren with a lost phone to request urgent money. Confirm their voice first."
		whatToDoNext = "1. Do not send money or gift cards.\n2. Call their regular phone number that you already have saved.\n3. Ask another family member to confirm where they are."
		holdOff = "Do not send funds or purchase vouchers until you have spoken to them directly."
		alternative = "Call the relative on their original phone number that was already saved in your contacts."
	case model.CategoryRomance:
		headline = "Do not send funds to online acquaintances."
		whatThisIs = "People met on dating apps or social media who invent emergencies are almost always running financial scams."
		whatToDoNext = "1. Do not transfer funds or send vouchers.\n2. Share the conversation with someone you trust.\n3. Stop contact if they keep asking for financial help."
		holdOff = "Do not send wire transfers, cryptocurrency, or gift cards to someone you have only met online."
		alternative = "Talk through the situation with a trusted friend or family member before making any financial decisions."
	case model.CategoryTax:
		headline = "Hang up on threats and do not transfer money."
		whatThisIs = "Government and law enforcement agencies never demand instant payment or threaten video call arrest. Real officers do not act this way."
		whatToDoNext = "1. Disconnect any active call immediately.\n2. Do not send wire transfers, UPI, or crypto.\n3. Contact your local police or official agency directly."
		holdOff = "Do not transfer money or disclose your tax identification number under threat of arrest."
		alternative = "Contact the official national tax agency directly through their verified government website."
	default:
		headline = "Do not reply or click the link."
		whatThisIs = "The sender is using pressure tactics or suspicious links to access your money or personal accounts."
		whatToDoNext = "1. Do not click links or send any payment.\n2. Verify the request through an independent official contact.\n3. Block the sender if they continue reaching out."
		holdOff = "Do not click links, send money, or share private personal information."
		alternative = "Verify the request through an established contact method you know and trust."
	}

	return model.SeniorOutput{
		Headline:        headline,
		WhatThisIs:      whatThisIs,
		WhatToDoNext:    whatToDoNext,
		HoldOff:         holdOff,
		SafeAlternative: alternative,
		Hotline:         hl,
	}
}
