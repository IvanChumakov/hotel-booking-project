package main

import (
	"context"
	"github.com/IvanChumakov/hotel-booking-project/internal/broker"
	"log"
)

func main() {
	topic := "booking-notifications"
	admin, err := broker.New("localhost:19092")
	if err != nil {
		log.Print("initializing redpanda client err:", err)
		return
	}
	defer admin.Close()

	ok, err := admin.TopicExists(context.Background(), topic)
	if err != nil {
		log.Print("checking topic err:", err)
		return
	}
	if !ok {
		err = admin.CreateTopic(context.Background(), topic)
	}

	consumer, err := broker.NewConsumer("localhost:19092", topic)
	if err != nil {
		log.Print("initializing redpanda client err:", err)
		return
	}
	consumer.ReadNotifications()
}
