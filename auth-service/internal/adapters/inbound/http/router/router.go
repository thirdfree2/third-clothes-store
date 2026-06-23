package httprouter

import (
	httpauth "auth-service/internal/adapters/inbound/http/auth"
	httpcommon "auth-service/internal/adapters/inbound/http/common"
	"auth-service/internal/application/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authService *services.AuthService,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", httpcommon.Health)

	authHandler := httpauth.NewHandler(authService)
	v1 := router.Group("/api/v1")
	{
		admin := v1.Group("/auth")
		{
			admin.POST("/register", authHandler.Register)
			admin.POST("/login", authHandler.Login)
		}
	}

	return router
}
