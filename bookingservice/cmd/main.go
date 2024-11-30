package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/api"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.New()

func init() {
	if err := godotenv.Load(); err != nil {
		log.Logger.Info("No .env file")
	}
}

func main() {
	port, _ := os.LookupEnv("METRICS_BOOKING")
	serverPort, _ := os.LookupEnv("BOOKING_PORT")
	prometheusHost, _ := os.LookupEnv("PROMETHEUS_HOST")

	mux := http.NewServeMux()

	mux.Handle("/get_bookings", http.HandlerFunc(api.GetBookings))
	mux.Handle("/get_bookings_by_name", http.HandlerFunc(api.GetBookingsByName))
	mux.Handle("/get_free_rooms", http.HandlerFunc(api.GetFreeRoomsByDate))
	mux.Handle("/add_booking", http.HandlerFunc(api.AddBooking))
	mux.Handle("/payment_callback", http.HandlerFunc(api.PaymentCallBack))
	http.Handle("/metrics", promhttp.Handler())

	wrappedMux := middleware.NewMiddleware(mux)
	go func() {
		log.Logger.Info(fmt.Sprintf("Starting listening metrics on %s", port))
		if err := http.ListenAndServe(prometheusHost + ":" + port, nil); err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	
	log.Logger.Info(fmt.Sprintf("Starting booking server on port %s", serverPort))
	if err := http.ListenAndServe(serverPort, wrappedMux); err != nil {
		log.Logger.Error(err.Error())
	}
}
