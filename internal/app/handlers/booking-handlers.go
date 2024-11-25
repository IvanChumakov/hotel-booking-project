package handlers

import (
	"encoding/json"
	"hotel-booking/internal/app/services"
	"hotel-booking/internal/database"
	"io"
	"log"
	"net/http"
)

func GetBookings(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_bookings")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	bookings, err := services.GetAllBookings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func GetBookingsByName(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_bookings_by_name")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var hotel database.Hotels
	err = json.Unmarshal(data, &hotel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bookings, err := services.GetaBookingByName(hotel.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func GetFreeRoomsByDate(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_free_rooms_by_date")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var booking database.Booking
	err = json.Unmarshal(data, &booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rooms, err := services.GetFreeRooms(booking)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
