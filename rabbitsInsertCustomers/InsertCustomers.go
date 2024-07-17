package main

import (
	"context"
	"encoding/json"
	"go-restapi-gin/celery"
	"go-restapi-gin/models"

	"golang.org/x/crypto/bcrypt"

	"os"
	"strconv"

	"fmt"
	pusher "go-restapi-gin/pusherconn"

	"log"
	"runtime"
	"sync"
	"time"

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

func failOnError(err error, msg string) string {
	if err != nil {
		return "Error"
	} else {
		return ""
	}
}

func CreateCustomersMongoMessage(idKonsumen string) string {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createCustomersMongo", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(idKonsumen),
		})
	msg = failOnError(err, "Failed to publish a message")
	if msg == "Error" {
		return "Error"
	}
	log.Printf(" [x] Sent %s\n", idKonsumen)

	return "success"

}

func saveData(data []byte) string {
	var customers []models.Customer

	models.DB.Find(&customers)

	ul := &models.Customer{}
	json.Unmarshal(data, ul)

	//turn password into hash
	hashedPassword, errx := bcrypt.GenerateFromPassword([]byte(ul.Password), bcrypt.DefaultCost)
	if errx != nil {
		return "error"
	}
	ul.Password = string(hashedPassword)

	Qry := models.Customer{Email: ul.Email, Name: ul.Name, Password: ul.Password, UserName: ul.UserName}

	models.DB.Create(&Qry)

	fmt.Println(Qry.Id)

	idStr := strconv.Itoa(int(Qry.Id))

	CreateCustomersMongoMessage(idStr)

	return "success"
}

func GetData() {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createCustomers", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
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
			cli, _ := celery.Connect()

			saveData := saveData(d.Body)
			cli.Register("save.data.customers", saveData)

			// start workers (non-blocking call)
			cli.StartWorker()

			// wait for client request
			time.Sleep(4 * time.Second)

			// stop workers gracefully (blocking call)
			cli.StopWorker()

			client, _ := pusher.Connect()

			client.Trigger("trigger.load", "konsumens", "all.konsumens")

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
