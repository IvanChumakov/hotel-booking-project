package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/internal/models"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
)

type Consumer struct {
	client *kgo.Client
	topic  string
}

func NewConsumer(brokers string, topic string) (*Consumer, error) {
	groupID := uuid.New().String()
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{client: client, topic: topic}, nil
}

func (c *Consumer) ReadNotifications() {
	ctx := context.Background()
	for {
		log.Print("waiting for notifications")
		fetches := c.client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var notification models.PaymentInfo
			if err := json.Unmarshal(record.Value, &notification); err != nil {
				fmt.Printf("Error decoding notification: %v\n", err)
				continue
			}
			log.Print(notification)
		}
	}
}
