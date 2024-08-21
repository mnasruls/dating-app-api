package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"dating-app-api/configs"
	"dating-app-api/entities/responses"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifySignature(env *configs.EnviConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		timeStamp := c.Get("timestamp")
		if timeStamp == "" {
			return c.Status(401).JSON(responses.Response{
				StatusCode: 401,
				Message:    "timestamp required",
			})
		}
		_, err := time.Parse("2006-01-02 15:04:05", timeStamp)
		if err != nil {
			return c.Status(401).JSON(responses.Response{
				StatusCode: 401,
				Message:    "invalid timestamp",
			})
		}
		signature := c.Get("signature")
		if signature == "" {
			return c.Status(401).JSON(responses.Response{
				StatusCode: 401,
				Message:    "signature required",
			})
		}
		body := string(c.Body())
		generatedSignature := generateSignature(env.ApiKey, body, timeStamp)
		if generatedSignature != signature {
			return c.Status(401).JSON(responses.Response{
				StatusCode: 401,
				Message:    "invalid signature",
			})
		}
		return c.Next()
	}
}

func generateSignature(apiKey, body, time string) string {
	bodyToEnc := body + ":" + time

	h := hmac.New(sha256.New, []byte(apiKey))
	h.Write([]byte(bodyToEnc))
	hash := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}
