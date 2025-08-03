package routes

import (
	"github.com/anggitaputra13/mcs-auth/internal/interfaces/api/handlers"
	"github.com/anggitaputra13/mcs-auth/internal/interfaces/api/middlewares"
	"github.com/anggitaputra13/mcs-auth/pkg/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwt *auth.JWT) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", middlewares.AuthMiddleware(jwt), authHandler.Logout)
		authGroup.GET("/validate", middlewares.AuthMiddleware(jwt), authHandler.Validate)
	}
}
