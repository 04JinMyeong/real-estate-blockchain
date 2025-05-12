package routes

import (
	"realestate/handler"

	"github.com/gin-gonic/gin"
)

func BrokerRoutes(r *gin.Engine) {
  b := r.Group("/api/broker")
  b.POST("/register", handler.RegisterBroker)
  b.POST("/verify", handler.VerifyBroker)
}
