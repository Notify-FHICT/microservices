package main

import (
	"github.com/Notify-FHICT/microservices/user/endpoints"
	"github.com/Notify-FHICT/microservices/user/service"
	"github.com/Notify-FHICT/microservices/user/service/storage"
)

func main() {

	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	srv := service.NewUserService(db)
	api := endpoints.NewAPIHandler(srv)

	api.Server() //run! :D

	select {}

}
