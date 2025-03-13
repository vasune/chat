package auth

import (
	"chat/internal/config"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// создание JWT ключа
func JWTCreate(userID uint) string {
	payload := jwt.MapClaims{
		"userid": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JwtSecretKey))
	if err != nil {
		log.Println("JWT Error:", err)
		return ""
	}
	return tokenString
}

func extractBearerToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	token := strings.Replace(tokenString, "Bearer ", "", 1)
	if token == "" {
		log.Println("Error: Empty token")
		return ""
	}
	return token
}

func JWTVerify(tokenString string) (*jwt.Token, float64, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecretKey), nil
	})
	if err != nil {
		log.Println("JWT Parse Error:", err)
		return nil, 0, err
	}

	return token, claims["userid"].(float64), nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractBearerToken(r)
		token, userID, err := JWTVerify(tokenString)
		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
