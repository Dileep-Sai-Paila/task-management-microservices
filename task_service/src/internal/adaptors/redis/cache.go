package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func NewCache(address string) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}

	return &Cache{client: client}, nil
}

func (c *Cache) SetUserValidation(ctx context.Context, userID int32) error {

	key := fmt.Sprintf("user_validated:%d", userID)
	err := c.client.Set(ctx, key, 1, 5*time.Minute).Err() // short expiration to periodically re-validate
	if err != nil {
		return fmt.Errorf("could not set user validation in cache: %w", err)
	}
	return nil
}

// to check if a user ID is present in the cache
func (c *Cache) GetUserValidation(ctx context.Context, userID int32) (bool, error) {
	key := fmt.Sprintf("user_validated:%d", userID)
	err := c.client.Get(ctx, key).Err()

	if err == redis.Nil {
		return false, nil //cache miss only, not a real problem
	} else if err != nil {
		return false, fmt.Errorf("could not get user validation from cache: %w", err) // any other error should be a real problem.
	}

	// If we are getting here, err was nil, meaning the key exists => cache hit
	return true, nil
}

// to send a message to the task_notifications channel.
func (c *Cache) PublishTaskNotification(ctx context.Context, message string) error {
	err := c.client.Publish(ctx, "task_notifications", message).Err()
	if err != nil {
		return fmt.Errorf("could not publish task notification: %w", err)
	}
	return nil
}
