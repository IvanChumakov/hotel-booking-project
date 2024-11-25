package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hotel-booking/internal/app/handlers"
	"hotel-booking/internal/database"
	pb "hotel-booking/protos"
	"log"
	"net/http"
	"os"
	"time"
)

func GetAllBookings() ([]database.Booking, error) {
	return database.GetAllBookings()
}

func GetaBookingByName(name string) ([]database.Booking, error) {
	return database.GetBookingsByHotelName(name)
}

func GetHotelRoomsWithPrice(booking database.Booking) ([]database.Room, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewBookingClient(conn)
	response, err := client.GetHotelData(context.Background(), &pb.HotelData{
		HotelName: booking.HotelName,
	})

	if err != nil {
		log.Fatalf("grpc response error: %v", err)
		return nil, err
	}
	rooms := make([]database.Room, 0)
	for _, room := range response.RoomData {
		rooms = append(rooms, database.Room{
			Id:         uuid.Nil,
			Price:      int(room.Price),
			RoomNumber: int(room.RoomNumber),
			HotelId:    uuid.Nil,
		})
	}
	return rooms, nil
}

func FilterRooms(booking database.Booking, allRooms []database.Room) ([]database.Room, error) {
	booked, err := GetaBookingByName(booking.HotelName)
	if err != nil {
		return allRooms, err
	}
	freeRooms := make([]database.Room, 0)
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

func AddBooking(booking database.Booking) error {
	return database.AddBooking(booking)
}

func MakePaymentOperation(booking database.Booking) error {
	rooms, err := GetHotelRoomsWithPrice(booking)
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
	paymentInfo := handlers.PaymentInfo{
		Price: finalPrice,
	}
	jsonData, err := json.Marshal(paymentInfo)
	if err != nil {
		return err
	}
	port, _ := os.LookupEnv("PAYMENT_PORT")

	client := http.Client{
		Timeout: 100 * time.Second,
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost%s/payment", port),
		bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	httpData, err := client.Do(request)
	if err != nil {
		return err
	}
	if httpData.StatusCode != http.StatusOK {
		return fmt.Errorf("payment operation failed with status code %d", httpData.StatusCode)
	}
	return nil
}
