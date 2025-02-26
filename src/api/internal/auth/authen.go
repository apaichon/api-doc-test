// Package auth provides authentication and authorization services
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"api/config"

	"github.com/dgrijalva/jwt-go"
)

// cors is a middleware function that sets the CORS headers
func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

// generateSalt generates a random salt for password hashing
func generateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// DecodeJWTToken decodes a JWT and verifies its signature with a secret key
func DecodeJWTToken(tokenString, secretKey string) (*JwtClaims, error) {
	config := config.NewConfig()
	// Parse the token and validate the signature
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is what we expect (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}

		// Return the secret key for signature verification
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and contains our expected claims
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}

// getTokenFromRequest retrieves the token from the request header or cookies
func getTokenFromRequest(r *http.Request) string {
	// Check if the token is present in the request header
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	// Check if the token is present in the request cookies
	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

// validateToken validates the JWT token
func validateToken(tokenString string) (*jwt.Token, error) {
	config := config.NewConfig()
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil // Replace "your-secret-key" with your actual secret key
	})

	// fmt.Println("Token", token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
