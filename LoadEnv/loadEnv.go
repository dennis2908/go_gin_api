package loadEnv

import (
	"go-restapi-gin/models"
	"log"

	"github.com/joho/godotenv"
)

func Connects() { // init instead of int

	err := godotenv.Load("app.env")
	if err != nil {
		err1 := godotenv.Load("../app.env")
		if err1 != nil {
			err2 := godotenv.Load("../../app.env")
			if err2 != nil {
				log.Fatal("Error loading .env file")
			}
		}
	}

	config, err := models.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	configx, errx := models.LoadConfig(".")
	if errx != nil {
		configxx, errxx := models.LoadConfig("../")
		if errxx != nil {
			configxxx, errxxx := models.LoadConfig("../../")
			if errxxx != nil {
				log.Fatal("Error loading .env file")
			} else {
				models.ConnectDB(&configxxx)
			}
		} else {
			models.ConnectDB(&configxx)
		}
	} else {
		models.ConnectDB(&configx)
	}

	models.ConnectDB(&config)

}
