package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {
	key := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(key); err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	// Output key in both hex and string format
	fmt.Printf("Hex key: %x\n", key)
	fmt.Printf("String key: %s\n", hex.EncodeToString(key))
}
