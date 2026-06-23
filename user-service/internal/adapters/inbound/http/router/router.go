package httprouter

import (
	httpcart "user-service/internal/adapters/inbound/http/cart"
	httpcommon "user-service/internal/adapters/inbound/http/common"
	httpcustomer "user-service/internal/adapters/inbound/http/customer"
	"user-service/internal/adapters/inbound/http/middleware"
	httppurchase "user-service/internal/adapters/inbound/http/purchase"
	"user-service/internal/application/ports"
	"user-service/internal/application/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	customerService *services.CustomerService,
	cartService *services.CartService,
	purchaseService *services.PurchaseService,
	clothesClient ports.ClothesClient,
	jwtSecret string,
	jwtIssuer string,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", httpcommon.Health)

	customerHandler := httpcustomer.NewHandler(customerService)
	cartHandler := httpcart.NewHandler(cartService, clothesClient)
	purchaseHandler := httppurchase.NewHandler(purchaseService)
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, jwtIssuer)

	v1 := router.Group("/api/v1")
	{
		me := v1.Group("/me")
		me.Use(authMiddleware.RequireAuth())
		{
			me.GET("", customerHandler.GetMe)
			me.PATCH("", customerHandler.UpdateMe)
		}

		cart := v1.Group("/cart")
		cart.Use(authMiddleware.RequireAuth())
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddItem)
			cart.PATCH("/items/:itemId", cartHandler.UpdateItem)
			cart.DELETE("/items/:itemId", cartHandler.DeleteItem)
			cart.DELETE("", cartHandler.ClearCart)
		}

		purchases := v1.Group("/purchases")
		purchases.Use(authMiddleware.RequireAuth())
		{
			purchases.GET("", purchaseHandler.ListPurchases)
			purchases.POST("", purchaseHandler.CreatePurchase)
		}
	}

	return router
}
