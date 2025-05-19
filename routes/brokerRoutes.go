package routes

import (
	"realestate/handler"

	"github.com/gin-gonic/gin"
)

func BrokerRoutes(r *gin.Engine) {
	b := r.Group("/api/brokers")
	{
		b.POST("/register-with-did", handler.SignUpBrokerAndIssueDID)
		b.POST("/verify", handler.VerifyBroker)
	}
}
