package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/joho/godotenv"
)

type Data struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

func main() {
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

	q, err := ch.QueueDeclare("Email-Queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Unable to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		log.Fatalf("Unable to consume data: %v", err)
	}

	log.Println("Waiting for messages...")

	// Use a goroutine to read messages
	go func() {
		for d := range msgs {
			var data Data
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Failed to decode JSON: %v", err)
				continue
			}
			log.Printf("Received Email: %s, OTP: %d", data.Email, data.Otp)
			SendEmail(data.Email, data.Otp)
		}
	}()

	// Block forever
	select {}
}

// / SendEmail sends an email using SMTP
func SendEmail(to string, otp int) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	host := "smtp.gmail.com"
	port := "587"

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP is: %d", otp)

	// Construct full email message
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", from, password, host)

	err = smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v\n", to, err)
		return
	}
	log.Printf("Email successfully sent to %s\n", to)
}
