package main

import (
	"fmt"
	loadEnv "go-restapi-gin/LoadEnv"
	"go-restapi-gin/middlewares"
	"go-restapi-gin/routers"
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

		loadEnv.Connects()
		r := gin.Default()

		routers.Router(r)

		protected := r.Group("")
		protected.Use(middlewares.JwtAuthMiddleware())

		r.Run()
	}
	wg.Wait()

}
