package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/auth"
	customErros "github.com/IvanChumakov/hotel-booking-project/hotel-lib/errors"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"io"
	"net/http"
	"time"

	"github.com/IvanChumakov/hotel-booking-project/bookingservice/internal/app"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/metrics"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
)

var metric = metrics.NewMetrics()
var log = logger.New()

// GetBookings godoc
// @Summary      Получить все бронирования
// @Description  Получить список всех бронирований
// @Tags         Bookings
// @Accept		 application/x-www-form-urlencoded
// @Security BearerAuth
// @Produce      json
// @Success      200  {array}  models.BookingSwag
// @Router       /get_bookings [get]
func GetBookings(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestGetBookings()
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

// GetBookingsByName godoc
// @Summary      Получить все бронирования по названию отеля
// @Description  Получить список всех бронирований по названию отеля
// @Tags         Bookings
// @Accept		 json
// @Param name    query     string  false  "имя отеля"
// @Security BearerAuth
// @Produce      json
// @Success      200  {array}  models.BookingSwag
// @Router       /get_bookings_by_name [get]
func GetBookingsByName(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestGetBookingsByName()
	ctx, span := tracing.StartTracerSpan(r.Context(), "get-bookings-by-name")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	name := r.URL.Query().Get("name")
	hotel := models.Hotels{
		Name: name,
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

// GetFreeRoomsByDate godoc
// @Summary      Получить свободные комнаты по дате и названию отеля
// @Description  Получить список всех свободных комнат по дате и названию отеля
// @Tags         Bookings
// @Accept		 json
// @Param name    query     string  false  "имя отеля"
// @Param from    query     string  false  "дата заезда в формате 2006-01-02"
// @Param to    query     string  false  "дата отъезда в формате 2006-01-02"
// @Security BearerAuth
// @Produce      json
// @Success      200  {array}  models.Room
// @Router       /get_free_rooms [get]
func GetFreeRoomsByDate(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestGetFreeRoomsByDate()
	ctx, span := tracing.StartTracerSpan(r.Context(), "get-free-rooms-by-date")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	fromParsed, err := time.Parse("2006-01-02", from)
	if err != nil {
		http.Error(w, "неправильный формат даты from", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toParsed, err := time.Parse("2006-01-02", to)
	if err != nil {
		http.Error(w, "неправильный формат даты from", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	booking := models.Booking{
		HotelName: name,
		From:      models.CustomDate{Time: fromParsed},
		To:        models.CustomDate{Time: toParsed},
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

// AddBooking godoc
// @Summary      Добавить бронирование
// @Description  Добавить информацию о бронировании в базу
// @Tags         Bookings
// @Accept		 json
// @Param bookingInfo body models.BookingSwag true "Информация о бронировании"
// @Security BearerAuth
// @Produce      json
// @Success      200
// @Router       /add_booking [post]
func AddBooking(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestAddBooking()
	ctx, span := tracing.StartTracerSpan(r.Context(), "add-booking")
	defer span.End()

	if http.MethodPost != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	role := r.Header.Get("role")
	if role == "hotelier" {
		http.Error(w, "Недостаточно привилегий", http.StatusUnauthorized)
		log.Logger.Error("недостаточно привилегий")
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

// Register godoc
// @Summary      Регистрация
// @Description  Зарегистрироваться в сервисе
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.User true "Создать пользователя"
// @Success      200  {string} string
// @Router       /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Error("error while reading")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Logger.Error("error while unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := auth.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var loginErr *customErros.LoginExistsError

		if errors.As(err, &loginErr) {
			_, _ = fmt.Fprintf(w, err.Error())
		}
		log.Logger.Error("error while registering: ", err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

// Login godoc
// @Summary      Вход
// @Description  Войти в аккаунт
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.UserLogin true "Данные пользователя"
// @Success      200  {string} string
// @Router       /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Error("error while reading")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userLogin models.UserLogin
	err = json.Unmarshal(data, &userLogin)
	if err != nil {
		log.Logger.Error("error while unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := auth.Login(userLogin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var loginErr *customErros.AuthError

		if errors.As(err, &loginErr) {
			_, _ = fmt.Fprintf(w, err.Error())
		}
		log.Logger.Error("error while login: ", err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}
