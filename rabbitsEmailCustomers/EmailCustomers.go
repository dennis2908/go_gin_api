package main

import (
	"encoding/json"
	"net/smtp"
	"os"
	"runtime"
	"sync"

	"fmt"
	"log"

	"go-restapi-gin/models"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	err := godotenv.Load("../app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	numberOfCores := runtime.NumCPU()
	fmt.Println(numberOfCores)
	runtime.GOMAXPROCS(numberOfCores)
	var wg sync.WaitGroup
	for i := 0; i < numberOfCores; i++ {
		wg.Add(1)

		config, err := models.LoadConfig("../")
		if err != nil {
			log.Fatal("? Could not load environment variables", err)
		}
		models.ConnectDB(&config)

		GetData()

	}
	wg.Wait()

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetData() {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sendEmailToCustomers", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ul := &models.Customer{}
			json.Unmarshal(d.Body, ul)

			// SMTP configuration
			if ul.Name == "" {
				ul.Name = "Customer"
			}

			username := os.Getenv("SMTP_USER")
			password := os.Getenv("SMTP_PASSWORD")
			host := os.Getenv("SMTP_HOST")
			port := os.Getenv("SMTP_PORT")

			// Subject and body
			subject := "Good Night, " + ul.Name
			body := "Hi, Good Night " + ul.Name

			// Sender and receiver
			from := os.Getenv("SMTP_FROM")
			to := []string{
				ul.Email,
			}

			// Build the message
			message := fmt.Sprintf("From: %s\r\n", from)
			message += fmt.Sprintf("To: %s\r\n", to)
			message += fmt.Sprintf("Subject: %s\r\n", subject)
			message += fmt.Sprintf("\r\n%s\r\n", body)

			// Authentication.
			auth := smtp.PlainAuth("", username, password, host)

			// Send email
			err := smtp.SendMail(host+":"+port, auth, from, to, []byte(message))
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Email sent successfully.")

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
