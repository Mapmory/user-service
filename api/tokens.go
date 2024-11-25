// api/tokens.go

package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "your_secret_key" // Replace with a secure key
)

// GenerateToken creates a new JWT token for a given userID
func GenerateToken(userID string) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * 1).Unix() // Token expires in 1 hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expiresAt,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// TokenVerifyMiddleware verifies the JWT token
func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		// Split the header to get the token part
		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		// Pass the userID to the next handler via context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
