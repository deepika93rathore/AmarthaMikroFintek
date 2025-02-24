package routes

import (
	"github.com/gin-gonic/gin"
	"billing-engine/handler"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/loan/:id/outstanding", handler.GetOutstanding)
	r.GET("/loan/:id/delinquent", handler.IsDelinquent)
	r.POST("/loan/:id/payment", handler.MakePayment)
}
