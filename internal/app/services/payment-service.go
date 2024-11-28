package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/internal/models"
	"log"
	"net/http"
	"os"
)

type CallBack struct {
	client http.Client
}

func NewCallBack(client http.Client) *CallBack {
	return &CallBack{client: client}
}

func (c *CallBack) SendCallback(paymentInfo models.PaymentInfo) (int, error) {
	port, _ := os.LookupEnv("BOOKING_PORT")
	jsonData, err := json.Marshal(paymentInfo)
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://localhost%s/payment_callback", port), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("sending callback: ", err)
		return 0, err
	}
	httpData, err := c.client.Do(request)
	if err != nil {
		log.Println("doing request: ", err)
		return 0, err
	}
	return httpData.StatusCode, nil
}
