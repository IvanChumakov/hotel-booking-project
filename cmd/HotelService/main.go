package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"hotel-booking/internal/app/handlers"
	"hotel-booking/internal/app/services"
	pb "hotel-booking/protos"
	"log"
	"net"
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
	port, _ := os.LookupEnv("HOTEL_PORT")

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/get_hotels", http.HandlerFunc(handlers.GetHotels))
	mux.Handle("/add_hotel", http.HandlerFunc(handlers.AddHotel))

	go func() {
		log.Printf("Starting server at port %s", port)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	grpcPort, _ := os.LookupEnv("PORT")
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServer(grpcServer, services.NewServer())

	log.Println("Starting server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
