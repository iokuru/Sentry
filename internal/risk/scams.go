package risk

import (
	"regexp"
)

var (
	damageControlRegex = regexp.MustCompile(`(?i)\b(i\s*clicked|i\s*(already\s*)?(sent|paid|wired|transferred)|i\s*gave(\s*them)?\s*(the\s*)?(otp|code|pin|digits)|entered\s*my\s*(card|debit|credit|cvv|details)|typed\s*in\s*my\s*(card|debit|credit|cvv|details)|i\s*let\s*(them|the\s*caller)\s*connect|installed\s*(the\s*)?apk|download\s*karke\s*install|maine\s*.*(paise|qr|pin|apk)|scanned\s*(the\s*)?qr|already\s*bought.*gift\s*cards?|read\s*(the\s*)?numbers\s*to|transferred\s*bitcoin|wired\s*(the\s*)?money|panicked\s*and\s*wired|spoke\s*to\s*someone.*gave.*code)\b`)

	deliveryFeeRegex = regexp.MustCompile(`(?i)\b(redelivery(\s*fee)?|postage\s*(fee|due)|customs\s*(fee|duty|clearance)|parcel\s*(is\s*)?on\s*hold|package\s*(is\s*)?on\s*hold|on\s*hold\s*due\s*to|unpaid\s*(duty|postage|fee)|pay\s*(fee|duty|postage|\$?\d+(\.\d+)?\s*(fee|postage|duty))|delivery\s*attempted.*fee|clear\s*(the\s*)?pkg|clearance\s*fee|customs\s*duty\s*of)\b`)

	bankSecurityRegex = regexp.MustCompile(`(?i)\b(card\s*(has\s*been\s*)?(temporarily\s*)?(suspended|blocked|locked|restricted)|account\s*(has\s*been\s*)?(suspended|locked|blocked|frozen)|unusual\s*(login|activity|charge|debit)|fraud\s*(charge|attempt|alert|transaction|activity)|unauthorized\s*(login|transaction|charge)|verify\s*(your\s*)?identity\s*(immediately|now|at)|unblock\s*now|reactivate\s*account|account\s*block\s*(kar|ho)|pan\s*card\s*update.*link|security\s*alert.*card)\b`)

	digitalArrestRegex = regexp.MustCompile(`(?i)\b(digital\s*arrest|narcotics\s*(wing|detected|seized)|customs\s*narcotics|remain\s*on\s*(video\s*)?call|cbi\s*(court|police|notice|warrant)|arrest\s*warrant|court\s*summons|police\s*will\s*(reach|arrest)|digital\s*arrest\s*warrant|money\s*laundering.*arrest|investigation.*police\s*will)\b`)

	upiQrRegex = regexp.MustCompile(`(?i)\b(scan(\s*this)?\s*qr\s*code\s*to\s*receive|enter(\s*your)?\s*(secret\s*)?upi\s*pin\s*to\s*(receive|credit|claim|add)|approve(\s*the)?\s*collect\s*request|phonepe.*upi\s*pin.*receive|googlepay.*upi\s*pin.*receive|qr\s*code\s*scan\s*karo.*(upi\s*)?pin|upi\s*pin\s*daalo.*account)\b`)

	kycPanRegex = regexp.MustCompile(`(?i)\b(kyc\s*(has\s*)?expired|update\s*(your\s*)?(pan|aadhaar|kyc)|pan\s*card\s*(update|details)|sim\s*(card\s*)?(will\s*be\s*)?(deactivated|blocked)|download\s*.*\.apk|install\s*.*\.apk|\.apk\b|aadhaar\s*pan\s*linking|sim\s*card\s*services\s*will\s*be\s*blocked|kyc\s*pending.*call)\b`)

	subscriptionTrapRegex = regexp.MustCompile(`(?i)\b(geek\s*squad.*(renew|charged|support)|norton.*(renew|debited|protection)|mcafee.*(renew|debited|protection)|amazon\s*prime.*(renew|charged|invoice)|auto-renew(ed)?\s*for\s*\$\d+|subscription\s*(has\s*)?auto-renewed|license\s*has\s*auto-renewed|call\s*(our\s*)?(billing\s*)?(helpdesk|support|cancellation|customer\s*desk).*to\s*(cancel|reverse|dispute)|call.*to\s*cancel\s*immediately)\b`)

	refundOverpayRegex = regexp.MustCompile(`(?i)\b(accidentally\s*(sent|credited)|sent\s*(extra|\$\d+\s*extra).*(by\s*mistake|send\s*back|refund)|overpayment\s*(error|notice)|wire\s*back.*(excess|difference)|refund\s*(the\s*)?(difference|excess)|refund\s*error.*sent|mailed\s*a.*(cashier\s*)?cheque.*(wire|difference)|transfer\s*the\s*extra|send\s*back\s*the\s*extra|issued\s*an\s*official\s*bank\s*cheque.*extra)\b`)

	prizeLotteryRegex = regexp.MustCompile(`(?i)\b(congratulations.*(won|winner)|sweepstakes\s*worth|claim\s*(your\s*)?(prize|winnings|cash\s*reward)|mega\s*jackpot|won\s*\$?\d+.*lottery|lucky\s*(grahak|draw)|pay\s*\$?\d+.*(processing|release|clearance)\s*fee|lottery\s*jeet\s*gaya|buy\s*\$?\d+\s*.*gift\s*card.*claim)\b`)

	techRemoteRegex = regexp.MustCompile(`(?i)\b(trojan(\s*spyware)?\s*detected|critical\s*system\s*breach|windows\s*defender\s*blocked|call\s*(windows\s*)?support\s*(immediately\s*)?at|download\s*anydesk|install\s*anydesk|teamviewer|remote\s*(session|removal|access)\s*(required|to\s*clean)|remote\s*screen\s*access|clean\s*device)\b`)

	utilityThreatRegex = regexp.MustCompile(`(?i)\b(notice\s*of\s*disconnection|disconnection\s*notice|(electricity|power|water|gas|service|broadband|internet)\s*(disconnection|scheduled\s*for\s*shutoff|will\s*terminate|will\s*be\s*suspended)|will\s*(be|b)\s*(shutoff|terminated|disconnected)\s*(today|tonight|in\s*\d+\s*hr)|power\s*will\s*be\s*disconnected|bijli.*connection.*kat\s*jayega|past\s*due\s*balance.*call.*immediately|power\s*shutoff|meter\s*will\s*(be|b)\s*shutoff)\b`)

	fakeJobRegex = regexp.MustCompile(`(?i)\b(remote\s*product\s*evaluators?|earn\s*\$?\d+[\s\-\$to\d]+daily.*(rating|home|reviewing)|like\s*youtube\s*videos|video\s*like\s*karke|work\s*from\s*home.*hiring\s*immediately|refundable\s*(training|equipment)\s*deposit|starter\s*kit\s*deposit|registration\s*fee.*start\s*work|ghar\s*baithe\s*kaam|hotel\s*rating\s*job)\b`)

	investmentSchemeRegex = regexp.MustCompile(`(?i)\b(cryptocurrency\s*vip\s*trading|automated\s*(arbitrage\s*)?algorithm\s*generates|daily\s*guaranteed\s*profit|10x\s*returns\s*in\s*\d+\s*days|guaranteed\s*15%\s*daily\s*return|secret\s*crypto\s*trading\s*club|zero\s*risk\s*guaranteed\s*return|send\s*0\.\d+\s*eth.*get|multiply\s*your\s*savings.*guaranteed)\b`)
)

