package api

import (
	"encoding/json"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"io"
	"log"
	"net/http"

	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/app"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/metrics"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
)

var metric = metrics.NewMetrics()

func GetBookings(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartTracerSpan(r.Context(), "get-all-bookings")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	bookings, err := app.GetAllBookings(ctx)
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
	ctx, span := tracing.StartTracerSpan(r.Context(), "get-bookings-by-name")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var hotel models.Hotels
	err = json.Unmarshal(data, &hotel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bookings, err := app.GetaBookingByName(hotel.Name, ctx)
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
	ctx, span := tracing.StartTracerSpan(r.Context(), "get-free-rooms-by-date")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var booking models.Booking
	err = json.Unmarshal(data, &booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rooms, _ := app.GetHotelRoomsWithPrice(booking, ctx)
	freeRooms, err := app.FilterRooms(booking, rooms, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(freeRooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddBooking(w http.ResponseWriter, r *http.Request) {
	log.Print("/add_booking")
	metric.IncRequestAddBooking()
	ctx, span := tracing.StartTracerSpan(r.Context(), "add-booking")
	defer span.End()

	if http.MethodPost != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var booking models.Booking
	err = json.Unmarshal(data, &booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.MakePaymentOperation(booking, ctx)
	if err != nil {
		http.Error(w, "Failed to make payment operation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = app.SendNotification(booking, ctx)
	w.WriteHeader(http.StatusOK)
}

func PaymentCallBack(w http.ResponseWriter, r *http.Request) {
	ctx, err := tracing.GetParentContextFromHeader(r.Context(), r.Header.Get("x-trace-id"))
	if err != nil {
		http.Error(w, "Failed to get parent context: "+err.Error(), http.StatusInternalServerError)
	}
	ctx, span := tracing.StartTracerSpan(ctx, "payment-call-back")
	defer span.End()

	log.Print("/payment_callback")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var paymentInfo models.PaymentInfo
	err = json.Unmarshal(body, &paymentInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.AddBooking(paymentInfo.Booking, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
