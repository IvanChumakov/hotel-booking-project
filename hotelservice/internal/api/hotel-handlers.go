package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/auth"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"

	customErros "github.com/IvanChumakov/hotel-booking-project/hotel-lib/errors"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/metrics"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotelservice/internal/app"
)

var metric = metrics.NewMetrics()
var log = logger.New()

// GetHotels godoc
// @Summary      Получить все отели.
// @Description  Получить список всех отелей с номерами
// @Tags         Hotels
// @Accept		 application/x-www-form-urlencoded
// @Security BearerAuth
// @Produce      json
// @Success      200  {array}  models.Hotels
// @Router       /get_hotels [get]
func GetHotels(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestAllHotels()
	ctx, span := tracing.StartTracerSpan(r.Context(), "get_hotels")
	defer span.End()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hotelsArr, err := app.GetAllHotels(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(hotelsArr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddHotel godoc
// @Summary      Добавить отель
// @Description  Добавить информацию об отеле в базу
// @Tags         Hotels
// @Accept		 json
// @Param 		 hotel body models.Hotels true "Добавить отель"
// @Security BearerAuth
// @Produce      json
// @Success      200
// @Router       /add_hotel [post]
func AddHotel(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestAddHotels()
	ctx, span := tracing.StartTracerSpan(r.Context(), "add_hotel")
	defer span.End()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	role := r.Header.Get("role")
	login := r.Header.Get("login")
	if role == "customer" {
		http.Error(w, "Недостаточно привилегий", http.StatusUnauthorized)
		log.Logger.Error("недостаточно привилегий")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Logger.Error("error while reading")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var row models.Hotels
	err = json.Unmarshal(data, &row)
	row.OwnerLogin = login

	if err != nil {
		log.Logger.Error("error while unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.AddHotel(row, ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
