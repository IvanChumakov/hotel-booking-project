package main

import (
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"net"
	"os"

	"github.com/IvanChumakov/hotel-booking-project/deliveryservice/internal/app"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var log = logger.New()

func init() {
	if err := godotenv.Load(); err != nil {
		log.Logger.Info("No .env file found")
	}
}

func main() {
	port := os.Getenv("DELIVERY_PORT")
	if err := tracing.NewTrace(); err != nil {
		log.Logger.Error(err.Error())
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("delivery-service%s", port))
	if err != nil {
		log.Logger.Error(fmt.Sprintf("failed to listen: %s", err.Error()))
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationDeliveryServer(grpcServer, app.NewServer())

	log.Logger.Info("Starting server on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Logger.Error(fmt.Sprintf("failed to serve: %s", err.Error()))
	}
}
