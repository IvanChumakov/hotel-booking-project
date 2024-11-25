package services

import (
	"context"
	"google.golang.org/grpc"
	"hotel-booking/internal/database"
	pb "hotel-booking/protos"
	"log"
)

func GetAllBookings() ([]database.Booking, error) {
	return database.GetAllBookings()
}

func GetaBookingByName(name string) ([]database.Booking, error) {
	return database.GetBookingsByHotelName(name)
}

func GetFreeRooms(booking database.Booking) ([]database.Room, error) {
	conn, err := grpc.NewClient("localhost:50051")
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
			Price:      int(room.Price),
			RoomNumber: int(room.RoomNumber),
		})
	}
	return rooms, nil
}
