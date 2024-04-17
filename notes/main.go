package main

import (
	"fmt"

	"github.com/Notify-FHICT/microservices/notes/endpoints"
	"github.com/Notify-FHICT/microservices/notes/storage"
)

func main() {
	// Initialize MongoDB database
	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	// Initialize API handler with the database
	api := endpoints.NewAPIHandler(db)

	// Start the server
	fmt.Println("Starting server")
	api.Server()

	// Block indefinitely to keep the program running
	select {}
}
