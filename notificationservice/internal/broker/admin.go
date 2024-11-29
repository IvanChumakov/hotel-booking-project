package broker

import (
	"context"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
)

type Admin struct {
	client *kadm.Client
}

func New(broker string) (*Admin, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
	)
	if err != nil {
		return nil, err
	}
	adminClient := kadm.NewClient(client)
	return &Admin{client: adminClient}, nil
}

func (a *Admin) TopicExists(ctx context.Context, topic string) (bool, error) {
	topics, err := a.client.ListTopics(ctx)
	if err != nil {
		return false, err
	}
	for _, topicData := range topics {
		if topicData.Topic == topic {
			log.Print("topic found")
			return true, nil
		}
	}
	log.Print("there is no topic")
	return false, nil
}

func (a *Admin) CreateTopic(ctx context.Context, topic string) error {
	resp, err := a.client.CreateTopic(ctx, 1, 1, nil, topic)
	if err != nil {
		log.Print("error creating topic, ", err.Error())
		return err
	}
	if resp.Err != nil {
		log.Print("error creating topic", resp.Err)
		return resp.Err
	}
	log.Print("topic created")
	return nil
}

func (a *Admin) Close() {
	a.client.Close()
}
