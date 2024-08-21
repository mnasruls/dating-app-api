package handlers

import (
	"dating-app-api/deliveries/validators"
	"dating-app-api/entities/models"
	"dating-app-api/entities/requests"
	"dating-app-api/entities/responses"
	"dating-app-api/helpers"
	"dating-app-api/services"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandlerInterface interface {
	RegisterUser(c *fiber.Ctx) error
	GetDetailUser(c *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
	CheckUsername(c *fiber.Ctx) error
	VerifyUser(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}

type userHandler struct {
	service services.UserServiceInterface
	resp    responses.CommondResponse
	db      *gorm.DB
}

func NewUserHandler(service services.UserServiceInterface, resp responses.CommondResponse, db *gorm.DB) UserHandlerInterface {
	return &userHandler{
		service: service,
		resp:    resp,
		db:      db,
	}
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	request := new(requests.CreateUserRequest)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("[userHandler][RegisterUser] parse request body error :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, err.Error()))
	}

	validate := request.ValiadateCreateUser()
	if validate != nil {
		log.Println("[userHandler][RegisterUser] validate request body :", helpers.JsonMinify(validate))
		return c.Status(400).JSON(h.resp.StatusBadRequest(validate, "invalid validation"))
	}

	err = validators.ValidatePassword(request.Password)
	if err != nil {
		log.Println("[userHandler][RegisterUser] validate password :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(map[string]string{"password": err.Error()}, "invalid validation"))
	}

	dbTx := h.db.Begin()
	if dbTx.Error != nil {
		log.Println("[userHandler][RegisterUser] error create db transaction :", dbTx.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}

	res := h.service.RegisterUser(request, dbTx)
	if res.StatusCode != http.StatusCreated {
		roll := dbTx.Rollback()
		if roll.Error != nil {
			log.Println("[userHandler][RegisterUser] error rollback db transaction :", roll.Error)
			return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
		}
		return c.Status(res.StatusCode).JSON(res)
	}

	comm := dbTx.Commit()
	if comm.Error != nil {
		log.Println("[userHandler][RegisterUser] error commit db transaction :", comm.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}

	return c.Status(res.StatusCode).JSON(res)
}

func (h *userHandler) GetDetailUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, "invalid id"))
	}

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, "invalid id"))
	}

	res := h.service.GetDetail(id)
	return c.Status(res.StatusCode).JSON(res)
}

func (h *userHandler) GetMe(c *fiber.Ctx) error {
	me := c.Locals("metadata").(models.TokenMetaData)
	res := h.service.GetDetail(me.Id)
	return c.Status(res.StatusCode).JSON(res)
}

func (h *userHandler) CheckUsername(c *fiber.Ctx) error {
	request := new(requests.CreateUserRequest)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("[userHandler][CheckUsername] parse request body error :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, err.Error()))
	}

	validate := request.ValiadateCheckUsername()
	if validate != nil {
		log.Println("[userHandler][CheckUsername] validate request body :", helpers.JsonMinify(validate))
		return c.Status(400).JSON(h.resp.StatusBadRequest(validate, "invalid validation"))
	}

	res := h.service.CheckUsername(request.Username)
	return c.Status(res.StatusCode).JSON(res)
}

func (h *userHandler) VerifyUser(c *fiber.Ctx) error {
	request := new(requests.VerifyUser)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("[userHandler][VerifyUser] parse request body error :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, err.Error()))
	}

	validate := request.VerifyUser()
	if validate != nil {
		log.Println("[userHandler][VerifyUser] validate request body :", helpers.JsonMinify(validate))
		return c.Status(400).JSON(h.resp.StatusBadRequest(validate, "invalid validation"))
	}

	dbTx := h.db.Begin()
	if dbTx.Error != nil {
		log.Println("[userHandler][VerifyUser] error create db transaction :", dbTx.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}

	res := h.service.VerifyUser(request, dbTx)
	if res.StatusCode != http.StatusCreated {
		roll := dbTx.Rollback()
		if roll.Error != nil {
			log.Println("[userHandler][VerifyUser] error rollback db transaction :", roll.Error)
			return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
		}
		return c.Status(res.StatusCode).JSON(res)
	}

	comm := dbTx.Commit()
	if comm.Error != nil {
		log.Println("[userHandler][VerifyUser] error commit db transaction :", comm.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}
	return c.Status(res.StatusCode).JSON(res)

}

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	request := new(requests.AuthRequest)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("[userHandler][ChangePassword] parse request body error :", err)
		return c.Status(400).JSON(h.resp.StatusBadRequest(nil, err.Error()))
	}

	validate := request.ValiadateChangePassword()
	if validate != nil {
		log.Println("[userHandler][ChangePassword] validate request body :", helpers.JsonMinify(validate))
		return c.Status(400).JSON(h.resp.StatusBadRequest(validate, "invalid validation"))
	}

	dbTx := h.db.Begin()
	if dbTx.Error != nil {
		log.Println("[userHandler][ChangePassword] error create db transaction :", dbTx.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}

	res := h.service.ChangePassword(c.Context(), request.Password, dbTx)
	if res.StatusCode != http.StatusOK {
		roll := dbTx.Rollback()
		if roll.Error != nil {
			log.Println("[userHandler][ChangePassword] error rollback db transaction :", roll.Error)
			return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
		}
		return c.Status(res.StatusCode).JSON(res)
	}

	comm := dbTx.Commit()
	if comm.Error != nil {
		log.Println("[userHandler][ChangePassword] error commit db transaction :", comm.Error)
		return c.Status(500).JSON(h.resp.StatusServerError("something went wrong"))
	}
	return c.Status(res.StatusCode).JSON(res)

}
