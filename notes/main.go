package main

import (
	"fmt"

	"github.com/Notify-FHICT/microservices/notes/endpoints"
	"github.com/Notify-FHICT/microservices/notes/storage"
)

func main() {

	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	api := endpoints.NewAPIHandler(db)

	fmt.Println("starting server")
	api.Server() //run! :D

	select {}
}
