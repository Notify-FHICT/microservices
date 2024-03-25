package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Notify-FHICT/microservices/agenda/storage"
	"github.com/Notify-FHICT/microservices/agenda/storage/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQHandler struct {
	c storage.DB
}

func NewRMQHandler(collection storage.DB) RMQHandler {
	return RMQHandler{
		collection,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (rmq *RMQHandler) MessageBus() {
	conn, err := amqp.Dial("amqp://guest:guest@10.101.45.75:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"Link_NoteEvent", // queue
		"",               // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              //args
	)
	failOnError(err, "Failed to declare a queue")

	// print consumed messages from queue
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
			var obj models.UpdateNoteID
			buf := bytes.NewBuffer(msg.Body)
			decoder := json.NewDecoder(buf)
			err := decoder.Decode(&obj)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(obj)
				err := rmq.c.LinkNoteID(obj)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("YIPPIE")
			}

		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
