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
