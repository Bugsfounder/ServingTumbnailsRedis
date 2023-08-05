package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Ping the Redis server to ensure connectivity
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// Subscribe to a channel
	channel := "thumbnail"
	pubSub := client.Subscribe(ctx, channel)
	defer pubSub.Close()

	// Create a channel to receive OS signals (e.g., Ctrl+C)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle incoming messages
	go func() {
		for {
			msg, err := pubSub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("Error receiving message:", err)
				return
			}
			fmt.Println("Received:", msg.Payload)
		}
	}()

	fmt.Println("Subscribed to channel:", channel)

	// Wait for a termination signal
	<-signals
	fmt.Println("Unsubscribing and exiting...")
}
