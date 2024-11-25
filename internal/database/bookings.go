package database

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Booking struct {
	Id         uuid.UUID  `json:"id"`
	HotelName  string     `json:"hotel_name"`
	RoomNumber int        `json:"room_number"`
	From       CustomDate `json:"from"`
	To         CustomDate `json:"to"`
}

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	c.Time, err = time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *CustomDate) MarshalJSON() ([]byte, error) {
	formattedDate := c.Time.Format(`"2006-01-02"`)
	return json.Marshal(formattedDate)
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
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From.Time, &booking.To.Time)

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
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From.Time, &booking.To.Time)
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func AddBooking(booking Booking) error {
	query := NewSqlBuilder()
	fromParsed := booking.From.Time.Format("2006-01-02")
	toParsed := booking.To.Time.Format("2006-01-02")

	query = query.CustomQuery(fmt.Sprintf("insert into bookings (id, hotel_name, room_number, \"from\", \"to\") "+
		"values ('%s', '%s', %d, '%s', '%s')",
		uuid.NewString(), booking.HotelName, booking.RoomNumber, fromParsed, toParsed)).
		Returning("id")

	db, err := InitConnection("bookings")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Insert(query)
	if err != nil {
		return err
	}
	return nil
}
