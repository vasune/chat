package auth

import (
	"chat/internal/config"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JwtSecretKey))
}

func extractBearerToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	token := strings.Replace(tokenString, "Bearer ", "", 1)
	if token == "" {
		return "", fmt.Errorf("no token found")
	}
	return token, nil
}

func JWTVerify(tokenString string) (*jwt.Token, float64, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecretKey), nil
	})
	if err != nil {
		return nil, 0, err
	}

	userID := claims["userid"].(float64)

	return token, userID, nil
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "Token not found", http.StatusUnauthorized)
			return
		}
		_, userID, err := JWTVerify(token)
		if err != nil {
			http.Error(w, "Wrong token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
