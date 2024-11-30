package main

import (
	"github.com/IvanChumakov/hotel-booking-project/notificationservice/internal/broker"
	"log"
)

func main() {
	topic := "booking-notifications"
	//admin, err := broker.New("redpanda-0:9092")
	//if err != nil {
	//	log.Print("initializing redpanda client err:", err)
	//	return
	//}
	//log.Print("connection initialized")
	//defer admin.Close()
	//
	//ok, err := admin.TopicExists(context.Background(), topic)
	//if err != nil {
	//	log.Print("checking topic err:", err)
	//	return
	//}
	//if !ok {
	//	err = admin.CreateTopic(context.Background(), topic)
	//}

	consumer, err := broker.NewConsumer("redpanda-0:9092", topic)
	if err != nil {
		log.Print("initializing redpanda client err:", err)
		return
	}
	consumer.ReadNotifications()
}
