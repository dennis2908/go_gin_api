package main

import (
	"context"
	"go-restapi-gin/models"
	"go-restapi-gin/mongoconn"
	"go-restapi-gin/structs"
	"os"

	"fmt"
	"log"
	"runtime"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/joho/godotenv"
)

func main() {

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

	err := godotenv.Load("../app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	FailOnError(err, "Failed to connect to RabbitMQ")
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
			post := structs.InsertCustomer{
				IdCustomer: string(d.Body),
				Operation:  "mongo insert data customer",
			}
			db, err := mongoconn.Connect()
			if err != nil {
				log.Fatal(err.Error())
			}

			var ctx = context.TODO()

			// Insert ke database
			_, errx := db.Collection("Customers").InsertOne(ctx, post)

			// Handle error
			if errx != nil {
				fmt.Printf("an error ocurred when connect to mongoDB : %v", err)
				panic(err)
			}

			fmt.Println("Proses insert berhasil...")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
