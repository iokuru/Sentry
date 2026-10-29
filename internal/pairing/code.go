package pairing

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// SafeAlphabet excludes lookalike characters: 0, O, 1, I, L, 8, B
const SafeAlphabet = "CDFGHJKMNPQRTVWXYZ234679"
const CodeLength = 6

func GenerateCode() (string, error) {
	alphabetLen := big.NewInt(int64(len(SafeAlphabet)))
	var b strings.Builder
	b.Grow(CodeLength)

	for i := 0; i < CodeLength; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		b.WriteByte(SafeAlphabet[idx.Int64()])
	}

	return b.String(), nil
}

func NormalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func IsValidCode(code string) bool {
	norm := NormalizeCode(code)
	if len(norm) != CodeLength {
		return false
	}
	for i := 0; i < len(norm); i++ {
		if !strings.ContainsRune(SafeAlphabet, rune(norm[i])) {
			return false
		}
	}
	return true
}
