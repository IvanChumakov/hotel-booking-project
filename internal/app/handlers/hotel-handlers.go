package handlers

import (
	"encoding/json"
	"hotel-booking/internal/app/services"
	"hotel-booking/internal/database"
	"io"
	"log"
	"net/http"
)

func GetHotels(w http.ResponseWriter, r *http.Request) {
	log.Print("/get_hotels")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hotelsArr, err := services.GetAllHotels()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(hotelsArr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func AddHotel(w http.ResponseWriter, r *http.Request) {
	log.Print("/add_hotel")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var row database.Hotels
	err = json.Unmarshal(data, &row)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = services.AddHotel(row)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
