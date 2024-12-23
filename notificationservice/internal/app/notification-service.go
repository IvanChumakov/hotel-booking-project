package app

import (
	"context"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"google.golang.org/grpc/metadata"
	"log"
	"os"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SendNotification(booking models.Booking, ctx context.Context) error {
	ctx, span := tracing.StartTracerSpan(ctx, "send-message")
	defer span.End()

	log.Print("sending notification to delivery system")
	port, _ := os.LookupEnv("DELIVERY_PORT")

	log.Printf("delivery-service%s", port)
	conn, err := grpc.NewClient("delivery-service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("new client: %v", err)
		return err
	}
	defer conn.Close()

	client := pb.NewNotificationDeliveryClient(conn)
	notification := pb.Notification{
		RoomNumber: int32(booking.RoomNumber),
		HotelName:  booking.HotelName,
		From:       timestamppb.New(booking.From.Time),
		To:         timestamppb.New(booking.To.Time),
	}

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)
	_, err = client.SendNotification(ctx, &notification)
	if err != nil {
		log.Fatalf("sending: %v", err)
		return err
	}
	log.Print("sucessfully sent")
	return nil
}
