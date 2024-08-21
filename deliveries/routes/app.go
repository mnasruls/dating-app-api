package routes

import (
	"dating-app-api/configs"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Build(route fiber.Router, env configs.EnviConfig, db *gorm.DB) {
	BuildUserRoute(route, env, db)
	BuidAuthRoute(route, env, db)
}
