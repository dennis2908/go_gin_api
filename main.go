package main

import (
	"fmt"
	"go-restapi-gin/controllers/customercontroller"
	"go-restapi-gin/controllers/ordercontroller"
	"go-restapi-gin/controllers/productcontroller"
	"go-restapi-gin/models"
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

		r.GET("/api/products", productcontroller.Index)
		r.GET("/api/product/:id", productcontroller.Show)
		r.POST("/api/product", productcontroller.Create)
		r.PUT("/api/product/:id", productcontroller.Update)
		r.DELETE("/api/product", productcontroller.Delete)

		r.GET("/api/customers", customercontroller.Index)
		r.GET("/api/customers/:id", customercontroller.Show)
		r.GET("/api/customers/export-excel/:filename", customercontroller.ExcelCustomers)
		r.POST("/api/customers", customercontroller.Create)
		r.PUT("/api/customers/:id", customercontroller.Update)
		r.DELETE("/api/customers", customercontroller.Delete)

		r.GET("/api/orders", ordercontroller.Index)
		r.GET("/api/orders/:id", ordercontroller.Show)
		r.POST("/api/orders", ordercontroller.Create)
		r.PUT("/api/orders/:id", ordercontroller.Update)
		r.DELETE("/api/orders", ordercontroller.Delete)

		r.Run()
	}
	wg.Wait()

}
