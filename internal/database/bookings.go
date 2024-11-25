package database

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Booking struct {
	Id         uuid.UUID `json:"id"`
	HotelName  string    `json:"hotel_name"`
	RoomNumber int       `json:"room_number"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
}

func GetAllBookings() ([]Booking, error) {
	query := NewSqlBuilder()
	query = query.Select(make([]string, 0)).From("bookings")

	db, err := InitConnection("bookings")
	if err != nil {
		log.Fatal("database connection error")
		return nil, err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		log.Fatal("database error")
		return nil, err
	}

	bookings := make([]Booking, 0)
	for rows.Next() {
		var booking Booking
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From, &booking.To)

		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func GetBookingsByHotelName(hotelName string) ([]Booking, error) {
	query := NewSqlBuilder()
	query = query.Select(make([]string, 0)).From("bookings").Where(fmt.Sprintf("hotel_name = '%s'", hotelName))

	db, err := InitConnection("bookings")
	if err != nil {
		log.Fatal("database connection error")
		return nil, err
	}
	defer db.Close()

	rows, err := db.GetAll(query)
	if err != nil {
		log.Fatal("database error")
		return nil, err
	}

	bookings := make([]Booking, 0)
	for rows.Next() {
		var booking Booking
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From, &booking.To)
		bookings = append(bookings, booking)
	}
	return bookings, nil
}
