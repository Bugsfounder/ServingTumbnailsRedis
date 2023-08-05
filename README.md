# ServingTumbnailsRedis
Explore redis and serve the thumbnail. Basically we need to check how the code should be serving a folder which contains thumbnail of the camera feed The change should take place from the backend, and then an event can be posted/pushed over redis with update filename The UI should listen and pickup the filename and refresh the image inside itself


# Writing Redis Pub/Sub in Go

## Introduction

Redis is an open-source, in-memory data structure store that can be used as a database, cache, and message broker. Redis supports the Publish/Subscribe (Pub/Sub) messaging pattern, which allows for communication between different parts of an application using channels. In this guide, we'll explore how to write a Redis Pub/Sub implementation in the Go programming language.

## Prerequisites

Before you begin, make sure you have the following prerequisites:

- Basic understanding of the Go programming language.
- Redis server installed and running locally or on a remote host.

## Setting Up Redis Client in Go

To interact with Redis in Go, you'll need a Redis client library. One popular choice is the `go-redis` library, which provides a high-level Redis client. You can install it using the following command:

```bash
go get github.com/go-redis/redis/v8
```

Replace `v8` with the appropriate version if needed.

## Writing the Publisher

Let's start by writing the publisher code, which will send messages to a Redis channel.

```go
package main

import (
	"context"
	"fmt"
	"time"

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

	// Publish messages to a channel
	channel := "my-channel"
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message %d", i)
		err := client.Publish(ctx, channel, message).Err()
		if err != nil {
			fmt.Println("Error publishing:", err)
		} else {
			fmt.Println("Published:", message)
		}
		time.Sleep(time.Second)
	}
}
```

## Writing the Subscriber

Next, let's create the subscriber code, which listens for messages on a Redis channel.

```go
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
	channel := "my-channel"
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
```

## Running the Publisher and Subscriber

1. Start your Redis server if it's not already running.
2. Run the publisher and subscriber programs in two separate terminals using `go run`.

```bash
# Terminal 1
go run publisher.go

# Terminal 2
go run subscriber.go
```

You should see the publisher sending messages, and the subscriber receiving and printing those messages.

## Conclusion

Congratulations! You've successfully implemented a Redis Pub/Sub system in Go. This allows you to establish communication between different parts of your application using channels. You can further enhance this basic implementation by adding error handling, authentication, and optimizations based on your specific use case.