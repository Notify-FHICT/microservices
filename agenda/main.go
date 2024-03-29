package main

import (
	"fmt"

	"github.com/Notify-FHICT/microservices/agenda/endpoints"
	"github.com/Notify-FHICT/microservices/agenda/storage"
)

func main() {

	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	api := endpoints.NewAPIHandler(db)
	rmq := endpoints.NewRMQHandler(db)

	go rmq.MessageBus()

	fmt.Println("starting server")
	api.Server() //run! :D

	select {}
}
