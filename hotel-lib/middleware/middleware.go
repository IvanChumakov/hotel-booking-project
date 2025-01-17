package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/redis"
	"github.com/golang-jwt/jwt/v5"
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

func CachedQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.New()

		redisClient, err := redis.NewClient()
		if err != nil {
			log.Logger.Error(fmt.Sprintf("ошибка подключения к redis: %s", err.Error()))
			http.Error(w, "Redis connection error", http.StatusBadGateway)
			return
		}

		ok, data := redisClient.GetData(r.URL.Query().Get("name"))
		if ok {
			w.Write([]byte(data))
			return
		}
		next.ServeHTTP(w, r)
	})
}
