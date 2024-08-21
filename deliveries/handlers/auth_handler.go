package handlers

import (
	"dating-app-api/entities/requests"
	"dating-app-api/entities/responses"
	"dating-app-api/helpers"
	"dating-app-api/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	LogOut(c *fiber.Ctx) error
}

type authHandler struct {
	service services.AuthServiceInterface
	resp    responses.CommondResponse
	db      *gorm.DB
}

func NewAuthHandler(service services.AuthServiceInterface, resp responses.CommondResponse, db *gorm.DB) AuthHandlerInterface {
	return &authHandler{
		service: service,
		resp:    resp,
		db:      db,
	}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	request := new(requests.AuthRequest)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("[authHandler][Login] parse request body error :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, err.Error()))
	}

	validate := request.ValiadateAuthLogin()
	if validate != nil {
		log.Println("[authHandler][Login] validate request body :", helpers.JsonMinify(validate))
		return c.Status(400).JSON(h.resp.StatusBadRequest(validate, "invalid validation"))
	}

	res := h.service.Login(request)
	return c.Status(res.StatusCode).JSON(res)
}

func (h *authHandler) LogOut(c *fiber.Ctx) error {
	res := h.service.LogOut(c.Context())
	return c.Status(res.StatusCode).JSON(res)
}

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	res := h.service.RefreshToken(c.Context())
	return c.Status(res.StatusCode).JSON(res)
}
