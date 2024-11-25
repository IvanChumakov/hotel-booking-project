package main

import (
	"github.com/joho/godotenv"
	"hotel-booking/internal/app/handlers"
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

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/payment", http.HandlerFunc(handlers.MakeOperation))

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
