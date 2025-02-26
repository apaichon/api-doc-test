package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input-file> <hex-key>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	hexKey := os.Args[2]

	// Decode hex key back to bytes
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		fmt.Printf("Error decoding key: %v\n", err)
		os.Exit(1)
	}

	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error creating cipher: %v\n", err)
		os.Exit(1)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("Error generating nonce: %v\n", err)
		os.Exit(1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Error creating GCM: %v\n", err)
		os.Exit(1)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	encoded := base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
	fmt.Println(encoded)
}
