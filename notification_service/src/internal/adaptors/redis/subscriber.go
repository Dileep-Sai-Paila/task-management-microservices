package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type Subscriber struct {
	client *redis.Client
}

func NewSubscriber(address string) (*Subscriber, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}

	log.Printf("Successfully connected to Redis at %s", address)
	return &Subscriber{client: client}, nil
}

func (s *Subscriber) Listen(ctx context.Context, channelName string) {
	pubsub := s.client.Subscribe(ctx, channelName)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Could not subscribe to channel '%s': %v", channelName, err)
	}

	log.Printf("Subscribed to '%s' channel. Waiting for messages...", channelName)
	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("[Notification Received]: %s", msg.Payload)
	}
}
