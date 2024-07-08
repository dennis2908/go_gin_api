package loadEnv

import (
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

}
