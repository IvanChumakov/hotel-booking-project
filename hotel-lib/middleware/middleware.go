package middleware

import (
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := logger.New()
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Logger.Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.RequestURI(), time.Since(start)))
}

func NewMiddleware(handler http.Handler) *Logger {
	return &Logger{
		handler: handler,
	}
}
