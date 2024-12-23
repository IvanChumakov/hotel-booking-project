package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(brokers string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{client: client, topic: topic}, nil
}

func (p *Producer) SendMessage(notification models.Booking, ctx context.Context) error {
	_, span := tracing.StartTracerSpan(ctx, "send-message")
	defer span.End()

	b, _ := json.Marshal(notification)
	log.Print("marshalling before sending: ", string(b))

	traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
	headers := []kgo.RecordHeader{
		{
			Key:   "x-trace-id",
			Value: []byte(traceId),
		},
	}
	p.client.Produce(context.Background(), &kgo.Record{Topic: p.topic, Value: b, Headers: headers}, nil)
	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}
