package main

import (
	"hotel-booking/internal/app/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/get_hotels", http.HandlerFunc(handlers.GetHotels))
	mux.Handle("/add_hotel", http.HandlerFunc(handlers.AddHotel))
	mux.Handle("/get_bookings", http.HandlerFunc(handlers.GetBookings))
	mux.Handle("/get_bookings_by_name", http.HandlerFunc(handlers.GetBookingsByName))
	mux.Handle("/get_free_rooms", http.HandlerFunc(handlers.GetFreeRoomsByDate))

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
