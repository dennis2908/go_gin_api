package routers

import (
	"go-restapi-gin/controllers/ordercontroller"

	"github.com/gin-gonic/gin"
)

func RouterOrder(r *gin.Engine) {

	r.GET("/api/orders", ordercontroller.Index)
	r.GET("/api/orders/:id", ordercontroller.Show)
	r.POST("/api/orders", ordercontroller.Create)
	r.PUT("/api/orders/:id", ordercontroller.Update)
	r.DELETE("/api/orders", ordercontroller.Delete)

}
