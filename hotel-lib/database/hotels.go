package database

import (
	"context"
	"fmt"
	"log"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"github.com/google/uuid"
)

func GetAllHotels(ctx context.Context) (map[uuid.UUID]*models.Hotels, error) {
	_, span := tracing.StartTracerSpan(ctx, "get-all-hotels-database")
	defer span.End()

	query := NewSqlBuilder()
	query = query.
		Select(make([]string, 0)).
		From("hotels").
		Join("rooms", "h.id = r.hotel_id")
	log.Println(query.query)

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return make(map[uuid.UUID]*models.Hotels), err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		return make(map[uuid.UUID]*models.Hotels), err
	}

	hotelsMap := make(map[uuid.UUID]*models.Hotels)

	for rows.Next() {
		var hotel models.Hotels
		var room models.Room
		err = rows.Scan(&hotel.Id, &hotel.Name, &hotel.OwnerLogin, &room.Id, &room.Price, &room.HotelId, &room.RoomNumber)

		if room.HotelId == uuid.Nil {
			hotelsMap[hotel.Id] = &hotel
			hotelsMap[hotel.Id].Rooms = make([]models.Room, 0)
			continue
		}

		if _, ok := hotelsMap[hotel.Id]; !ok {
			hotelsMap[hotel.Id] = &hotel
			hotelsMap[hotel.Id].Rooms = []models.Room{room}
		} else {
			hotelsMap[hotel.Id].Rooms = append(hotelsMap[hotel.Id].Rooms, room)
		}
	}

	return hotelsMap, nil
}

func AddHotel(hotel models.Hotels, ctx context.Context) error {
	_, span := tracing.StartTracerSpan(ctx, "add-hotel-database")
	defer span.End()

	user, err := GetUser(hotel.OwnerLogin)
	if err != nil {
		return err
	}

	query := NewSqlBuilder()
	query = query.CustomQuery(fmt.Sprintf("insert into hotels (id, name, owner_id) "+
		"values ('%s', '%s', '%s')", uuid.NewString(), hotel.Name, user.Id)).Returning("id")

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return err
	}
	defer db.Close()
	log.Print(query.query)

	hotelId, err := db.Insert(query)
	if err != nil {
		return err
	}
	log.Print("inserted data")
	for _, room := range hotel.Rooms {
		query.Clear()
		query = query.CustomQuery(fmt.Sprintf("insert into rooms (id, price, hotel_id, room_number) values ('%s', %d, '%s', %d)",
			uuid.NewString(), room.Price, hotelId, room.RoomNumber)).Returning("id")
		_, err = db.Insert(query)

		if err != nil {
			return err
		}
	}
	return nil
}

func GetRoomsByName(name string, ctx context.Context) ([]models.Room, error) {
	_, span := tracing.StartTracerSpan(ctx, "get-hotel-data-database")
	defer span.End()

	query := NewSqlBuilder()
	query = query.
		Select([]string{"r.room_number", "r.price"}).
		From("hotels").
		Join("rooms", "h.id = r.hotel_id").
		Where(fmt.Sprintf("name = '%s'", name))
	log.Println(query.query)

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		return nil, err
	}
	rooms := make([]models.Room, 0)
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.RoomNumber, &room.Price)

		rooms = append(rooms, room)
	}
	return rooms, nil
}
