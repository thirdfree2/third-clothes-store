package main

import (
	httprouter "auth-service/internal/adapters/inbound/http/router"
	jwtservice "auth-service/internal/adapters/outbound/jwt"
	"auth-service/internal/adapters/outbound/password"
	postgresdb "auth-service/internal/adapters/outbound/postgres/db"
	postgresuser "auth-service/internal/adapters/outbound/postgres/user"
	"auth-service/internal/application/services"
	"auth-service/internal/config"
	"log"

	"github.com/gin-gonic/gin"
)

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

	passwordHasher := password.NewBcryptHasher()
	userRepo := postgresuser.NewRepository(db)
	tokenService := jwtservice.NewTokenService(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTAccessTokenTTL(),
	)

	authService := services.NewAuthService(
		userRepo,
		passwordHasher,
		tokenService,
	)

	router := httprouter.NewRouter(authService)

	log.Printf("server listening on %s", cfg.HTTPAddr())
	if err := router.Run(cfg.HTTPAddr()); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
