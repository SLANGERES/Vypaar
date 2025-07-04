package RabbitMQ

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type data struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

func SendMail(Email string, Otp int) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Unable to dial the connection: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Unable to open channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("Email-Queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Unable to declare queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newData := data{
		Email: Email,
		Otp:   Otp,
	}
	jsonData, err := json.Marshal(newData)
	if err != nil {
		log.Fatalf("Unable to convert data to JSON: %v", err)
	}

	err = ch.PublishWithContext(
		ctx,
		"",            // exchange
		"Email-Queue", // routing key (queue name)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return "Failed to publish message"
	}

	return "Message published successfully"

}
