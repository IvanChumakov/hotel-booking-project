package main

import (
	m "github.com/IvanChumakov/hotel-booking-project/hotel-lib/middleware"
	"log"
	"net"
	"net/http"
	"os"

	tracer "github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	_ "github.com/IvanChumakov/hotel-booking-project/hotelservice/docs"
	"github.com/IvanChumakov/hotel-booking-project/hotelservice/internal/api"
	"github.com/IvanChumakov/hotel-booking-project/hotelservice/internal/app"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"github.com/joho/godotenv"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	port, _ := os.LookupEnv("HOTEL_PORT")
	prometheusHost, _ := os.LookupEnv("PROMETHEUS_HOST")
	prometheusPort, _ := os.LookupEnv("PROMETHEUS_PORT")

	if err := tracer.NewTrace(); err != nil {
		log.Fatalf("Error initializing tracing: %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/get_hotels", m.JWTTokenVerify(m.LoggerMiddleware(http.HandlerFunc(api.GetHotels))))
	mux.Handle("/add_hotel", m.JWTTokenVerify(m.LoggerMiddleware(http.HandlerFunc(api.AddHotel))))
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))
	mux.Handle("/register", m.LoggerMiddleware(http.HandlerFunc(api.Register)))
	mux.Handle("/login", m.LoggerMiddleware(http.HandlerFunc(api.Login)))
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Printf("Starting server at port %s", port)
		err := http.ListenAndServe(port, mux)
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

	lis, err := net.Listen("tcp", "hotel-service:50051")
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
