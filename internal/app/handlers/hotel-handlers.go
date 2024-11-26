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

func AddHotel(w http.ResponseWriter, r *http.Request) {
	log.Print("/add_hotel")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Print("method checked")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("error while reading")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(string(data))

	var row database.Hotels
	err = json.Unmarshal(data, &row)
	if err != nil {
		log.Print("error while unmarshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = services.AddHotel(row)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
