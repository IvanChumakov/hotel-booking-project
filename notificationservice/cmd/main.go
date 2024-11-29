package main

import (
	"context"
	broker2 "github.com/IvanChumakov/hotel-booking-project/notificationservice/internal/broker"
	"log"
)

func main() {
	topic := "booking-notifications"
	admin, err := broker2.New("localhost:19092")
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

	consumer, err := broker2.NewConsumer("localhost:19092", topic)
	if err != nil {
		log.Print("initializing redpanda client err:", err)
		return
	}
	consumer.ReadNotifications()
}
