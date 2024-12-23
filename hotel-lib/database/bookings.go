package database

import (
	"context"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"github.com/google/uuid"
	"log"
)

func GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	_, span := tracing.StartTracerSpan(ctx, "get-all-bookings-database")
	defer span.End()

	query := NewSqlBuilder()
	query = query.Select(make([]string, 0)).From("bookings")

	db, err := InitConnection("hotel-bookings")
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

	bookings := make([]models.Booking, 0)
	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From.Time, &booking.To.Time)

		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func GetBookingsByHotelName(hotelName string, ctx context.Context) ([]models.Booking, error) {
	_, span := tracing.StartTracerSpan(ctx, "get-bookings-by-name-database")
	defer span.End()

	query := NewSqlBuilder()
	query = query.Select(make([]string, 0)).From("bookings").Where(fmt.Sprintf("hotel_name = '%s'", hotelName))

	db, err := InitConnection("hotel-bookings")
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

	bookings := make([]models.Booking, 0)
	for rows.Next() {
		var booking models.Booking
		err = rows.Scan(&booking.Id, &booking.HotelName, &booking.RoomNumber, &booking.From.Time, &booking.To.Time)
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func AddBooking(booking models.Booking, ctx context.Context) error {
	_, span := tracing.StartTracerSpan(ctx, "add-booking-database")
	defer span.End()

	query := NewSqlBuilder()
	fromParsed := booking.From.Time.Format("2006-01-02")
	toParsed := booking.To.Time.Format("2006-01-02")

	query = query.CustomQuery(fmt.Sprintf("insert into bookings (id, hotel_name, room_number, from_date, to_date) "+
		"values ('%s', '%s', %d, '%s', '%s')",
		uuid.NewString(), booking.HotelName, booking.RoomNumber, fromParsed, toParsed)).
		Returning("id")

	db, err := InitConnection("hotel-bookings")
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
