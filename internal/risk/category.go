package risk

import (
	"regexp"
	"strings"

	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/model"
)

var (
	familyRegex   = regexp.MustCompile(`(?i)\b(hi\s*mum|hi\s*mom|hey\s*dad|hello\s*mum|hello\s*mom|new\s*number|broke\s*my\s*phone|dropped\s*my\s*phone|lost\s*my\s*phone|grandson|granddaughter|daughter|son|grandma|grandpa|family\s*dinner)\b`)
	taxRegex      = regexp.MustCompile(`(?i)\b(tax\s*refund|irs|ato|hmrc|unpaid\s*taxes|tax\s*audit|taxation\s*office|revenue\s*service|police|cbi|digital\s*arrest|arrest\s*warrant|narcotics|court\s*summons|lawsuit|national\s*insurance)\b`)
	techRegex     = regexp.MustCompile(`(?i)\b(microsoft|windows\s*defender|apple\s*support|geek\s*squad|norton|mcafee|antivirus|computer\s*virus|trojan|malware|anydesk|teamviewer|remote\s*access|remote\s*screen|firewall\s*breached|system\s*breach|tech\s*support|(subscription|annual)\s*renewal)\b`)
	romanceRegex  = regexp.MustCompile(`(?i)\b(soulmate|sweetheart|military\s*doctor|oil\s*rig|send\s*money\s*for\s*flight|my\s*love|darling|meet\s*you\s*in\s*person)\b`)
	prizeRegex    = regexp.MustCompile(`(?i)\b(lottery|winner|won|prize|sweepstakes|reward\s*points|jackpot|claim\s*reward|lucky\s*draw|congratulations\s*you\s*won|kbc\s*jio|cash\s*prize)\b`)
	utilityRegex  = regexp.MustCompile(`(?i)\b(electric|electricity|power|gas\s*bill|water\s*service|water\s*dept|broadband|fiber|spectrum|verizon|jio|airtel|sim\s*card|bijli|utility|shutoff|meter|disconnection)\b`)
	deliveryRegex = regexp.MustCompile(`(?i)\b(parcel|package|redelivery|courier|postage|customs\s*fee|customs\s*duty|clearance\s*fee|usps|ups|fedex|auspost|australia\s*post|royal\s*mail|dhl)\b`)
	bankRegex     = regexp.MustCompile(`(?i)\b(bank(\s*account)?|checking|savings|overdraft|wire|transfer|credit\s*card|debit\s*card|chase|wells\s*fargo|bank\s*of\s*america|bofa|citi|citibank|barclays|commbank|anz|nab|westpac|sbi|hdfc|icici|pan\s*card|aadhaar|upi|qr\s*code|phonepe|google\s*pay|googlepay|gpay|paytm|netbanking|overpayment|cashier\s*cheque|cheque|fraud\s*(charge|alert|department)|security\s*passcode|pending.*charge)\b`)
)

