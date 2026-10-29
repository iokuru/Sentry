package pairing

import (
	"strings"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := GenerateCode()
		if err != nil {
			t.Fatalf("GenerateCode error: %v", err)
		}
		if len(code) != CodeLength {
			t.Fatalf("expected code length %d, got %d", CodeLength, len(code))
		}
		if !IsValidCode(code) {
			t.Fatalf("generated invalid code: %s", code)
		}
		for _, ch := range code {
			if strings.ContainsRune("0O1IL8B", ch) {
				t.Fatalf("code contains lookalike character %c in %s", ch, code)
			}
		}
		seen[code] = true
	}
	if len(seen) < 95 {
		t.Errorf("low randomness in generated codes: %d unique out of 100", len(seen))
	}
}

func TestIsValidCode(t *testing.T) {
	validCode, _ := GenerateCode()
	if !IsValidCode(validCode) {
		t.Errorf("expected %s to be valid", validCode)
	}
	if !IsValidCode(strings.ToLower(validCode)) {
		t.Errorf("expected lowercase %s to normalize and be valid", validCode)
	}
	if IsValidCode("ABCDE") {
		t.Errorf("expected 5 character code to be invalid")
	}
	if IsValidCode("ABCDEF7") {
		t.Errorf("expected 7 character code to be invalid")
	}
	if IsValidCode("012345") {
		t.Errorf("expected lookalike characters 0, 1 to be invalid")
	}
	if IsValidCode("ILLO8B") {
		t.Errorf("expected lookalike characters to be invalid")
	}
}
