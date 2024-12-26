package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
)

var customerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibG9naW4iLCJyb2xlIjoiY3VzdG9tZXIifQ.ITYJDLmEatyyllSmnH9q2ryiV9z1xeT2u4IoHKfGdCM"
var hotelierToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmV3bG9naW4iLCJyb2xlIjoiaG90ZWxpZXIifQ.baZbEAvF5bS9rWOT6JNujz2zJNsJImCTOprjuKgpazI"
var client = &http.Client{}

func TestGetHotels(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/get_hotels", nil)
	request.Header.Set("Authorization", customerToken)

	response, err := client.Do(request)
	if err != nil {
		log.Logger.Error("error:", err.Error(), nil)
		return
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestAddHotel(t *testing.T) {
	hotel := models.Hotels{
		Name: "Test Hotel",
		Rooms: []models.Room{
			{
				Price:      200,
				RoomNumber: 111,
			},
		},
	}
	data, _ := json.Marshal(hotel)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/add_hotel", bytes.NewBuffer(data))
	req.Header.Set("Authorization", hotelierToken)

	response, err := client.Do(req)
	if err != nil {
		log.Logger.Error("error:", err.Error(), nil)
		return
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestAddHotelUnauthorized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/add_hotel", nil)

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRegister(t *testing.T) {
	user := models.User{
		Password: "test_password",
		Login:    "test_login",
		Role:     "customer",
	}
	data, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/register", bytes.NewReader(data))

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestLogin(t *testing.T) {
	userLogin := models.UserLogin{
		Password: "test_password",
		Login:    "test_login",
	}
	data, _ := json.Marshal(userLogin)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/login", bytes.NewReader(data))

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
