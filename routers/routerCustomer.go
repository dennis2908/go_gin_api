package routers

import (
	"go-restapi-gin/middlewares"

	"go-restapi-gin/controllers/customercontroller"

	"github.com/gin-gonic/gin"
)

func RouterCustomer(r *gin.Engine) {
	protected := r.Group("")
	protected.Use(middlewares.JwtAuthMiddleware())
	r.GET("/api/customers", customercontroller.Index)
	protected.GET("/api/customers/:id", customercontroller.Show)
	protected.GET("/api/customers/export-excel/:filename", customercontroller.ExcelCustomers)
	protected.GET("/api/customers/import-excel/:filename", customercontroller.ReadExcelKonsumens)
	r.POST("/api/customers", customercontroller.Create)
	r.POST("/api/customers/login", customercontroller.Login)
	r.POST("/api/customers/refresh/token", customercontroller.GenerateRefreshToken)
	protected.PUT("/api/customers/:id", customercontroller.Update)
	protected.DELETE("/api/customers", customercontroller.Delete)

}
