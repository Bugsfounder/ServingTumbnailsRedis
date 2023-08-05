# Documentation: Redis Pub/Sub Integration

This documentation provides an overview of a codebase that demonstrates the integration of Redis Publish/Subscribe (Pub/Sub) mechanism in a Go backend (API publisher) and a Next.js frontend (API subscriber). The Pub/Sub mechanism allows real-time communication between different parts of an application by broadcasting messages to multiple subscribers.

### Ensure Redis Server is Running
Before proceeding, ensure that a Redis server is running on your device. If you don't have Redis set up, follow these steps:

1. Pull the Redis Docker image:
   ```bash
   docker pull redis
   ```

2. Run the Redis Docker container:
   ```bash
   docker run --name redis-server -d -p 6379:6379 redis
   ```

This will set up a Redis server using Docker. Make sure to follow these steps if you need a Redis server for your application.

## 1. Pub/Sub Concept

Redis Pub/Sub is a messaging paradigm where a publisher sends messages to a specific channel, and multiple subscribers listen to that channel to receive and react to messages. It enables asynchronous communication between different components of an application.

## 2. Go Backend (API Publisher)

### `handler/handler.go`
```go
// handler/handler.go

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	server_ctxt *gin.Engine
}

func NewApiHandler(router *gin.Engine) *ApiHandler {
	return &ApiHandler{
		server_ctxt: router,
	}
}

func (api *ApiHandler) RegisterApiHandlers() (int, error) {
	api.server_ctxt.GET("/api", api.Index)

	// go api.PublishRandomImages()

	return 1, nil
}

func (api *ApiHandler) Index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World")
}

// func (api *ApiHandler) PublishRandomImages() {

// }
```

This package defines an API handler in a Go backend. It uses the Gin web framework for HTTP handling.

- `ApiHandler`: Struct that holds a reference to the Gin engine.
- `NewApiHandler`: Initializes a new `ApiHandler` instance.
- `RegisterApiHandlers`: Registers API routes (endpoints) and handlers.
- `Index`: Handles requests to the root API endpoint ("/api") and responds with "Hello World".

### `publisher/publisher.go`
```go
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
```

This package handles Redis Pub/Sub functionality for publishing random images.

- `Publisher`: Struct that holds a Redis client instance.
- `NewPublisher`: Creates a new `Publisher` instance and establishes a connection to the Redis server.
- `getRandomImage`: Returns a random image URL from the provided list.
- `PublishRandomImage`: Publishes a random image URL to a specified channel in Redis.

### `main.go`
```go
// main.go

package main

import (
	"flag"
	"time"

	"redis_pub_sub/publisher"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	api_gw_server_address := flag.String("api_gw_server_address", "127.0.0.1:8000", "Api Gateway server address. Default:127.0.0.1:8000")
	flag.Parse()

	router := gin.Default()

	// Create the CORS middleware with appropriate configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://0.0.0.0:3000"} // Replace with your frontend URL
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}

	// Use the CORS middleware
	router.Use(cors.New(corsConfig))

	// Start a goroutine to periodically call testPublisher()
	go func() {
		for {
			testPublisher()
			time.Sleep(3 * time.Second)
		}
	}()

	// Run the Gin server
	router.Run(*api_gw_server_address)
}
func testPublisher() {
	// testing publisher here
	imageLinks := []string{
		"https://cdn.pixabay.com/photo/2023/05/13/14/35/white-flower-7990645_960_720.jpg",
		"https://cdn.pixabay.com/photo/2019/09/02/11/00/frog-4446995_960_720.jpg",
		"https://cdn.pixabay.com/photo/2015/07/05/13/44/beach-832346_960_720.jpg",
		"https://cdn.pixabay.com/photo/2016/10/18/21/22/beach-1751455_960_720.jpg",
		"https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg",
	}

	pub, err := publisher.NewPublisher()
	if err != nil {
		panic(err)
	}

	channel := "thumbnail"
	err = pub.PublishRandomImage(channel, imageLinks)
	if err != nil {
		panic(err)
	}
}
```
The main entry point of the Go backend application.

- Parses command-line flags to configure API server address.
- Configures CORS (Cross-Origin Resource Sharing) to allow frontend access.
- Starts a Goroutine to periodically publish random images using the `testPublisher` function.
- Starts the Gin server to listen for API requests.
### Run main.go file
```go
go run main.go
```
output
```
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8000
Random Image: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
Published: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
Random Image: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
Published: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
Random Image: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
```
## 3. Next.js Frontend (API Subscriber)

