package pusherconn

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	pusher "github.com/pusher/pusher-http-go/v5"
)

var (
	pusherClient *pusher.Client
)

func Connect() (pusher.Client, error) {
	err := godotenv.Load("../app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pusherClient := pusher.Client{

		AppID: os.Getenv("pusher_appId"),

		Key: os.Getenv("pusher_key"),

		Secret: os.Getenv("pusher_secret"),

		Cluster: os.Getenv("pusher_cluster"),

		Secure: true,
	}
	return pusherClient, nil

}
