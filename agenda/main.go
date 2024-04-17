package main

import (
	"fmt"

	"github.com/Notify-FHICT/microservices/agenda/endpoints"
	"github.com/Notify-FHICT/microservices/agenda/storage"
)

func main() {
	// Establish connection to MongoDB database
	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	// Create API and RMQ handlers with the connected database
	api := endpoints.NewAPIHandler(db)
	rmq := endpoints.NewRMQHandler(db)

	// Start the RabbitMQ message bus in a separate goroutine
	go rmq.MessageBus()

	// Start the API server
	fmt.Println("starting server")
	api.Server()

	// Block indefinitely to keep the main goroutine alive
	select {}
}
