package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
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

func (c *CallBack) SendCallback(paymentInfo models.PaymentInfo, ctx context.Context) (int, error) {
	ctx, span := tracing.StartTracerSpan(ctx, "send-callback")
	defer span.End()

	port, _ := os.LookupEnv("BOOKING_PORT")
	jsonData, err := json.Marshal(paymentInfo)
	if err != nil {
		return 0, err
	}

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())

	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://booking-service%s/payment_callback", port), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("sending callback: ", err)
		return 0, err
	}
	request.Header.Set("x-trace-id", traceId)
	httpData, err := c.client.Do(request)
	if err != nil {
		log.Println("doing request: ", err)
		return 0, err
	}
	return httpData.StatusCode, nil
}