### `pages/api/subscribe.js`
```js
import { NextApiRequest, NextApiResponse } from "next";
import Redis from 'ioredis';

const handler = async (req, res) => {
  res.setHeader('Cache-Control', 'no-cache');
  res.setHeader('Content-Type', 'text/event-stream');
  res.setHeader('Connection', 'keep-alive');
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Credentials', 'true');
  res.setHeader('Access-Control-Allow-Methods', 'GET, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  const redisConfig = {
    host: 'localhost',
    port: 6379,
  };
  const redis = new Redis(redisConfig);

  // Subscribe to a channel
  await redis.subscribe('thumbnail');

  // Listen for incoming messages
  redis.on('message', (channel, message) => {
    console.log(`Received message in channel ${channel}: ${message}`);
    res.write(`data: ${message}\n\n`);
  });

  // Keep the connection open
  req.on('close', () => {
    redis.unsubscribe('thumbnail');
    redis.quit();
  });
};

export default handler;
```
This Next.js API route establishes a Server-Sent Events (SSE) connection to listen for messages published on the "thumbnail" Redis channel.

- Sets necessary headers for SSE and CORS.
- Initializes a connection to the Redis server and subscribes to the "thumbnail" channel.
- Listens for incoming messages and sends them to the client over the SSE connection.
- Closes the Redis subscription and connection when the client closes the SSE connection.

### `app/Subscriber.js`
```js
"use client"
import React, { useState, useEffect } from 'react';

const Subscriber = () => {
    const [message, setMessage] = useState('');

    useEffect(() => {
        console.log('EventSource connecting...');
        const eventSource = new EventSource('/api/subscribe');

        eventSource.onmessage = (event) => {
            console.log('EventSource connected.');
            const message = event.data;
            console.log("message", message);
            setMessage(message);
        };

        return () => {
            eventSource.close();
        };
    }, []);

    return (
        <div>
            <h2>Redis Subscription Example</h2>
            <ul>
                <p>{message}</p>
                <img src={message} alt="" />
            </ul>
        </div>
    );
};

export default Subscriber;
```
A React component that subscribes to the SSE connection and displays received messages (image URLs).

- Sets up an SSE connection to the `/api/subscribe` endpoint.
- Listens for messages and updates the component's state with the received message (image URL).
- Renders the received image URL.

### `page.js`
```js
// pages/index.js

import React from "react";
import Subscription from './Subscriber';
export default function Home() {
  return (
    <div>
      <main>
        dsfsd
        <Subscription />
      </main>
    </div>
  );
}
```

A Next.js page that renders the `Subscriber` component.
#### run react js 
```
npm run dev
```
output
```
- wait compiling /api/subscribe (client and server)...
- event compiled successfully in 518 ms (47 modules)
- warn "next" should not be imported directly, imported in /home/bugs/workspace/POC_Codes/ServingTumbnailsRedis/thumbnail-serving-frontend/.next/server/pages/api/subscribe.js
See more info here: https://nextjs.org/docs/messages/import-next
API resolved without sending a response for /api/subscribe, this may result in stalled requests.
Received message in channel thumbnail: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
Received message in channel thumbnail: https://cdn.pixabay.com/photo/2016/10/18/21/22/beach-1751455_960_720.jpg
Received message in channel thumbnail: https://cdn.pixabay.com/photo/2016/10/18/21/22/beach-1751455_960_720.jpg
Received message in channel thumbnail: https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885_960_720.jpg
```
## 4. Redis Pub/Sub Workflow

1. The Go backend (`handler/handler.go`) publishes random image URLs to the "thumbnail" Redis channel using the `PublishRandomImage` method.
2. The Next.js frontend (`pages/api/subscribe.js`) establishes an SSE connection to the `/api/subscribe` endpoint and listens for messages on the "thumbnail" channel.
3. When a new image URL is published to the Redis channel, the frontend receives the message and updates the UI to display the image.




## 5. Conclusion

This codebase demonstrates the implementation of Redis Pub/Sub to achieve real-time communication between a Go backend and a Next.js frontend. The backend publishes random images to a Redis channel, and the frontend subscribes to the channel to receive and display the images in real time. This Pub/Sub architecture enables efficient and asynchronous communication between different parts of the application.