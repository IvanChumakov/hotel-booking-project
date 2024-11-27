package handlers

import (
	"crypto/rand"
	"encoding/json"
	"hotel-booking/internal/app/services"
	"io"
	"math"
	"math/big"
	"net/http"
)

func MakeOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var paymentInfo services.PaymentInfo
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
	if rnd.Int64()%2 == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
