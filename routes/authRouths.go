package routes

import (
	"go-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}
}
