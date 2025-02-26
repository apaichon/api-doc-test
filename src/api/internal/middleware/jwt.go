package middleware

import (
	"api/config"
	"api/internal/auth"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTMiddleware(excludedRoutes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println("JWT Middleware called for route:", r.URL.Path)
			// Check if the request path is in the excluded routes
			for _, route := range excludedRoutes {
				if strings.HasPrefix(r.URL.Path, route) {
					// fmt.Println("Skipping JWT validation for route:", route, r.URL.Path)
					next.ServeHTTP(w, r) // Skip JWT validation
					return
				}
			}

			// Validate the JWT token
			tokenString := getTokenFromRequest(r)
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := validateToken(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authorized, err := auth.HasUserApiPermission(r)
			// fmt.Printf("authorized:%v, err:%v", authorized, err)
			if err != nil || !authorized {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If valid, pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

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
