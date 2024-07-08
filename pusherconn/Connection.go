package pusherconn

import (
	pusher "github.com/pusher/pusher-http-go/v5"
)

var (
	pusherClient *pusher.Client
)

func Connect() (pusher.Client, error) {
	pusherClient := pusher.Client{

		AppID: "705473",

		Key: "b9e4d6190581d989a6e2",

		Secret: "629f95f4aa4563d80845",

		Cluster: "ap1",

		Secure: true,
	}
	return pusherClient, nil

}
