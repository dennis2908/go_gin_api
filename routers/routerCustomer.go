package routers

import (
	"go-restapi-gin/controllers/customercontroller"

	"github.com/gin-gonic/gin"
)

func RouterCustomer(r *gin.Engine) {
	r.GET("/api/customers", customercontroller.Index)
	r.GET("/api/customers/:id", customercontroller.Show)
	r.GET("/api/customers/export-excel/:filename", customercontroller.ExcelCustomers)
	r.GET("/api/customers/import-excel/:filename", customercontroller.ReadExcelKonsumens)
	r.POST("/api/customers", customercontroller.Create)
	r.PUT("/api/customers/:id", customercontroller.Update)
	r.DELETE("/api/customers", customercontroller.Delete)

}
