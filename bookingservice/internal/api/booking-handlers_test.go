package api

import (
	"net/http"
	"testing"
)

var customerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibG9naW4iLCJyb2xlIjoiY3VzdG9tZXIifQ.ITYJDLmEatyyllSmnH9q2ryiV9z1xeT2u4IoHKfGdCM"
var hotelierToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibmV3bG9naW4iLCJyb2xlIjoiaG90ZWxpZXIifQ.baZbEAvF5bS9rWOT6JNujz2zJNsJImCTOprjuKgpazI"
var client = &http.Client{}

func TestGetBookings(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/get_bookings", nil)
	req.Header.Set("Authorization", customerToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetBookingsByName(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/get_bookings_by_name?name=TestHotel", nil)
	req.Header.Set("Authorization", customerToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetFreeRoomsByDate(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/get_free_rooms?name=TestHotel&from=2024-01-01&to=2024-01-10", nil)
	req.Header.Set("Authorization", customerToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
