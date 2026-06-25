package httprouter

import (
	httpcategory "clothing-service/internal/adapters/inbound/http/category"
	httpclothes "clothing-service/internal/adapters/inbound/http/clothes"
	httpcolor "clothing-service/internal/adapters/inbound/http/color"
	httpcommon "clothing-service/internal/adapters/inbound/http/common"
	"clothing-service/internal/adapters/inbound/http/middleware"
	"clothing-service/internal/application/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	clothesService *services.ClothesService,
	clothesImageService *services.ClothesImageService,
	colorService *services.ColorService,
	categoryService *services.CategoryService,
	minIOPublicBaseURL string,
	jwtSecret string,
	jwtIssuer string,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", httpcommon.Health)

	clothesHandler := httpclothes.NewHandler(
		clothesService,
		clothesImageService,
		minIOPublicBaseURL,
	)

	colorHandler := httpcolor.NewHandler(colorService)
	categoryHandler := httpcategory.NewHandler(categoryService)
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, jwtIssuer)

	v1 := router.Group("/api/v1")
	publicClothes := v1.Group("/clothes")
	{
		publicClothes.GET("", clothesHandler.List)
		publicClothes.GET("/:id", clothesHandler.FindByID)
	}
	publicCategory := v1.Group("/categories")
	{
		publicCategory.GET("", categoryHandler.CategoryDropdown)
	}

	v1.Use(authMiddleware.RequireAuth())
	privateClothes := v1.Group("/clothes")
	privateClothes.Use(authMiddleware.RequireAuth())
	{
		privateClothes.POST("", authMiddleware.RequirePermission("clothes:write"), clothesHandler.Create)
		privateClothes.PUT("/:id", authMiddleware.RequirePermission("clothes:write"), clothesHandler.Update)
		privateClothes.DELETE("/:id", authMiddleware.RequirePermission("clothes:write"), clothesHandler.Delete)
		privateClothes.POST("/:id/images", authMiddleware.RequirePermission("images:write"), clothesHandler.UploadImages)
		privateClothes.DELETE("/:id/images/:image_id", authMiddleware.RequirePermission("images:write"), clothesHandler.DeleteImage)
	}

	privateColor := v1.Group("/colors")
	{
		privateColor.GET("", authMiddleware.RequirePermission("clothes:read"), colorHandler.ColorDropdown)
	}
	return router
}
