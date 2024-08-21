package main

import (
	"dating-app-api/configs"
	"dating-app-api/deliveries/routes"
	"dating-app-api/deliveries/validators"
	"dating-app-api/entities/responses"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	env, errs := configs.InitEnv()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		log.Fatalln("error init env")
	}

	db, err := configs.InitDb(env)
	if err != nil {
		log.Fatalln("error init db :", err)
	}

	err = configs.InitMigrations(&env)
	if err != nil {
		log.Fatalln("error run migrations :", err.Error())
	}

	validators.AddValidatorLibs()

	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   10 * 1024 * 1024,
	})

	fiberApp.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			AllowMethods: "POST, GET, OPTIONS, PUT, DELETE",
		},
	))
	// Configure rate limiter middleware
	fiberApp.Use(limiter.New(limiter.Config{
		Expiration: 30 * time.Second, // Duration of the rate limit window
		Max:        5,                // Maximum number of requests per IP during the rate limit window
		// LimiterMiddleware: limiter.ConfigDefault.LimiterMiddleware,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Use IP address as the unique identifier for rate limiting
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(responses.Response{
				StatusCode: fiber.StatusTooManyRequests,
				Message:    "Too many request",
			})
		},
	}))

	route := fiberApp.Group(fmt.Sprintf("/api/%s/", env.AppVersion))
	route.Use(logger.New(logger.Config{
		Format: `{"host":"${host}","pid":"${pid}","time":"${time}","request-id":"${locals:requestid}","status":"${status}","method":"${method}","latency":"${latency}","path":"${path}",` +
			`"user-agent":"${ua}","response-body":"${resBody}"}` + "\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Jakarta",
	}))

	// init app
	routes.Build(route, env, db)

	err = fiberApp.Listen(fmt.Sprintf("%v:%v", env.AppHost, env.AppPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err.Error())
	}
}
