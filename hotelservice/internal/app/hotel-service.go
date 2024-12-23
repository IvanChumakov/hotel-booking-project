package app

import (
	"context"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/database"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
)

type Server struct {
	pb.UnimplementedBookingServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetHotelData(ctx context.Context, hotelData *pb.HotelData) (*pb.RoomsDataArray, error) {
	ctx, err := tracing.GetParentContext(ctx)
	if err != nil {
		return nil, err
	}
	ctx, span := tracing.StartTracerSpan(ctx, "get-hotel-data")
	defer span.End()

	rooms, err := GetHotelRoomsByName(hotelData.HotelName, ctx)
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

func GetAllHotels(ctx context.Context) ([]models.Hotels, error) {
	ctx, span := tracing.StartTracerSpan(ctx, "get_all_hotels_app")
	defer span.End()

	hotels, err := database.GetAllHotels(ctx)
	if err != nil {
		return nil, err
	}

	hotelsArr := make([]models.Hotels, 0)
	for _, value := range hotels {
		hotelsArr = append(hotelsArr, *value)
	}
	return hotelsArr, nil
}

func GetHotelRoomsByName(name string, ctx context.Context) ([]models.Room, error) {
	ctx, span := tracing.StartTracerSpan(ctx, "get-hotel-data-app")
	defer span.End()

	return database.GetRoomsByName(name, ctx)
}

func AddHotel(row models.Hotels, ctx context.Context) error {
	ctx, span := tracing.StartTracerSpan(ctx, "add_hotel_app")
	defer span.End()

	return database.AddHotel(row, ctx)
}
