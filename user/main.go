package main

import (
	"time"

	"github.com/Notify-FHICT/microservices/user/endpoints"
	"github.com/Notify-FHICT/microservices/user/service"
	"github.com/Notify-FHICT/microservices/user/service/storage"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func main() {

	recordMetrics()

	db, err := storage.NewMongoDB()
	if err != nil {
		panic(err)
	}

	srv := service.NewUserService(db)
	api := endpoints.NewAPIHandler(srv)

	api.Server() //run! :D

	select {}

}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(3 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "processed_ops_3s",
		Help: "Amount of requests made in the last 3 seconds",
	})
)
