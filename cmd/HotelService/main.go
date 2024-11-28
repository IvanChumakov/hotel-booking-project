package main

import (
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	_ "hotel-booking/cmd/HotelService/docs"
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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
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
	mux.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	go func() {
		log.Printf("Starting server at port %s", port)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//grpcPort, _ := os.LookupEnv("PORT")
	lis, err := net.Listen("tcp", "localhost:50051")
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
