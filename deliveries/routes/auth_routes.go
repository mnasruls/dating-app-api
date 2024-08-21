package routes

import (
	"dating-app-api/configs"
	"dating-app-api/deliveries/handlers"
	"dating-app-api/deliveries/middlewares"
	"dating-app-api/entities/responses"
	"dating-app-api/repositories"
	"dating-app-api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BuidAuthRoute(route fiber.Router, env configs.EnviConfig, db *gorm.DB) {
	common := responses.NewResponseAPI()
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, *common, env.Redis, &env)
	authHandler := handlers.NewAuthHandler(authService, *common, db)

	userVerify := middlewares.UserVerify(&env)
	userRtVerify := middlewares.RefreshTokenVerify(&env)
	signatureVerify := middlewares.VerifySignature(&env)

	route.Post("/auth/login", signatureVerify, authHandler.Login)
	route.Post("/auth/logout", userVerify, authHandler.LogOut)
	route.Post("/auth/refresh-token", userRtVerify, authHandler.RefreshToken)
}
