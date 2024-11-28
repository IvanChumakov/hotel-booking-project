package handlers

import (
	"crypto/rand"
	"encoding/json"
	"github.com/IvanChumakov/hotel-booking-project/internal/app/services"
	"github.com/IvanChumakov/hotel-booking-project/internal/models"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"time"
)

func MakeOperation(w http.ResponseWriter, r *http.Request) {
	log.Print("/payment_operation")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var paymentInfo models.PaymentInfo
	err = json.Unmarshal(data, &paymentInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rnd, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := services.NewCallBack(http.Client{
		Timeout: time.Second * 5,
	})

	var statusCode int
	if rnd.Int64()%2 == 0 {
		statusCode, err = client.SendCallback(paymentInfo)
	} else {
		log.Print("payment failure (преднамеренный)")
		http.Error(w, "Payment failure (преднамеренный)", http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if statusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(statusCode)
}
