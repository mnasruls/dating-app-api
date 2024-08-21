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

func BuildUserRoute(route fiber.Router, env configs.EnviConfig, db *gorm.DB) {
	common := responses.NewResponseAPI()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, *common, env.Redis, &env)
	userHandler := handlers.NewUserHandler(userService, *common, db)

	userVerify := middlewares.UserVerify(&env)
	signatureVerify := middlewares.VerifySignature(&env)

	route.Post("/user/register", signatureVerify, userHandler.RegisterUser)
	route.Post("/user/verify", signatureVerify, userHandler.VerifyUser)
	route.Post("/user/check-username", signatureVerify, userHandler.CheckUsername)
	route.Put("/user/change-password", userVerify, userHandler.ChangePassword)
	route.Get("/user/detail/:id", userVerify, userHandler.GetDetailUser)
	route.Get("/user/me", userVerify, userHandler.GetMe)
}
