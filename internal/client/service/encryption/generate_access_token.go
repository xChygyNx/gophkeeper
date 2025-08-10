package encryption

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAccessToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
