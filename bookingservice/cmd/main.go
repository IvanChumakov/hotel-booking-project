package main

import (
	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	port, _ := os.LookupEnv("BOOKING_PORT")

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/get_bookings", http.HandlerFunc(api.GetBookings))
	mux.Handle("/get_bookings_by_name", http.HandlerFunc(api.GetBookingsByName))
	mux.Handle("/get_free_rooms", http.HandlerFunc(api.GetFreeRoomsByDate))
	mux.Handle("/add_booking", http.HandlerFunc(api.AddBooking))
	mux.Handle("/payment_callback", http.HandlerFunc(api.PaymentCallBack))

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
