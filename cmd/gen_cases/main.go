package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TestCase struct {
	ID                    string   `json:"id"`
	InputText             string   `json:"input_text"`
	Context               string   `json:"context,omitempty"`
	Source                string   `json:"source"`
	ExpectedLevel         string   `json:"expected_level"`
	ExpectedCategory      string   `json:"expected_category"`
	RequiredSignals       []string `json:"required_signals"`
	RequiredWordingPoints []string `json:"required_wording_points"`
	CaregiverAlert        bool     `json:"caregiver_alert"`
}

func main() {
	casesDir := filepath.Join("testdata", "cases")
	holdoutDir := filepath.Join("testdata", "holdout")

	if err := os.MkdirAll(casesDir, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(holdoutDir, 0755); err != nil {
		panic(err)
	}

	cases := buildCases()
	for _, c := range cases {
		data, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			panic(err)
		}
		path := filepath.Join(casesDir, c.ID+".json")
		if err := os.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Generated %d cases in %s\n", len(cases), casesDir)

	holdout := buildHoldout()
	for _, c := range holdout {
		data, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			panic(err)
		}
		path := filepath.Join(holdoutDir, c.ID+".json")
		if err := os.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Generated %d holdout cases in %s\n", len(holdout), holdoutDir)
}

func buildCases() []TestCase {
	var list []TestCase

	add := func(id, text, ctx, source, level, cat string, signals, wording []string, alert bool) {
		list = append(list, TestCase{
			ID:                    id,
			InputText:             text,
			Context:               ctx,
			Source:                source,
			ExpectedLevel:         level,
			ExpectedCategory:      cat,
			RequiredSignals:       signals,
			RequiredWordingPoints: wording,
			CaregiverAlert:        alert,
		})
	}

	// 1. Delivery fee with link
	add("delivery-fee-clean", "USPS: Your package has arrived at our sorting facility but is on hold due to an unpaid redelivery fee of $1.99. Pay here: https://bit.ly/usps-redelivery", "SMS message from unknown number", "text", "high", "delivery", []string{"suspicious_link", "delivery_fee_or_redirection"}, []string{"carrier", "fee"}, true)
	add("delivery-fee-sloppy", "fedex parcel update: yr package is on hold at depo coz of 2.10 postage fee visit bit.ly/fdx-deliv now", "Received on phone", "text", "high", "delivery", []string{"suspicious_link", "delivery_fee_or_redirection"}, []string{"carrier", "fee"}, true)
	add("delivery-fee-caps", "ROYAL MAIL ALERT!! 🚨 YOUR PARCEL CANNOT BE DELIVERED WITHOUT 1.50 GBP CUSTOMS FEE. PAY IMMEDIATELY AT HTTPS://BIT.LY/RM-FEE-PAY", "Text message", "text", "high", "delivery", []string{"suspicious_link", "delivery_fee_or_redirection"}, []string{"carrier", "fee"}, true)

	// 2. Bank account blocked / suspended
	add("bank-blocked-clean", "Chase Alert: Your debit card ending in 4012 has been temporarily suspended due to suspicious activity. Verify your identity immediately at https://chase-security-verify.com", "SMS alert", "text", "high", "bank", []string{"suspicious_link", "bank_security_alert"}, []string{"bank", "card"}, true)
	add("bank-blocked-sloppy", "bofa notice: ur checking accnt is locked due to unusual login. pls verify ur details right now at bit.ly/bofa-unlock", "Text message", "text", "high", "bank", []string{"suspicious_link", "bank_security_alert"}, []string{"bank", "card"}, true)
	add("bank-blocked-caps", "WELLS FARGO WARNING: IMMEDIATE ACTION REQUIRED! YOUR ACCOUNT HAS BEEN BLOCKED DUE TO FRAUD ATTEMPT. UNBLOCK NOW: HTTPS://BIT.LY/WF-SECURE-NOW", "SMS alert", "text", "high", "bank", []string{"suspicious_link", "bank_security_alert"}, []string{"bank", "card"}, true)

	// 3. OTP or PIN request
	add("otp-request-clean", "Your account security team: A cancellation request was initiated for your wire transfer. Please share the 6-digit verification code you just received to cancel the transaction.", "Caller sent SMS during call", "text", "high", "bank", []string{"otp_request"}, []string{"verification", "code"}, true)
	add("otp-request-sloppy", "hey this is support agent mark, tell me the otp pin sent to ur phone quickly so i can stop the fraud charge", "Incoming chat", "text", "high", "bank", []string{"otp_request"}, []string{"verification", "code"}, true)
	add("otp-request-caps", "URGENT FRAUD ALERT: DO NOT IGNORE! SEND THE 6 DIGIT SECURITY PASSCODE IMMEDIATELY TO REVERSE THE PENDING $850 CHARGE!!", "SMS", "text", "high", "bank", []string{"otp_request"}, []string{"verification", "code"}, true)

	// 4. New number family claim
	add("family-claim-clean", "Hi Mom, I dropped my phone in the water and this is my temporary number. I have an urgent bill that needs to be paid before 5 PM today, could you please transfer $650 for me?", "WhatsApp message", "text", "high", "relative in need", []string{"family_impersonation", "untraceable_payment"}, []string{"family", "number"}, true)
	add("family-claim-sloppy", "hey dad broke my screen text me on this new number lost my wallet need 400 bucks for rent right now pls wire it", "Text message", "text", "high", "relative in need", []string{"family_impersonation", "untraceable_payment"}, []string{"family", "number"}, true)
	add("family-claim-caps", "HELLO MUM!! LOST MY PHONE THIS IS MY FRIEND PHONE NUMBER! EMERGENCY PLEASE SEND 500 POUNDS VIA BANK TRANSFER RIGHT NOW I AM STRANDED!!", "WhatsApp message", "text", "high", "relative in need", []string{"family_impersonation", "untraceable_payment"}, []string{"family", "number"}, true)

	// 5. Prize / lottery / gift card fee
	add("prize-win-clean", "Congratulations! You have been selected as the 2nd place winner of the Walmart Community Sweepstakes worth $5,000. To claim your prize, wire a $75 administrative processing fee to our claims agent.", "Letter and text", "text", "high", "prize", []string{"untraceable_payment", "prize_lottery_claim"}, []string{"prize", "processing"}, true)
	add("prize-win-sloppy", "u won 50k in national lottery draw!! just buy 100 apple gift card to pay clearance tax and claim your cash reward", "SMS", "text", "high", "prize", []string{"untraceable_payment", "prize_lottery_claim"}, []string{"prize", "gift"}, true)
	add("prize-win-caps", "WINNER NOTIFICATION!! 🏆 YOU WON $250,000 IN THE MEGA JACKPOT!! PAY $250 RELEASE FEE VIA BITCOIN TO CLAIM YOUR FUNDS TODAY!!", "SMS message", "text", "high", "prize", []string{"untraceable_payment", "prize_lottery_claim"}, []string{"prize", "fee"}, true)

	// 6. Tech support pop up / remote access
	add("tech-support-clean", "Microsoft Security Alert: Trojan spyware detected on your computer. Your financial data is at risk. Call Windows Support immediately at +1-800-555-0199 and do not restart your device.", "Browser pop up", "text", "high", "tech support", []string{"tech_support_remote_access", "authority_impersonation"}, []string{"screen", "remote"}, true)
	add("tech-support-sloppy", "apple support warning: virus found in system files download anydesk immediately and call 18885550144 to clean device", "Website warning", "text", "high", "tech support", []string{"tech_support_remote_access"}, []string{"screen", "remote"}, true)
	add("tech-support-caps", "CRITICAL SYSTEM BREACH DETECTED!! ⚠️ WINDOWS DEFENDER BLOCKED HACKER ATTACK. CALL SUPPORT NOW AT 1-800-555-0122 FOR REMOTE REMOVAL!!", "Screen message", "text", "high", "tech support", []string{"tech_support_remote_access", "authority_impersonation"}, []string{"screen", "remote"}, true)

	// 7. Police, tax, court threat / digital arrest
	add("police-threat-clean", "Internal Revenue Service: You have unpaid federal taxes totaling $4,820. An arrest warrant has been issued in your district. To avoid federal custody, wire settlement funds today.", "Voicemail transcript", "text", "high", "tax", []string{"authority_impersonation", "untraceable_payment"}, []string{"tax", "arrest"}, true)
	add("police-threat-sloppy", "cbi police notice: narcotics detected in yr parcel sent from mumbai. u are under digital arrest remain on video call or police will reach your home", "WhatsApp video notice", "text", "high", "tax", []string{"authority_impersonation", "digital_arrest_or_legal_threat"}, []string{"police", "arrest"}, true)
	add("police-threat-caps", "AUSTRALIAN TAXATION OFFICE: FINAL NOTICE BEFORE COURT SUMMONS! PAY OUTSTANDING DEBT OF $3400 VIA BITCOIN ATM WITHIN 2 HOURS OR POLICE WILL ARREST YOU!!", "SMS text", "text", "high", "tax", []string{"authority_impersonation", "untraceable_payment"}, []string{"tax", "arrest"}, true)

	// 8. Investment / crypto / trading group
	add("investment-crypto-clean", "Exclusive invitation: Join our elite cryptocurrency VIP trading desk. Our automated arbitrage algorithm generates 15% daily guaranteed profit. Deposit $500 in USDT to start trading.", "Telegram message", "text", "high", "other", []string{"investment_trading_scheme", "untraceable_payment"}, []string{"trading", "deposit"}, true)
	add("investment-crypto-sloppy", "whatsapp trading group: sir make 10x returns in 3 days on bitcoin invest 200 today withdraw 2000 tomorrow 100% genuine guaranteed", "WhatsApp invite", "text", "high", "other", []string{"investment_trading_scheme"}, []string{"trading", "deposit"}, true)
	add("investment-crypto-caps", "INVESTMENT OPPORTUNITY OF THE YEAR!! 🚀 JOIN OUR SECRET CRYPTO TRADING CLUB AND MULTIPLY YOUR SAVINGS! ZERO RISK GUARANTEED RETURN!!", "SMS text", "text", "high", "other", []string{"investment_trading_scheme"}, []string{"trading", "deposit"}, true)

	// 9. Romance or friend in need asking for transfer (medium or high)
	add("romance-need-clean", "My love, my deployment as a military doctor in Syria is almost complete. My emergency medical trunk is stuck at customs and I need $1,200 for clearance so I can fly home to you.", "Dating app match message", "text", "medium", "romance", []string{"untraceable_payment"}, []string{"funds", "online"}, false)
	add("romance-need-sloppy", "darling my oil rig had a breakdown and my credit card isnt working here can u send wire transfer 800 for my flight tickets to meet u", "Email message", "text", "medium", "romance", []string{"untraceable_payment"}, []string{"funds", "online"}, false)
	add("romance-need-caps", "SWEETHEART I AM DETAINED AT AIRPORT IMMIGRATION! THEY DEMAND $950 CASH VIA WESTERN UNION TO RELEASE MY PASSPORT PLEASE HELP ME MY LOVE!!", "Direct message", "text", "medium", "romance", []string{"untraceable_payment"}, []string{"funds", "online"}, false)

	// 10. Fake job offer / task scam
	add("fake-job-clean", "Global Media Partners is hiring remote product evaluators. Earn $300 to $500 per day working 1 hour from home. Submit a refundable equipment deposit of $150 to receive your tasks.", "Job inquiry response", "text", "high", "other", []string{"fake_job_or_task", "untraceable_payment"}, []string{"job", "fee"}, true)
	add("fake-job-sloppy", "part time job offer: like youtube videos and get 50 rs per like daily payout 3000 rs send 500 registration fee on upi to start work", "Telegram chat", "text", "high", "other", []string{"fake_job_or_task"}, []string{"job", "fee"}, true)
	add("fake-job-caps", "WORK FROM HOME HIRING IMMEDIATELY!! 💼 EARN $4000 WEEKLY REVIEWING HOTELS ONLINE! PAY $99 STARTER KIT DEPOSIT TO BEGIN WORK TODAY!!", "SMS text", "text", "high", "other", []string{"fake_job_or_task"}, []string{"job", "fee"}, true)

	// 11. Utility disconnection threat
	add("utility-threat-clean", "Notice of Disconnection: Your electric service account #91823 has a past due balance of $214.50. Service will be terminated today unless payment is made by calling +1-800-555-0188.", "Automated text message", "text", "high", "utility", []string{"utility_disconnection_threat", "urgency_wording"}, []string{"utility", "bill"}, true)
	add("utility-threat-sloppy", "water dept urgent: ur water meter will b shutoff in 1 hr due to unpaid bill pay immediately with gift card or call 8885550199", "SMS alert", "text", "high", "utility", []string{"utility_disconnection_threat", "untraceable_payment"}, []string{"utility", "bill"}, true)
	add("utility-threat-caps", "FINAL NOTICE POWER SHUTOFF!! ⚡ ELECTRICITY DISCONNECTION SCHEDULED FOR 6 PM TODAY! CALL IMMEDIATE PAYMENT LINE 1-800-555-0177 TO PREVENT SHUTOFF!!", "Text alert", "text", "high", "utility", []string{"utility_disconnection_threat", "urgency_wording"}, []string{"utility", "bill"}, true)

	// 12. UPI collect request / QR scan to receive
	add("upi-collect-clean", "To receive your OLX sofa sale payment of Rs 15,000, please scan this QR code on Google Pay or PhonePe and enter your secret UPI PIN to credit your bank account.", "Marketplace chat", "text", "high", "bank", []string{"upi_collect_or_qr_receive"}, []string{"QR", "PIN"}, true)
	add("upi-collect-sloppy", "sir i made payment for table on phonepe approve the collect request of 8000 and enter upi pin so money comes in your account", "WhatsApp message", "text", "high", "bank", []string{"upi_collect_or_qr_receive"}, []string{"QR", "PIN"}, true)
	add("upi-collect-caps", "SCAN QR CODE TO RECEIVE RS 25,000 CASH REFUND ON GOOGLEPAY!! ENTER UPI PIN IMMEDIATELY TO CLAIM INSTANT BANK TRANSFER!!", "SMS text", "text", "high", "bank", []string{"upi_collect_or_qr_receive"}, []string{"QR", "PIN"}, true)

	// 13. KYC / Aadhaar / PAN update deadline with APK link
	add("kyc-pan-clean", "Dear Customer, your State Bank account KYC has expired today. Your netbanking access will be blocked within 24 hours. Update your PAN card details now by installing: http://bit.ly/sbi-kyc-update", "SMS from unknown header", "text", "high", "bank", []string{"kyc_pan_deadline", "suspicious_link"}, []string{"KYC", "PAN"}, true)
	add("kyc-pan-sloppy", "ur jio sim card will be deactivated tonight kyc not verified download jio-support.apk from bit.ly/jio-kyc to prevent disconnection", "SMS message", "text", "high", "utility", []string{"kyc_pan_deadline", "suspicious_link"}, []string{"KYC", "PAN"}, true)
	add("kyc-pan-caps", "HDFC BANK ALERT!! ⚠️ AADHAAR PAN LINKING MANDATORY! YOUR ACCOUNT WILL BE FROZEN TODAY IF NOT UPDATED. DOWNLOAD APK: HTTPS://BIT.LY/HDFC-PAN-APP", "SMS alert", "text", "high", "bank", []string{"kyc_pan_deadline", "suspicious_link"}, []string{"KYC", "PAN"}, true)

	// 14. Refund or overpayment claim
	add("refund-overpayment-clean", "We accidentally sent an overpayment of $2,400 instead of $240 to your checking account. Please wire back the excess $2,160 immediately to avoid legal prosecution.", "Email with attachment", "text", "high", "bank", []string{"refund_overpayment_claim", "untraceable_payment"}, []string{"overpaid", "difference"}, true)
	add("refund-overpayment-sloppy", "sir check your bank we sent 3500 extra by mistake please refund the difference 3000 to our upi id right now", "Caller messaging", "text", "high", "bank", []string{"refund_overpayment_claim"}, []string{"overpaid", "difference"}, true)
	add("refund-overpayment-caps", "ACCIDENTAL REFUND ERROR: WE SENT $4,500 BY MISTAKE! RETURN THE EXCESS $4,000 VIA CASH APP OR BITCOIN IMMEDIATELY OR WE REPORT FRAUD TO POLICE!!", "SMS message", "text", "high", "bank", []string{"refund_overpayment_claim", "untraceable_payment"}, []string{"overpaid", "difference"}, true)

	// 15. Subscription renewal auto-charge notice
	add("subscription-renewal-clean", "Geek Squad Order Confirmation: Your annual Total Tech Support subscription has auto-renewed for $399.99 and charged to your account. If you did not authorize this charge, call +1-800-555-0144 to cancel immediately.", "Email receipt", "text", "high", "tech support", []string{"subscription_renewal_trap"}, []string{"renewal", "cancellation"}, true)
	add("subscription-renewal-sloppy", "norton antivirus renewal notice: 449.00 usd debited from bank card for 3 yr plan. to dispute or cancel call billing support at 18885550183", "Email text", "text", "high", "tech support", []string{"subscription_renewal_trap"}, []string{"renewal", "cancellation"}, true)
	add("subscription-renewal-caps", "AMAZON TECH SUPPORT ANNUAL RENEWAL INVOICE #89128: $499.00 CHARGED TODAY! NOT YOU? CALL CANCELLATION HELPLINE IMMEDIATELY AT 1-800-555-0166 TO REVERSE PAYMENT!!", "Email invoice text", "text", "high", "tech support", []string{"subscription_renewal_trap"}, []string{"renewal", "cancellation"}, true)

	// 16. Marketplace buyer overpaying
	add("marketplace-buyer-clean", "I have mailed a certified cashier cheque for $1,200 for your dining table ($400). Please wire the remaining $800 to my mover so they can pick up the furniture tomorrow.", "Marketplace message", "text", "high", "bank", []string{"refund_overpayment_claim", "untraceable_payment"}, []string{"funds", "extra"}, true)
	add("marketplace-buyer-sloppy", "hey i sent 500 extra for shipping on paypal pls send back the extra 500 with zelle before delivery guy comes", "Text message", "text", "high", "bank", []string{"refund_overpayment_claim"}, []string{"funds", "extra"}, true)
	add("marketplace-buyer-caps", "CRAIGSLIST BUYER PAYMENT RECEIPT: SENT $1500 VIA OVERSEAS WIRE FOR $500 ITEM. WIRE BACK $1000 DIFFERENCE VIA WESTERN UNION IMMEDIATELY!!", "Chat message", "text", "high", "bank", []string{"refund_overpayment_claim", "untraceable_payment"}, []string{"funds", "extra"}, true)

	// Indian Context / Hinglish Variants
	add("hinglish-electricity-clean", "Dear Consumer, your electricity power will be disconnected at 9:30 PM from the power house because your previous month bill was not updated. Please immediately contact our power officer at 9876543210.", "SMS in Hinglish", "text", "high", "utility", []string{"utility_disconnection_threat"}, []string{"utility", "bill"}, true)
	add("hinglish-electricity-sloppy", "bijli vibhag notice: aapka bijli connection aaj raat kat jayega bill update nahi hua call karo officer ko 9876501234", "WhatsApp forward", "text", "high", "utility", []string{"utility_disconnection_threat"}, []string{"utility", "bill"}, true)
	add("hinglish-digital-arrest", "Customs narcotics wing Delhi: Parcel seized containing MDMA under your Aadhaar number. CBI court digital arrest warrant active. Connect on Skype video call immediately or police team will arrest.", "Incoming message", "text", "high", "tax", []string{"authority_impersonation", "digital_arrest_or_legal_threat"}, []string{"police", "arrest"}, true)
	add("hinglish-upi-qr-fraud", "Bhaiya maine 5000 PhonePe pe bhej diya hai, ye QR code scan karo aur apna UPI PIN daalo to paisa turant aapke bank account me add ho jayega.", "Marketplace message", "text", "high", "bank", []string{"upi_collect_or_qr_receive"}, []string{"QR", "PIN"}, true)
	add("hinglish-bank-blocked", "Aapka SBI yono account block kar diya gaya hai. Turant pan card update karne ke liye link pe click kare: bit.ly/sbi-pan-link", "SMS", "text", "high", "bank", []string{"kyc_pan_deadline", "suspicious_link"}, []string{"KYC", "PAN"}, true)
	add("hinglish-lottery-kbc", "KBC Jio Lucky Grahak Award: Badhai ho aapka Jio number 25 lakh lottery jeet gaya hai. Certificate lene ke liye registration charge 1500 gpay karo.", "WhatsApp text", "text", "high", "prize", []string{"untraceable_payment", "prize_lottery_claim"}, []string{"prize", "processing"}, true)

	// SAFE CASES (Expected Low, CaregiverAlert false)
	add("safe-bank-2fa-clean", "Your Chase one-time verification code is 492011. Do not share this code with anyone. Our employees will never call or text you to ask for it.", "2FA SMS from bank shortcode", "text", "low", "bank", []string{}, []string{"review", "expected"}, false)
	add("safe-bank-2fa-sloppy", "google verification code: 821940. do not share your code with anyone.", "Sign in prompt", "text", "low", "other", []string{}, []string{"review", "expected"}, false)
	add("safe-bank-2fa-caps", "NETFLIX SIGN IN PASSCODE: 772910. NEVER SHARE THIS CODE. ENTER IT ON YOUR DEVICE SCREEN TO CONTINUE.", "Sign-in code", "text", "low", "other", []string{}, []string{"review", "expected"}, false)

	add("safe-shipping-carrier-clean", "USPS Tracking: Package #9400111899562544 was delivered at 2:15 PM at your front door. Thank you for using USPS.", "Delivery confirmation SMS", "text", "low", "delivery", []string{}, []string{"standard", "update"}, false)
	add("safe-shipping-carrier-sloppy", "fedex update: pkg 772189901 has arrived at pickup locker. pin to open locker is on receipt.", "Carrier alert", "text", "low", "delivery", []string{}, []string{"standard", "update"}, false)
	add("safe-shipping-carrier-caps", "AUSTRALIA POST: YOUR DELIVERY #AP98124 IS OUT FOR DELIVERY TODAY WITH YOUR POSTIE. NO SIGNATURE REQUIRED.", "Delivery text", "text", "low", "delivery", []string{}, []string{"standard", "update"}, false)

	add("safe-order-confirmation-clean", "Thank you for your order at Powell Books! Order #78291 has been confirmed with 2 items. Expected delivery Tuesday.", "Order confirmation email", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-order-confirmation-sloppy", "target receipt: order 1829910 ready for store pickup. items: paper towels, laundry soap.", "Order notice", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-order-confirmation-caps", "AMAZON.COM: ORDER CONFIRMED! 1 ITEM DISPATCHED TO YOUR ADDRESS. TRACK YOUR SHIPMENT IN THE AMAZON APP.", "Order alert", "text", "low", "other", []string{}, []string{"standard", "update"}, false)

	add("safe-family-chat-clean", "Hey Grandma, hope you had a good week! Are we still having family dinner this Sunday at 1:00 PM?", "Family SMS", "text", "low", "relative in need", []string{}, []string{"standard", "update"}, false)
	add("safe-family-chat-sloppy", "hi mom did you find my blue sweater at your house? love you talk later", "Chat from daughter", "text", "low", "relative in need", []string{}, []string{"standard", "update"}, false)
	add("safe-family-chat-caps", "HEY DAD! HAPPY BIRTHDAY!! CALL US LATER WHEN YOU HAVE A MINUTE TO CATCH UP!", "Birthday text from son", "text", "low", "relative in need", []string{}, []string{"standard", "update"}, false)

	add("safe-calendar-invite-clean", "Calendar Invitation: Annual Garden Club Meeting on Thursday, October 15 from 10:00 AM to 11:30 AM at the community center.", "Calendar notification", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-calendar-invite-sloppy", "reminder: book club discussion tomorrow 4pm at the library meeting room bring your copy", "Event reminder", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-calendar-invite-caps", "ROTARY CLUB BREAKFAST MEETING: WEDNESDAY 7:30 AM AT CITY HALL. RSVP ON THE CLUB BULLETIN.", "Meeting notice", "text", "low", "other", []string{}, []string{"standard", "update"}, false)

	add("safe-pharmacy-clean", "CVS Pharmacy: Your prescription for Lisinopril is ready for pickup at 145 Main Street. Pharmacy hours today: 9:00 AM to 7:00 PM.", "Automated pharmacy text", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-pharmacy-sloppy", "walgreens rx ready: order 98211 ready at drive thru window pharmacy closes at 8pm", "Pharmacy SMS", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-pharmacy-caps", "BOOTS PHARMACY UK: YOUR REPEAT PRESCRIPTION HAS BEEN DISPENSED AND IS READY FOR COLLECTION AT HIGH STREET STORE.", "Pharmacy text", "text", "low", "other", []string{}, []string{"standard", "update"}, false)

	add("safe-verification-requested-clean", "Uber: 382910 is your verification code. You requested this code to sign in to your Uber driver account.", "Sign in verification requested right now by user", "text", "low", "other", []string{}, []string{"review", "expected"}, false)
	add("safe-verification-requested-sloppy", "apple id verification: 610283. use this code to complete your password reset on your ipad.", "User triggered reset", "text", "low", "other", []string{}, []string{"review", "expected"}, false)
	add("safe-verification-requested-caps", "MICROSOFT SECURITY CODE: 918230. USE TO SIGN IN TO YOUR OUTLOOK ACCOUNT. NEVER SHARE WITH THIRD PARTIES.", "Verification text", "text", "low", "tech support", []string{}, []string{"review", "expected"}, false)

	add("safe-appointment-reminder-clean", "Dr. Miller Office: Reminder for your dental cleaning appointment tomorrow at 10:30 AM. Reply C to confirm or call 555-0122 to reschedule.", "Healthcare appointment SMS", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-appointment-reminder-sloppy", "optometry appt reminder: eye exam scheduled friday 3pm with dr patel reply 1 to confirm", "Clinic SMS", "text", "low", "other", []string{}, []string{"standard", "update"}, false)
	add("safe-appointment-reminder-caps", "VALLEY HEALTH CLINIC: REMINDER FOR YOUR LAB TEST TOMORROW AT 8:00 AM. PLEASE REMAIN FASTING AFTER MIDNIGHT.", "Clinic alert", "text", "low", "other", []string{}, []string{"standard", "update"}, false)

	// DAMAGE CONTROL CASES (User states they already fell for it: clicked link, paid money, gave PIN)
	add("damage-control-clicked-link", "USPS: Redelivery fee pending. Pay at http://bit.ly/usps-fake-hold", "I clicked the link and entered my name and home address already", "text", "high", "delivery", []string{"damage_control_stated", "suspicious_link"}, []string{"containment", "bank"}, true)
	add("damage-control-gave-otp", "Chase Fraud team: Share the 6 digit code sent to cancel the wire transfer.", "I got scared and gave the caller the 6-digit code on the phone", "text", "high", "bank", []string{"damage_control_stated", "otp_request"}, []string{"containment", "bank"}, true)
	add("damage-control-sent-gift-cards", "Geek Squad: Virus detected, buy 3 Target gift cards for $500 each to pay cleanup technician.", "I already bought the gift cards and read the numbers to the man on the phone", "text", "high", "tech support", []string{"damage_control_stated", "untraceable_payment"}, []string{"containment", "bank"}, true)
	add("damage-control-wire-transfer", "Internal Revenue Service: Pay $3,500 tax audit settlement immediately to avoid arrest.", "I panicked and wired the money from my credit union an hour ago", "text", "high", "tax", []string{"damage_control_stated", "authority_impersonation"}, []string{"containment", "bank"}, true)
	add("damage-control-upi-qr-scanned", "Scan QR code to receive Rs 10,000 for your old refrigerator sale.", "Maine QR code scan karke apna UPI PIN daal diya aur mere account se 10,000 kat gaye", "text", "high", "bank", []string{"damage_control_stated", "upi_collect_or_qr_receive"}, []string{"containment", "bank"}, true)
	add("damage-control-card-entered", "Netflix: Payment declined, update credit card at https://bit.ly/netflix-billing-fix", "I typed in my debit card number and CVV on that website before checking", "text", "high", "bank", []string{"damage_control_stated", "suspicious_link"}, []string{"containment", "bank"}, true)
	add("damage-control-remote-installed", "Windows Support: Your computer has a trojan, install AnyDesk so we can clean it.", "I let the caller connect to my computer and they opened a black screen", "text", "high", "tech support", []string{"damage_control_stated", "tech_support_remote_access"}, []string{"containment", "bank"}, true)
	add("damage-control-crypto-deposited", "VIP Trading Group: Deposit 0.05 BTC to double your money in 24 hours.", "I transferred Bitcoin to the wallet address they gave me yesterday", "text", "high", "other", []string{"damage_control_stated", "untraceable_payment"}, []string{"containment", "bank"}, true)
	add("damage-control-romance-sent", "My love, I am trapped in the Istanbul airport, please wire $1,000 via Western Union.", "I sent the Western Union transfer this morning because he sounded desperate", "text", "high", "romance", []string{"damage_control_stated", "untraceable_payment"}, []string{"containment", "bank"}, true)
	add("damage-control-pan-apk-installed", "SBI Alert: Download sbi-update.apk to complete mandatory PAN link today.", "Maine wo APK file phone me download karke install kar li thi", "text", "high", "bank", []string{"damage_control_stated", "kyc_pan_deadline"}, []string{"containment", "bank"}, true)

	// Additional realistic scam cases
	add("telecom-broadband-threat-clean", "Verizon Support: Your home broadband and fiber connection will be suspended today due to past due balance of $128. Call +1-800-555-0182 immediately to pay.", "SMS alert", "text", "high", "utility", []string{"utility_disconnection_threat"}, []string{"utility", "bill"}, true)
	add("telecom-broadband-threat-sloppy", "spectrum notice: internet service will terminate tonight unpaid bill pay right now at bit.ly/spec-pay-now", "Text message", "text", "high", "utility", []string{"utility_disconnection_threat", "suspicious_link"}, []string{"utility", "bill"}, true)
	add("fake-job-hotel-review", "Earn $200-$400 daily from home by rating luxury hotels on TripAdvisor. No experience needed, daily instant withdrawal via Telegram: https://bit.ly/hotel-rating-job", "Job post invite", "text", "high", "other", []string{"fake_job_or_task", "suspicious_link"}, []string{"job", "fee"}, true)
	add("crypto-doubling-whatsapp", "Join our automated Binance trading desk. Send 0.1 ETH and get 0.2 ETH credited back within 6 hours. Zero risk with smart contracts.", "WhatsApp group invite", "text", "high", "other", []string{"investment_trading_scheme", "untraceable_payment"}, []string{"trading", "deposit"}, true)
	add("refund-overpayment-geek", "Geek Squad billing: We mistakenly processed a refund of $899 to your checking account instead of $89. You must send back $810 immediately to our agent via Apple gift cards.", "Email invoice", "text", "high", "tech support", []string{"subscription_renewal_trap", "untraceable_payment"}, []string{"overpaid", "difference"}, true)
	add("kyc-sim-block-clean", "Jio Alert: Your SIM card services will be blocked within 12 hours due to incomplete KYC verification. Click to update Aadhaar: http://bit.ly/jio-aadhaar-update", "SMS alert", "text", "high", "utility", []string{"kyc_pan_deadline", "suspicious_link"}, []string{"KYC", "PAN"}, true)
	add("hinglish-police-customs-clean", "Delhi Customs Office: Foreign currency parcel detected with your name and passport. Digital arrest order issued. Join Skype investigation immediately or local police will seize property.", "WhatsApp official notice", "text", "high", "tax", []string{"authority_impersonation", "digital_arrest_or_legal_threat"}, []string{"police", "arrest"}, true)
	add("hinglish-job-telegram", "Ghar baithe kaam kare aur rozana 2000 se 5000 kamaye. Sirf video like karke screenshot bhejna hai. Registration fees 300 rs pay kare.", "Telegram message", "text", "high", "other", []string{"fake_job_or_task"}, []string{"job", "fee"}, true)
	add("marketplace-fake-cheque", "I have issued an official bank cheque for $950 for your couch ($350). Please transfer the extra $600 to my transport agent via Zelle today.", "Facebook marketplace message", "text", "high", "bank", []string{"refund_overpayment_claim"}, []string{"funds", "extra"}, true)
	add("delivery-customs-duty-auspost", "AusPost: We were unable to deliver your international parcel #AU918231 due to unpaid duty of $4.20. Pay duty online here: https://bit.ly/auspost-duty-pay", "SMS alert", "text", "high", "delivery", []string{"suspicious_link", "delivery_fee_or_redirection"}, []string{"carrier", "fee"}, true)

	return list
}

func buildHoldout() []TestCase {
	return []TestCase{
		{
			ID:                    "holdout-dhl-customs-sloppy",
			InputText:             "dhl express notice: shipment 982190 on hold at hub pay 3.50 clearance fee now bit.ly/dhl-clear-pkg",
			Context:               "SMS text",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "delivery",
			RequiredSignals:       []string{"suspicious_link", "delivery_fee_or_redirection"},
			RequiredWordingPoints: []string{"carrier", "fee"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-citibank-blocked-caps",
			InputText:             "CITIBANK CRITICAL SECURITY ALERT: DEBIT CARD HAS BEEN RESTRICTED! REACTIVATE ACCOUNT AT HTTPS://BIT.LY/CITI-REACTIVE-NOW",
			Context:               "Phone message",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "bank",
			RequiredSignals:       []string{"suspicious_link", "bank_security_alert"},
			RequiredWordingPoints: []string{"bank", "card"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-hmrc-threat-clean",
			InputText:             "HM Revenue and Customs: A lawsuit has been filed under your National Insurance number. To avoid immediate court summons, settle outstanding balance of 1,850 GBP by calling 0300 555 0192.",
			Context:               "Voicemail text",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "tax",
			RequiredSignals:       []string{"authority_impersonation"},
			RequiredWordingPoints: []string{"tax", "arrest"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-telecom-cut-sloppy",
			InputText:             "airtel alert: yr mobile number will be disconnected tonight kyc pending call support officer 9812304918",
			Context:               "SMS",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "utility",
			RequiredSignals:       []string{"kyc_pan_deadline"},
			RequiredWordingPoints: []string{"KYC", "PAN"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-mcafee-renewal-clean",
			InputText:             "McAfee Antivirus Protection: Your 2-year multi-device license has auto-renewed for $299.00. If you wish to cancel this charge, call our billing helpdesk at +1-888-555-0177 within 24 hours.",
			Context:               "Email receipt",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "tech support",
			RequiredSignals:       []string{"subscription_renewal_trap"},
			RequiredWordingPoints: []string{"renewal", "cancellation"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-safe-hospital-reminder",
			InputText:             "St. Jude Medical Center: Appointment reminder for cardiology consultation on Friday at 9:15 AM with Dr. Adams. Please arrive 15 minutes early.",
			Context:               "Appointment SMS",
			Source:                "text",
			ExpectedLevel:         "low",
			ExpectedCategory:      "other",
			RequiredSignals:       []string{},
			RequiredWordingPoints: []string{"standard", "update"},
			CaregiverAlert:        false,
		},
		{
			ID:                    "holdout-safe-school-chat",
			InputText:             "Hi Grandpa, soccer practice was moved to 5 PM today. Can you still pick me up afterwards?",
			Context:               "Message from grandson",
			Source:                "text",
			ExpectedLevel:         "low",
			ExpectedCategory:      "relative in need",
			RequiredSignals:       []string{},
			RequiredWordingPoints: []string{"standard", "update"},
			CaregiverAlert:        false,
		},
		{
			ID:                    "holdout-safe-bank-statement",
			InputText:             "Barclays: Your monthly account statement for October is now available to view in your online banking app.",
			Context:               "Monthly statement notice",
			Source:                "text",
			ExpectedLevel:         "low",
			ExpectedCategory:      "bank",
			RequiredSignals:       []string{},
			RequiredWordingPoints: []string{"standard", "update"},
			CaregiverAlert:        false,
		},
		{
			ID:                    "holdout-damage-control-pin-sent",
			InputText:             "Your bank verification code is 891230. Provide to representative to verify refund.",
			Context:               "I spoke to someone claiming to be fraud support and gave them my 6 digit code",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "bank",
			RequiredSignals:       []string{"damage_control_stated", "otp_request"},
			RequiredWordingPoints: []string{"containment", "bank"},
			CaregiverAlert:        true,
		},
		{
			ID:                    "holdout-digital-arrest-hindi",
			InputText:             "Mumbai Police Cyber Crime: Aapka bank account illegal money laundering me use hua hai. Turant video call pe join kare digital arrest proceed karne ke liye.",
			Context:               "Incoming notice",
			Source:                "text",
			ExpectedLevel:         "high",
			ExpectedCategory:      "tax",
			RequiredSignals:       []string{"authority_impersonation", "digital_arrest_or_legal_threat"},
			RequiredWordingPoints: []string{"police", "arrest"},
			CaregiverAlert:        true,
		},
	}
}
