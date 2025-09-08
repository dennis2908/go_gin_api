package routers

import (
	"go-restapi-gin/controllers/ordercontroller"
	"go-restapi-gin/middlewares"

	"github.com/gin-gonic/gin"
)

func RouterOrder(r *gin.Engine) {

	protected := r.Group("")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/api/orders", ordercontroller.Index)
	protected.GET("/api/orders/:id", ordercontroller.Show)
	protected.POST("/api/orders", ordercontroller.Create)
	protected.PUT("/api/orders/:id", ordercontroller.Update)
	protected.DELETE("/api/orders", ordercontroller.Delete)

}
