package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func Encrypt(inputText string, hexKey string) (string, error) {
	// Decode hex key back to bytes
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("error decoding key: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %v", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %v", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(inputText), nil)
	encoded := base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
	return encoded, nil
}
