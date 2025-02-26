package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("error generating key: %v", err)
	}
	return hex.EncodeToString(key), nil
}
