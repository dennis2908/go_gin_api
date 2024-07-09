package main

import (
	"fmt"
	"go-restapi-gin/models"
	"go-restapi-gin/routers"
	"log"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {

	numberOfCores := runtime.NumCPU()
	fmt.Println(numberOfCores)
	runtime.GOMAXPROCS(numberOfCores)
	var wg sync.WaitGroup
	for i := 0; i < numberOfCores; i++ {
		wg.Add(1)

		config, err := models.LoadConfig(".")
		if err != nil {
			log.Fatal("? Could not load environment variables", err)
		}
		r := gin.Default()
		models.ConnectDB(&config)

		routers.Router(r)

		r.Run()
	}
	wg.Wait()

}
