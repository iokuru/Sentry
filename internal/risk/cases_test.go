package risk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iokuru/sentry/internal/model"
)

type CaseFile struct {
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

func TestExecuteCaseSuite(t *testing.T) {
	casesDir := filepath.Join("..", "..", "testdata", "cases")
	files, err := filepath.Glob(filepath.Join(casesDir, "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no test case files found in %s: %v", casesDir, err)
	}

	totalCases := 0
	scamCases := 0
	safeCases := 0
	severeThreatCases := 0

	scamsClassifiedMediumOrHigh := 0
	severeThreatsClassifiedLow := 0
	safeClassifiedMediumOrHigh := 0
	fullMatches := 0

	fmt.Println("\n| ID | Scam type | Variant | Source | Expected level | Actual level | Expected cat | Actual cat | Match? | Duration |")
	fmt.Println("|---|---|---|---|---|---|---|---|---|---|")

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read case file %s: %v", f, err)
		}

		var tc CaseFile
		if err := json.Unmarshal(data, &tc); err != nil {
			t.Fatalf("parse case file %s: %v", f, err)
		}

		totalCases++
		start := time.Now()

		in := Input{
			Text:    tc.InputText,
			Context: tc.Context,
		}

		score, level, matches := CalculateScore(in)
		actualCat := string(ClassifyCategory(in))
		duration := time.Since(start)

		levelMatch := strings.EqualFold(string(level), tc.ExpectedLevel)
		catMatch := strings.EqualFold(actualCat, tc.ExpectedCategory)
		overallMatch := levelMatch && catMatch

		if overallMatch {
			fullMatches++
		}

		isSafe := strings.EqualFold(tc.ExpectedLevel, "low")
		if isSafe {
			safeCases++
			if level == model.RiskMedium || level == model.RiskHigh {
				safeClassifiedMediumOrHigh++
			}
		} else {
			scamCases++
			if level == model.RiskMedium || level == model.RiskHigh {
				scamsClassifiedMediumOrHigh++
			}
		}

		// Check severe threats
		isSevere := strings.Contains(tc.ID, "otp") ||
			strings.Contains(tc.ID, "digital-arrest") ||
			strings.Contains(tc.ID, "bank-blocked")
		if isSevere {
			severeThreatCases++
			if level == model.RiskLow {
				severeThreatsClassifiedLow++
			}
		}

		matchStr := "PASS"
		if !overallMatch {
			matchStr = "FAIL"
		}

		// Derive scam type and variant from id
		scamType := tc.ExpectedCategory
		variant := "clean"
		if strings.Contains(tc.ID, "sloppy") {
			variant = "sloppy"
		} else if strings.Contains(tc.ID, "caps") {
			variant = "caps"
		} else if strings.Contains(tc.ID, "damage") {
			variant = "damage-control"
		} else if strings.Contains(tc.ID, "hinglish") {
			variant = "hinglish"
		}

		fmt.Printf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			tc.ID, scamType, variant, tc.Source,
			tc.ExpectedLevel, string(level),
			tc.ExpectedCategory, actualCat,
			matchStr, duration.Round(time.Microsecond),
		)

		// Verify wording lint on the output
		var isDamageControl bool
		for _, m := range matches {
			if m.ID == "damage_control_stated" {
				isDamageControl = true
				break
			}
		}
		out := BuildSeniorOutput(model.Category(actualCat), level, "US", isDamageControl)
		wordingErrs := ValidateResultWording(out)
		if len(wordingErrs) > 0 {
			t.Errorf("wording lint errors for case %s: %v", tc.ID, wordingErrs)
		}
		_ = score
	}

	scamPrecision := float64(scamsClassifiedMediumOrHigh) / float64(scamCases) * 100.0
	safeFalseAlarmRate := float64(safeClassifiedMediumOrHigh) / float64(safeCases) * 100.0

	fmt.Printf("\n--- Case Suite Summary ---\n")
	fmt.Printf("Total cases:              %d\n", totalCases)
	fmt.Printf("Scam cases:               %d\n", scamCases)
	fmt.Printf("Safe cases:               %d\n", safeCases)
	fmt.Printf("Full matches:             %d/%d (%.1f%%)\n", fullMatches, totalCases, float64(fullMatches)/float64(totalCases)*100.0)
	fmt.Printf("Scam detection rate:      %d/%d (%.1f%%, target >= 95%%)\n", scamsClassifiedMediumOrHigh, scamCases, scamPrecision)
	fmt.Printf("Severe threat misses:     %d (target 0)\n", severeThreatsClassifiedLow)
	fmt.Printf("Safe cases false alarms:  %d/%d (%.1f%%, target 0%%)\n", safeClassifiedMediumOrHigh, safeCases, safeFalseAlarmRate)

	if scamPrecision < 95.0 {
		t.Errorf("scam precision %.1f%% below target of 95%%", scamPrecision)
	}
	if severeThreatsClassifiedLow > 0 {
		t.Errorf("severe threat misses %d > 0", severeThreatsClassifiedLow)
	}
	if safeClassifiedMediumOrHigh > 0 {
		t.Errorf("safe false alarm count %d > 0", safeClassifiedMediumOrHigh)
	}
}

func TestExecuteHoldoutSuite(t *testing.T) {
	holdoutDir := filepath.Join("..", "..", "testdata", "holdout")
	files, err := filepath.Glob(filepath.Join(holdoutDir, "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no holdout case files found in %s: %v", holdoutDir, err)
	}

	fmt.Println("\n================ HOLDOUT SUITE EVALUATION ================")
	fmt.Println("| ID | Expected Level | Actual Level | Expected Cat | Actual Cat | Match? |")
	fmt.Println("|---|---|---|---|---|---|")

	matches := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read holdout file %s: %v", f, err)
		}

		var tc CaseFile
		if err := json.Unmarshal(data, &tc); err != nil {
			t.Fatalf("parse holdout file %s: %v", f, err)
		}

		in := Input{
			Text:    tc.InputText,
			Context: tc.Context,
		}

		_, level, _ := CalculateScore(in)
		actualCat := string(ClassifyCategory(in))

		levelMatch := strings.EqualFold(string(level), tc.ExpectedLevel)
		catMatch := strings.EqualFold(actualCat, tc.ExpectedCategory)
		isMatch := levelMatch && catMatch
		if isMatch {
			matches++
		}

		fmt.Printf("| %s | %s | %s | %s | %s | %v |\n",
			tc.ID, tc.ExpectedLevel, string(level), tc.ExpectedCategory, actualCat, isMatch)
	}

	fmt.Printf("Holdout Accuracy: %d/%d (%.1f%%)\n", matches, len(files), float64(matches)/float64(len(files))*100.0)
	if matches < len(files) {
		t.Logf("Holdout evaluated with %d/%d passing", matches, len(files))
	}
}
