package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/IvanChumakov/hotel-booking-project/notificationservice/internal/app"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
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
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
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
			log.Print(string(record.Value))

			var notification models.Booking
			if err := json.Unmarshal(record.Value, &notification); err != nil {
				log.Print("Error unmarshaling booking")
				continue
			}

			app.SendNotification(notification)
			c.client.MarkCommitRecords(record)
		}
	}
}