func checkDamageControl(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if damageControlRegex.MatchString(combined) {
		return true, "user states funds were transferred, sensitive codes shared, or links opened"
	}
	return false, ""
}

func checkDeliveryFeeOrRedirection(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if deliveryFeeRegex.MatchString(combined) {
		return true, "demands payment of redelivery fee, postage charge, or customs duty"
	}
	return false, ""
}

func checkBankSecurityAlert(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if bankSecurityRegex.MatchString(combined) {
		return true, "claims banking card or account is suspended or compromised"
	}
	return false, ""
}

func checkDigitalArrestOrLegalThreat(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if digitalArrestRegex.MatchString(combined) {
		return true, "threatens digital arrest, narcotics seizure, or immediate law enforcement action"
	}
	return false, ""
}

func checkUPICollectOrQRReceive(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if upiQrRegex.MatchString(combined) {
		return true, "requests QR code scan or UPI PIN entry to supposedly receive funds"
	}
	return false, ""
}

func checkKYCPANDeadline(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if kycPanRegex.MatchString(combined) {
		return true, "demands immediate KYC, PAN, or Aadhaar update with threat of deactivation"
	}
	return false, ""
}

func checkSubscriptionRenewalTrap(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if subscriptionTrapRegex.MatchString(combined) {
		return true, "fictitious auto-renewal invoice providing a fraudulent cancellation helpline"
	}
	return false, ""
}

func checkRefundOverpaymentClaim(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if refundOverpayRegex.MatchString(combined) {
		return true, "claims an accidental overpayment and demands the difference be sent back"
	}
	return false, ""
}

func checkPrizeLotteryClaim(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if prizeLotteryRegex.MatchString(combined) {
		return true, "claims lottery or sweepstakes win requiring upfront fees or deposits"
	}
	return false, ""
}

func checkTechSupportRemoteAccess(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if techRemoteRegex.MatchString(combined) {
		return true, "claims device infection and prompts remote desktop software installation"
	}
	return false, ""
}

func checkUtilityDisconnectionThreat(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if utilityThreatRegex.MatchString(combined) {
		return true, "threatens immediate electricity, power, or utility disconnection"
	}
	return false, ""
}

func checkFakeJobOrTask(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if fakeJobRegex.MatchString(combined) {
		return true, "promises high work-from-home earnings while demanding starter deposits or fees"
	}
	return false, ""
}

func checkInvestmentTradingScheme(in Input) (bool, string) {
	combined := in.Text + " " + in.Context
	if investmentSchemeRegex.MatchString(combined) {
		return true, "promises guaranteed crypto or investment trading returns"
	}
	return false, ""
}
