package risk_test

import (
	"testing"

	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/risk"
)

type scenarioFixture struct {
	name     string
	category model.Category
	country  string
	text     string
	minScore int
}

var knownFixtures = []scenarioFixture{
	{
		name:     "Chase Bank Account Suspension",
		category: model.CategoryBank,
		country:  "US",
		text:     "URGENT: Chase bank security alert. Your debit card has been locked due to unauthorized activity. Call 1-800-555-0199 or verify your PIN now at chase-verify-security.xyz to avoid permanent account closure.",
		minScore: 60,
	},
	{
		name:     "USPS Package Redelivery",
		category: model.CategoryDelivery,
		country:  "US",
		text:     "USPS: Your package could not be delivered due to incomplete address. Please update tracking details at bit.ly/usps-redelivery-99 and enter the one-time verification code to pay $1.99 redelivery fee.",
		minScore: 30,
	},
	{
		name:     "Electric Utility Disconnection",
		category: model.CategoryUtility,
		country:  "US",
		text:     "FINAL NOTICE: Electric service will be terminated today in 30 minutes due to unpaid utility bill of $428. Settle balance immediately using an Apple gift card or prepaid card.",
		minScore: 30,
	},
	{
		name:     "Lottery Sweepstakes Winner",
		category: model.CategoryPrize,
		country:  "US",
		text:     "Official Sweepstakes Winner! You have won $1,250,000 in the lucky draw lottery. Send $350 processing fee via Bitcoin or cryptocurrency to claim your prize jackpot.",
		minScore: 30,
	},
	{
		name:     "Grandson Bail Emergency",
		category: model.CategoryRelativeInNeed,
		country:  "AU",
		text:     "Hi mum, broke my phone this morning and got arrested for an accident. Text me on this new number right now, I need help with urgent bail money, please wire $1500.",
		minScore: 50,
	},
	{
		name:     "IRS Warrant and Arrest Demand",
		category: model.CategoryTax,
		country:  "US",
		text:     "IRS final warning: A federal warrant has been issued for your arrest due to unpaid back taxes. Federal revenue officers are en route to your residence. Immediate payment required.",
		minScore: 30,
	},
	{
		name:     "Online Romance Customs Fee",
		category: model.CategoryRomance,
		country:  "US",
		text:     "My darling soulmate, I want to meet you in person but I am stuck at customs on my flight. They are demanding $800 wire transfer for luggage release before letting me see you.",
		minScore: 30,
	},
	{
		name:     "Microsoft Trojan Defender Spyware",
		category: model.CategoryTechSupport,
		country:  "GB",
		text:     "Microsoft Windows Defender Alert: Trojan spyware detected on your computer. Your bank logins and passwords are compromised. Immediate action required, call support at 1-800-555-0144 to verify one-time code.",
		minScore: 50,
	},
}

func TestRulesOnlyAnalysisMatchesFixtures(t *testing.T) {
	for _, f := range knownFixtures {
		t.Run(f.name, func(t *testing.T) {
			// rules-only input: no Comprehend entities, neutral sentiment, no keyphrases
			rulesInput := risk.Input{
				Text:       f.text,
				Context:    "",
				Sentiment:  "NEUTRAL",
				KeyPhrases: nil,
				Entities:   nil,
			}

			rulesScore, rulesLevel, rulesSignals := risk.CalculateScore(rulesInput)
			rulesCat := risk.ClassifyCategory(rulesInput)
			rulesOutput := risk.BuildSeniorOutput(rulesCat, rulesLevel, f.country)

			// simulated comprehend input
			compInput := risk.Input{
				Text:       f.text,
				Context:    "",
				Sentiment:  "NEGATIVE",
				KeyPhrases: []string{"account", "alert", "security"},
				Entities:   nil,
			}

			compScore, compLevel, _ := risk.CalculateScore(compInput)
			compCat := risk.ClassifyCategory(compInput)
			compOutput := risk.BuildSeniorOutput(compCat, compLevel, f.country)

			if rulesCat != f.category {
				t.Fatalf("expected category %s, got %s", f.category, rulesCat)
			}
			if rulesCat != compCat {
				t.Fatalf("category mismatch: rules %s vs comp %s", rulesCat, compCat)
			}

			if rulesScore < f.minScore {
				t.Fatalf("score %d below expected min %d", rulesScore, f.minScore)
			}
			if rulesScore != compScore {
				t.Fatalf("score mismatch: rules %d vs comp %d", rulesScore, compScore)
			}

			if rulesLevel != compLevel {
				t.Fatalf("level mismatch: rules %s vs comp %s", rulesLevel, compLevel)
			}

			if rulesOutput.HoldOff != compOutput.HoldOff {
				t.Fatalf("hold off advice mismatch")
			}
			if rulesOutput.SafeAlternative != compOutput.SafeAlternative {
				t.Fatalf("safe alternative mismatch")
			}
			if rulesOutput.Hotline.Number != compOutput.Hotline.Number {
				t.Fatalf("hotline number mismatch")
			}

			if len(rulesSignals) == 0 {
				t.Fatalf("expected detected signals for scam scenario, got 0")
			}
		})
	}
}
