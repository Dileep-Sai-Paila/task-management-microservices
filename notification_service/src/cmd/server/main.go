package main

import (
	"context"
	"log"
	"notification_service/internal/adaptors/redis"
	"notification_service/internal/config"
)

func main() {
	log.Println("Starting Notification Service...")

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	subscriber, err := redis.NewSubscriber(cfg.RedisAddress)
	if err != nil {
		log.Fatalf("could not create redis subscriber: %v", err)
	}

	subscriber.Listen(context.Background(), "task_notifications")
}
