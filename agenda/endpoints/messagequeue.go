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

// RMQHandler handles RabbitMQ message consumption
type RMQHandler struct {
	c storage.DB
}

// NewRMQHandler creates a new RMQHandler with the provided collection
func NewRMQHandler(collection storage.DB) RMQHandler {
	return RMQHandler{
		collection,
	}
}

// failOnError is a utility function to log and panic on errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// MessageBus starts consuming messages from RabbitMQ
func (rmq *RMQHandler) MessageBus() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@10.101.45.75:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	msgs, err := ch.Consume(
		"Link_NoteEvent", // queue
		"",               // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              // args
	)
	failOnError(err, "Failed to declare a queue")

	// Print consumed messages from the queue
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
				if obj.ID.IsZero() {
					err := rmq.c.UnlinkNoteID(obj)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					err := rmq.c.LinkNoteID(obj)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
