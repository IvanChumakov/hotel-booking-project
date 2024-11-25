package database

import (
	"fmt"
	"github.com/google/uuid"
)

type Hotels struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Rooms []Room    `json:"room"`
}

type Room struct {
	Id         uuid.UUID `json:"id"`
	Price      int       `json:"price"`
	HotelId    uuid.UUID `json:"hotel_id"`
	RoomNumber int       `json:"room_number"`
}

func GetAllHotels() (map[uuid.UUID]*Hotels, error) {
	query := NewSqlBuilder()
	query = query.
		Select(make([]string, 0)).
		From("hotels").
		Join("rooms", "h.id = r.hotel_id")

	db, err := InitConnection("hotel-booking")
	if err != nil {
		return make(map[uuid.UUID]*Hotels), err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		return make(map[uuid.UUID]*Hotels), err
	}

	hotelsMap := make(map[uuid.UUID]*Hotels)

	for rows.Next() {
		var hotel Hotels
		var room Room
		err = rows.Scan(&hotel.Id, &hotel.Name, &room.Id, &room.Price, &room.HotelId, &room.RoomNumber)

		if room.HotelId == uuid.Nil {
			hotelsMap[hotel.Id] = &hotel
			hotelsMap[hotel.Id].Rooms = make([]Room, 0)
			continue
		}

		if _, ok := hotelsMap[hotel.Id]; !ok {
			hotelsMap[hotel.Id] = &hotel
			hotelsMap[hotel.Id].Rooms = []Room{room}
		} else {
			hotelsMap[hotel.Id].Rooms = append(hotelsMap[hotel.Id].Rooms, room)
		}
	}

	return hotelsMap, nil
}

func AddHotel(hotel Hotels) error {
	query := NewSqlBuilder()
	query = query.CustomQuery(fmt.Sprintf("insert into hotels (id, name) "+
		"values ('%s', '%s')", uuid.NewString(), hotel.Name)).Returning("id")

	db, err := InitConnection("hotel-booking")
	if err != nil {
		return err
	}
	defer db.Close()

	hotelId, err := db.Insert(query)
	if err != nil {
		return err
	}
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

func GetRoomsByName(name string) ([]Room, error) {
	query := NewSqlBuilder()
	query = query.
		Select([]string{"r.room_number", "r.price"}).
		From("hotels").
		Join("rooms", "h.id = r.hotel_id").
		Where(fmt.Sprintf("name = '%s'", name))

	db, err := InitConnection("hotel-booking")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		return nil, err
	}
	rooms := make([]Room, 0)
	for rows.Next() {
		var room Room
		err = rows.Scan(&room.RoomNumber, &room.Price)

		rooms = append(rooms, room)
	}
	return rooms, nil
}
