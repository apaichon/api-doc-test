package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func Decrypt(encryptedStr string, hexKey string) (string, error) {
	// Decode hex key
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("error decoding key: %v", err)
	}

	// Decode base64 encrypted data
	encrypted, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}

	// Split nonce and ciphertext
	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting: %v", err)
	}

	return string(plaintext), nil
}