func ClassifyCategory(in Input) model.Category {
	combined := strings.ToLower(in.Text + " " + in.Context)

	// Job and investment scams belong to CategoryOther
	if fakeJobRegex.MatchString(combined) || investmentSchemeRegex.MatchString(combined) {
		return model.CategoryOther
	}

	// Refund and overpayment claims belong to CategoryBank
	if refundOverpayRegex.MatchString(combined) {
		return model.CategoryBank
	}

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

	var headline, whatThisIs, whatToDoNext, holdOff, alternative string
	switch category {
	case model.CategoryBank:
		if level == model.RiskMedium {
			headline = "Check this bank notice through your official app."
			whatThisIs = "This banking notice asks for fast action or verification. Verify your account status directly with your bank before responding."
			whatToDoNext = "1. Do not use phone numbers or links in this message.\n2. Sign in to your regular banking app.\n3. Call the phone number on your bank card if unsure."
			holdOff = "Do not share passwords, PIN numbers, or one-time codes with anyone who contacts you."
			alternative = "Check your account directly in your trusted banking app or call the number on your card."
		} else {
			headline = "Do not reply or click the link."
			whatThisIs = "This message mimics a bank to gain access to your money. Legitimate banks do not ask for codes or passwords by SMS."
			whatToDoNext = "1. Do not tap links or reply to the message.\n2. Call the number printed on the back of your card.\n3. Review your recent transactions online."
			holdOff = "Do not tap message links, transfer money, or share any card PINs or passwords."
			alternative = "Call the trusted phone number printed on the back of your bank card."
		}
	case model.CategoryDelivery:
		if level == model.RiskMedium {
			headline = "Verify delivery notices through official tracking."
			whatThisIs = "Delivery notifications with short links or vague parcel notices can be deceptive. Check tracking through your merchant order history."
			whatToDoNext = "1. Avoid clicking the link in this message.\n2. Look up the shipment in your original purchase receipt.\n3. Visit the carrier official website directly."
			holdOff = "Do not enter card details or personal addresses on unfamiliar tracking web pages."
			alternative = "Open the carrier official website directly and enter your tracking code from your receipt."
		} else {
			headline = "Do not pay redelivery fees from links."
			whatThisIs = "Carriers do not hold parcels for small surprise fees sent by text. These links lead to fake payment pages."
			whatToDoNext = "1. Do not click the tracking link.\n2. Track your parcel on the official carrier site or app.\n3. Delete the message if you have no pending delivery."
			holdOff = "Do not click tracking links or pay unexpected redelivery fees."
			alternative = "Track the item directly on the carrier official website or app using your original order number."
		}
	case model.CategoryUtility:
		if level == model.RiskMedium {
			headline = "Confirm utility billing notices independently."
			whatThisIs = "Utility alerts that press for quick payment should be checked against your current bill before you respond."
			whatToDoNext = "1. Do not call numbers provided inside the alert.\n2. Review your last utility statement or online account.\n3. Contact the utility provider using their known phone line."
			holdOff = "Do not send immediate wire payments or voucher codes to resolve utility notices."
			alternative = "Log in to your utility provider account directly or dial the number on your billing statement."
		} else {
			headline = "Do not pay shutoff fees over the phone."
			whatThisIs = "Utility providers send written notices and do not demand gift cards or immediate wire payment to avoid shutoff."
			whatToDoNext = "1. Do not make immediate payments through the link.\n2. Check your balance on your regular monthly bill.\n3. Call customer service using the number on your bill."
			holdOff = "Do not make immediate card payments over the phone or buy gift cards to stop shutoffs."
			alternative = "Check your utility account online or call the customer service number on your last paper bill."
		}
	case model.CategoryPrize:
		if level == model.RiskMedium {
			headline = "Approach unexpected prize claims with caution."
			whatThisIs = "Draws, gifts, and rewards that ask for personal information or small claim fees are often promotional lures or scams."
			whatToDoNext = "1. Do not give payment details or personal information.\n2. Ask yourself if you entered this promotion recently.\n3. Delete notifications that demand fees to claim rewards."
			holdOff = "Do not pay any fee or deposit to release rewards, gifts, or contest prizes."
			alternative = "Ignore the notification or check the official company website to see if an official event exists."
		} else {
			headline = "Do not send money to claim any prize."
			whatThisIs = "Legitimate prizes never require an upfront fee or tax deposit to release winnings. This is a payment trap."
			whatToDoNext = "1. Do not send money, gift cards, or crypto.\n2. Delete the message and block the sender.\n3. Remember you cannot win a contest you never entered."
			holdOff = "Do not pay any processing fee, tax deposit, or shipping fee to claim winnings."
			alternative = "Delete the message and remember legitimate lotteries never require payment to release prizes."
		}
	case model.CategoryTechSupport:
		if level == model.RiskMedium {
			headline = "Do not call unverified tech support numbers."
			whatThisIs = "Security warnings and software renewal notices can be designed to panic you into calling a bogus helpline."
			whatToDoNext = "1. Close the message or alert window without clicking links.\n2. Run your built in device security scanner.\n3. Ask a family member or trusted local technician for help."
			holdOff = "Do not install unknown support programs or provide passwords over the phone."
			alternative = "Consult an authorized repair center or speak with a knowledgeable family member."
		} else {
			headline = "Do not give screen access or call the number."
			whatThisIs = "Tech companies do not send popups or texts asking you to call for computer repair. They want remote access to your device."
			whatToDoNext = "1. Close the browser window or restart your device.\n2. Do not grant remote screen access to anyone.\n3. Have a trusted local professional inspect your device."
			holdOff = "Do not grant screen access, install remote software, or give callers your computer password."
			alternative = "Shut down your computer or ask a family member or local technician to inspect it in person."
		}
	case model.CategoryRelativeInNeed:
		if level == model.RiskMedium {
			headline = "Confirm the identity of family contacts."
			whatThisIs = "Messages claiming to come from friends or family from unknown numbers or accounts require personal confirmation."
			whatToDoNext = "1. Do not send cash or make transfers based on text alone.\n2. Reach out through their known phone number or social account.\n3. Ask a personal question only they would know."
			holdOff = "Do not send money or share account details until you verify their voice in person or on the phone."
			alternative = "Ring the person on their known telephone number or ask another family member."
		} else {
			headline = "Call your family member on their old number."
			whatThisIs = "Scammers pretend to be children or grandchildren with a lost phone to request urgent money. Confirm their voice first."
			whatToDoNext = "1. Do not send money or gift cards.\n2. Call their regular phone number that you already have saved.\n3. Ask another family member to confirm where they are."
			holdOff = "Do not send funds or purchase vouchers until you have spoken to them directly."
			alternative = "Call the relative on their original phone number that was already saved in your contacts."
		}
	case model.CategoryRomance:
		if level == model.RiskMedium {
			headline = "Take time before helping online contacts financially."
			whatThisIs = "Online acquaintances who share personal problems and bring up money should be treated with healthy hesitation."
			whatToDoNext = "1. Keep financial and personal details completely private.\n2. Never send emergency funds to someone you have not met in person.\n3. Discuss the conversation with a family member."
			holdOff = "Do not send electronic funds, gift vouchers, or parcels to online acquaintances."
			alternative = "Discuss the situation openly with a close relative or confidant before considering any assistance."
		} else {
			headline = "Do not send funds to online acquaintances."
			whatThisIs = "People met on dating apps or social media who invent emergencies are almost always running financial scams."
			whatToDoNext = "1. Do not transfer funds or send vouchers.\n2. Share the conversation with someone you trust.\n3. Stop contact if they keep asking for financial help."
			holdOff = "Do not send wire transfers, cryptocurrency, or gift cards to someone you have only met online."
			alternative = "Talk through the situation with a trusted friend or family member before making any financial decisions."
		}
	case model.CategoryTax:
		if level == model.RiskMedium {
			headline = "Verify official government notices independently."
			whatThisIs = "Official tax and revenue departments contact citizens through postal mail, not sudden messages requesting urgent settlement."
			whatToDoNext = "1. Do not call telephone numbers listed in the alert.\n2. Check your filing status on the official revenue portal.\n3. Contact the department using their published directory."
			holdOff = "Do not disclose identification numbers or initiate payments under pressure."
			alternative = "Look up the official tax agency phone number on their verified dot gov website and call them."
		} else {
			headline = "Hang up on threats and do not transfer money."
			whatThisIs = "Government and law enforcement agencies never demand instant payment or threaten video call arrest. Real officers do not act this way."
			whatToDoNext = "1. Disconnect any active call immediately.\n2. Do not send wire transfers, UPI, or crypto.\n3. Contact your local police or official agency directly."
			holdOff = "Do not transfer money or disclose your tax identification number under threat of arrest."
			alternative = "Contact the official national tax agency directly through their verified government website."
		}
	default:
		if level == model.RiskMedium {
			headline = "Pause and verify this unusual request."
			whatThisIs = "The message contains unexpected requests for action, payment, or details. It is safer to confirm who sent it first."
			whatToDoNext = "1. Do not click links or respond immediately.\n2. Reach out to the sender using a phone number you already know.\n3. Ask a family member for a second opinion."
			holdOff = "Pause before clicking unfamiliar links, sharing personal details, or transferring funds."
			alternative = "Reach out to the organization or person using contact information from a trusted previous bill or address book."
		} else {
			headline = "Do not reply or click the link."
			whatThisIs = "The sender is using pressure tactics or suspicious links to access your money or personal accounts."
			whatToDoNext = "1. Do not click links or send any payment.\n2. Verify the request through an independent official contact.\n3. Block the sender if they continue reaching out."
			holdOff = "Do not click links, send money, or share private personal information."
			alternative = "Verify the request through an established contact method you know and trust."
		}
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
