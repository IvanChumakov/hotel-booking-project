package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/middleware"
	"github.com/IvanChumakov/hotel-booking-project/hotelservice/internal/api"
	"github.com/IvanChumakov/hotel-booking-project/hotelservice/internal/app"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"github.com/joho/godotenv"
	_ "github.com/IvanChumakov/hotel-booking-project/hotelservice/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// @title Swagger Hotel Service API
// @version 1.0
// @host localhost:8081
// @BasePath /
func main() {
	port, _ := os.LookupEnv("HOTEL_PORT")
	prometheusHost, _ := os.LookupEnv("PROMETHEUS_HOST")
	prometheusPort, _ := os.LookupEnv("PROMETHEUS_PORT")

	mux := http.NewServeMux()

	mux.Handle("/get_hotels", http.HandlerFunc(api.GetHotels))
	mux.Handle("/add_hotel", http.HandlerFunc(api.AddHotel))
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))
	http.Handle("/metrics", promhttp.Handler())

	wrappedMux := middleware.NewMiddleware(mux)
	go func() {
		log.Printf("Starting server at port %s", port)
		err := http.ListenAndServe(port, wrappedMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Printf("Starting prometheus server at port %s", prometheusPort)
		err := http.ListenAndServe(prometheusHost+":"+prometheusPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServer(grpcServer, app.NewServer())

	log.Println("Starting server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
