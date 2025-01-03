package middleware

import (
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type Claims struct {
	Role  string `json:"role"`
	Login string `json:"name"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret-key")

func JWTTokenVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		log := logger.New()
		if tokenString == "" {
			log.Logger.Error("токен пустой")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			log.Logger.Error("токен невалидный")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		r.Header.Set("role", claims.Role)
		r.Header.Set("login", claims.Login)

		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.New()
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Logger.Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.RequestURI(), time.Since(start)))
	})
}
