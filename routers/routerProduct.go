package routers

import (
	"go-restapi-gin/controllers/productcontroller"

	"go-restapi-gin/middlewares"

	"github.com/gin-gonic/gin"
)

func RouterProduct(r *gin.Engine) {

	protected := r.Group("")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/api/products", productcontroller.Index)
	protected.GET("/api/products/:id", productcontroller.Show)
	protected.POST("/api/products", productcontroller.Create)
	protected.PUT("/api/products/:id", productcontroller.Update)
	protected.DELETE("/api/products", productcontroller.Delete)

}
