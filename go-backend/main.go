package main

import (
	"flag"
	"log"
	"redis-thumb/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	api_gw_server_address := flag.String("api_gw_server_address", "127.0.0.1:8000", "Api Gateway server address. Default:127.0.0.1:8000")
	flag.Parse()
	router := gin.Default()

	// Create the API handler and register handlers
	apiHandler := handler.NewApiHandler(router)
	_, err := apiHandler.RegisterApiHandlers()
	if err != nil {
		log.Fatal("Unable to register/create APIs callbacks")
	}

	// Run the Gin server
	router.Run(*api_gw_server_address)
}
