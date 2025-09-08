package routers

import (
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {

	RouterCustomer(r)
	RouterProduct(r)
	RouterOrder(r)
}
