package services

import (
	"context"
	"hotel-booking/internal/database"
	pb "hotel-booking/protos"
)

type Server struct {
	pb.UnimplementedBookingServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetHotelData(_ context.Context, hotelData *pb.HotelData) (*pb.RoomsDataArray, error) {
	rooms, err := GetHotelRoomsByName(hotelData.HotelName)
	if err != nil {
		return nil, err
	}

	roomsDataArr := pb.RoomsDataArray{
		RoomData: make([]*pb.RoomData, 0),
	}
	for _, room := range rooms {
		roomPb := pb.RoomData{
			Price:      int32(room.Price),
			RoomNumber: int32(room.RoomNumber),
		}
		roomsDataArr.RoomData = append(roomsDataArr.RoomData, &roomPb)
	}
	return &roomsDataArr, nil
}

func GetAllHotels() ([]database.Hotels, error) {
	hotels, err := database.GetAllHotels()
	if err != nil {
		return nil, err
	}

	hotelsArr := make([]database.Hotels, 0)
	for _, value := range hotels {
		hotelsArr = append(hotelsArr, *value)
	}
	return hotelsArr, nil
}

func GetHotelRoomsByName(name string) ([]database.Room, error) {
	return database.GetRoomsByName(name)
}

func AddHotel(row database.Hotels) error {
	return database.AddHotel(row)
}
