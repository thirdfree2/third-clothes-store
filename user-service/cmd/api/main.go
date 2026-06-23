package main

import (
	"log"
	httprouter "user-service/internal/adapters/inbound/http/router"
	httpclothing "user-service/internal/adapters/outbound/http/clothing"
	postgrescart "user-service/internal/adapters/outbound/postgres/cart"
	postgrescustomerprofile "user-service/internal/adapters/outbound/postgres/customer_profile"
	postgresdb "user-service/internal/adapters/outbound/postgres/db"
	postgrespurchase "user-service/internal/adapters/outbound/postgres/purchase"
	"user-service/internal/application/services"
	"user-service/internal/config"

	"github.com/gin-gonic/gin"
)

// @title User Service API
// @version 1.0
func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := postgresdb.NewDB(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	customerProfileRepo := postgrescustomerprofile.NewRepository(db)
	cartRepo := postgrescart.NewRepository(db)
	purchaseRepo := postgrespurchase.NewRepository(db, cfg.MinIOPublicURL)
	clothesClient := httpclothing.NewClient(cfg.ClothingServiceBaseURL)
	customerService := services.NewCustomerService(customerProfileRepo)
	cartService := services.NewCartService(cartRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, clothesClient)

	router := httprouter.NewRouter(customerService, cartService, purchaseService, clothesClient, cfg.JWTSecret, cfg.JWTIssuer)

	log.Printf("server listening on %s", cfg.HTTPAddr())
	if err := router.Run(cfg.HTTPAddr()); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
