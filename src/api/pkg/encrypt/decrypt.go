package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run decrypt.go <encrypted-string> <hex-key>")
		os.Exit(1)
	}

	encryptedStr := os.Args[1]
	hexKey := os.Args[2]

	// Decode hex key
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		fmt.Printf("Error decoding key: %v\n", err)
		os.Exit(1)
	}

	// Decode base64 encrypted data
	encrypted, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		fmt.Printf("Error decoding base64: %v\n", err)
		os.Exit(1)
	}

	// Split nonce and ciphertext
	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error creating cipher: %v\n", err)
		os.Exit(1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Error creating GCM: %v\n", err)
		os.Exit(1)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Printf("Error decrypting: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(plaintext))
}
