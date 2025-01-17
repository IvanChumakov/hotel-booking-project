package main

import (
	"fmt"
	m "github.com/IvanChumakov/hotel-booking-project/hotel-lib/middleware"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"os"

	_ "github.com/IvanChumakov/hotel-booking-project/bookingservice/docs"
	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/api"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.New()

func init() {
	if err := godotenv.Load(); err != nil {
		log.Logger.Info("No .env file")
	}
}

// @title Swagger Booking Service API
// @version 1.0
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	port, _ := os.LookupEnv("METRICS_BOOKING")
	serverPort, _ := os.LookupEnv("BOOKING_PORT")
	prometheusHost, _ := os.LookupEnv("PROMETHEUS_HOST")

	if err := tracing.NewTrace(); err != nil {
		log.Logger.Error("Failed to create tracing", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/get_bookings", m.JWTTokenVerify(m.LoggerMiddleware(http.HandlerFunc(api.GetBookings))))
	mux.Handle("/get_bookings_by_name", m.JWTTokenVerify(m.LoggerMiddleware(m.CachedQuery(http.HandlerFunc(api.GetBookingsByName)))))
	mux.Handle("/get_free_rooms", m.JWTTokenVerify(m.LoggerMiddleware(http.HandlerFunc(api.GetFreeRoomsByDate))))
	mux.Handle("/add_booking", m.JWTTokenVerify(m.LoggerMiddleware(http.HandlerFunc(api.AddBooking))))
	mux.Handle("/payment_callback", m.LoggerMiddleware(http.HandlerFunc(api.PaymentCallBack)))
	mux.Handle("/register", m.LoggerMiddleware(http.HandlerFunc(api.Register)))
	mux.Handle("/login", m.LoggerMiddleware(http.HandlerFunc(api.Login)))
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Logger.Info(fmt.Sprintf("Starting listening metrics on %s", port))
		if err := http.ListenAndServe(prometheusHost+":"+port, nil); err != nil {
			log.Logger.Error(err.Error())
		}
	}()

	log.Logger.Info(fmt.Sprintf("Starting booking server on port %s", serverPort))
	if err := http.ListenAndServe(serverPort, mux); err != nil {
		log.Logger.Error(err.Error())
	}
}
