package broker

import (
	"context"
	"encoding/json"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/twmb/franz-go/pkg/kgo"
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

func (p *Producer) SendMessage(notification models.Booking) {
	ctx := context.Background()
	b, _ := json.Marshal(notification)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, nil)
}

func (p *Producer) Close() {
	p.client.Close()
}
