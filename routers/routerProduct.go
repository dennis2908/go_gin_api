package routers

import (
	"go-restapi-gin/controllers/productcontroller"

	"github.com/gin-gonic/gin"
)

func RouterProduct(r *gin.Engine) {

	r.GET("/api/products", productcontroller.Index)
	r.GET("/api/products/:id", productcontroller.Show)
	r.POST("/api/products", productcontroller.Create)
	r.PUT("/api/products/:id", productcontroller.Update)
	r.DELETE("/api/products", productcontroller.Delete)

}
