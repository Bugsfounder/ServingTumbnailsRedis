package publisher

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

type Publisher struct {
	client *redis.Client
}

func NewPublisher() (*Publisher, error) {
	ctx := context.Background()

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Ping the Redis server to ensure connectivity
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		client: client,
	}, nil
}

func (p *Publisher) getRandomImage(imageLinks []string) string {
	randomIndex := rand.Intn(len(imageLinks))
	return imageLinks[randomIndex]
}

func (p *Publisher) PublishRandomImage(channel string, imageLinks []string) error {
	ctx := context.Background()

	randomImage := p.getRandomImage(imageLinks)
	fmt.Println("Random Image:", randomImage)
	message := fmt.Sprintf("%v", randomImage)
	err := p.client.Publish(ctx, channel, message).Err()
	if err != nil {
		fmt.Println("Error publishing:", err)
	} else {
		fmt.Println("Published:", message)
	}
	time.Sleep(time.Second)

	return err
}
