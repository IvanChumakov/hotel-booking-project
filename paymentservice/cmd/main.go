package main

import (
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"github.com/IvanChumakov/hotel-booking-project/paymentservice/internal/api"
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
	port, _ := os.LookupEnv("PAYMENT_PORT")
	if err := tracing.NewTrace(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/payment", http.HandlerFunc(api.MakeOperation))

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
