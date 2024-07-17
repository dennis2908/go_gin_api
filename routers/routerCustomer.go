package routers

import (
	"go-restapi-gin/controllers/customercontroller"
	"go-restapi-gin/middlewares"

	"github.com/gin-gonic/gin"
)

func RouterCustomer(r *gin.Engine) {
	protected := r.Group("")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/api/customers", customercontroller.Index)
	protected.GET("/api/customers/:id", customercontroller.Show)
	protected.GET("/api/customers/export-excel/:filename", customercontroller.ExcelCustomers)
	protected.GET("/api/customers/import-excel/:filename", customercontroller.ReadExcelKonsumens)
	protected.POST("/api/customers", customercontroller.Create)
	r.POST("/api/customers/login", customercontroller.Login)
	protected.PUT("/api/customers/:id", customercontroller.Update)
	protected.DELETE("/api/customers", customercontroller.Delete)

}
