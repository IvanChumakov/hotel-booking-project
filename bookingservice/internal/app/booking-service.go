package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tracer "github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/broker"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/database"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	ctx, span := tracer.StartTracerSpan(ctx, "get-all-bookings-app")
	defer span.End()

	return database.GetAllBookings(ctx)
}

func GetaBookingByName(name string, ctx context.Context) ([]models.Booking, error) {
	ctx, span := tracer.StartTracerSpan(ctx, "get-booking-by-name-app")
	defer span.End()

	return database.GetBookingsByHotelName(name, ctx)
}

func GetHotelRoomsWithPrice(booking models.Booking, ctx context.Context) ([]models.Room, error) {
	ctx, span := tracer.StartTracerSpan(ctx, "get-hotel-rooms-with-price")
	defer span.End()
	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

	conn, err := grpc.NewClient("hotel-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewBookingClient(conn)
	response, err := client.GetHotelData(ctx, &pb.HotelData{
		HotelName: booking.HotelName,
	})

	if err != nil {
		log.Fatalf("grpc response error: %v", err)
		return nil, err
	}
	rooms := make([]models.Room, 0)
	for _, room := range response.RoomData {
		rooms = append(rooms, models.Room{
			Id:         uuid.Nil,
			Price:      int(room.Price),
			RoomNumber: int(room.RoomNumber),
			HotelId:    uuid.Nil,
		})
	}
	return rooms, nil
}

func FilterRooms(booking models.Booking, allRooms []models.Room, ctx context.Context) ([]models.Room, error) {
	ctx, span := tracer.StartTracerSpan(ctx, "filter-rooms-app")
	defer span.End()

	booked, err := GetaBookingByName(booking.HotelName, ctx)
	if err != nil {
		return allRooms, err
	}
	freeRooms := make([]models.Room, 0)
	for _, room := range allRooms {
		flag := true
		for _, piece := range booked {
			if piece.RoomNumber == room.RoomNumber {
				if !(booking.From.Before(piece.From.Time) && booking.To.Before(piece.To.Time) ||
					booking.From.After(piece.From.Time) && booking.To.After(piece.To.Time)) {
					flag = false
					break
				}
			}
		}
		if flag {
			freeRooms = append(freeRooms, room)
		}
	}
	return freeRooms, nil
}

func AddBooking(booking models.Booking, ctx context.Context) error {
	return database.AddBooking(booking, ctx)
}

func MakePaymentOperation(booking models.Booking, ctx context.Context) error {
	ctx, span := tracer.StartTracerSpan(ctx, "make-payment-operation-app")
	defer span.End()

	rooms, err := GetHotelRoomsWithPrice(booking, ctx)
	if err != nil {
		return err
	}
	var finalPrice int
	for _, room := range rooms {
		if room.RoomNumber == booking.RoomNumber {
			finalPrice = room.Price * int(booking.To.Time.Sub(booking.From.Time).Hours()/24.0)
			break
		}
	}
	paymentInfo := models.PaymentInfo{
		Price:   finalPrice,
		Booking: booking,
	}
	jsonData, err := json.Marshal(paymentInfo)
	if err != nil {
		return err
	}
	port, _ := os.LookupEnv("PAYMENT_PORT")

	client := http.Client{
		Timeout: 100 * time.Second,
	}

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://payment-service%s/payment", port),
		bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	request.Header.Set("x-trace-id", traceId)
	httpData, err := client.Do(request)
	if err != nil {
		return err
	}
	if httpData.StatusCode != http.StatusOK {
		body, err := io.ReadAll(httpData.Body)
		if err != nil {
			return fmt.Errorf("payment operation failed (неожиданно)")
		}
		return fmt.Errorf("payment operation failed with error: %s", body)
	}
	return nil
}

func SendNotification(booking models.Booking, ctx context.Context) error {
	producer, err := broker.NewProducer("redpanda:9092", "new-topic")
	if err != nil {
		log.Print("sending message error: ", err.Error())
		return err
	}
	log.Print("connection initialized")
	err = producer.SendMessage(booking, ctx)
	if err != nil {
		log.Print("sending message error: ", err.Error())
	}
	return nil
}
