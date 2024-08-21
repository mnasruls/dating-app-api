package services

import (
	"context"
	"dating-app-api/configs"
	middlewares "dating-app-api/deliveries/middlewares"
	"dating-app-api/entities/models"
	"dating-app-api/entities/requests"
	"dating-app-api/entities/responses"
	"dating-app-api/helpers"
	"dating-app-api/repositories"
	"dating-app-api/utils"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Login(req *requests.AuthRequest) responses.Response
	RefreshToken(ctx context.Context) responses.Response
	LogOut(ctx context.Context) responses.Response
}

type authService struct {
	userRepo  repositories.UserRepositoryInterface
	common    responses.CommondResponse
	redisUtil *utils.Redis
	envs      *configs.EnviConfig
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, common responses.CommondResponse, redisUtil *utils.Redis, envs *configs.EnviConfig) AuthServiceInterface {
	return &authService{
		userRepo:  userRepo,
		common:    common,
		redisUtil: redisUtil,
		envs:      envs,
	}
}

func (service *authService) Login(req *requests.AuthRequest) responses.Response {

	whereClause := map[string]interface{}{
		"username": req.Username,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[authService][Login] error get detail user :", err)
		return service.common.StatusServerError(err.Error())
	}

	if user == nil {
		log.Println("[authService][Login] username not found")
		return service.common.StatusBadRequest(nil, "username or password wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println("[authService][Login] password not valid")
		return service.common.StatusBadRequest(nil, "username or password wrong")
	}

	var userResponse responses.UserResponse
	err = helpers.Unmarshal(user, &userResponse)
	if err != nil {
		log.Println("[authService][Login] error unmarshal user model to responses :", err)
		return service.common.StatusServerError("something went wrong")
	}

	meta := models.TokenMetaData{
		Id:     user.Id,
		Verify: user.Verified,
		RtId:   "rt",
	}

	token, err := middlewares.GenerateToken(service.envs, meta, false)
	if err != nil {
		log.Println("[authService][Login] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	rToken, err := middlewares.GenerateToken(service.envs, meta, true)
	if err != nil {
		log.Println("[authService][Login] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	authResponse := responses.AuthResponse{
		User:         &userResponse,
		AccessToken:  token,
		RefreshToken: rToken,
	}

	return service.common.StatusOk(authResponse, nil, "login successfully")
}

func (service *authService) RefreshToken(ctx context.Context) responses.Response {
	meta := ctx.Value("metadata").(models.TokenMetaData)
	whereClause := map[string]interface{}{
		"id": meta.Id,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[authService][RefreshToken] error get detail user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	var userResponse responses.UserResponse
	err = helpers.Unmarshal(user, &userResponse)
	if err != nil {
		log.Println("[authService][RefreshToken] error unmarshal user model to responses :", err)
		return service.common.StatusServerError("something went wrong")
	}

	newMeta := models.TokenMetaData{
		Id:     user.Id,
		Verify: user.Verified,
		RtId:   "rt",
	}

	token, err := middlewares.GenerateToken(service.envs, newMeta, false)
	if err != nil {
		log.Println("[authService][RefreshToken] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	rToken, err := middlewares.GenerateToken(service.envs, newMeta, true)
	if err != nil {
		log.Println("[authService][RefreshToken] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	if err := service.redisUtil.DeleteDataFromRedis(fmt.Sprintf("metart:%v:%v", meta.Id, meta.RtId)); err != nil {
		log.Println("[authService][RefreshToken] error delete previous token :", err)
		return service.common.StatusServerError("something went wrong")

	}

	authResponse := responses.AuthResponse{
		User:         &userResponse,
		AccessToken:  token,
		RefreshToken: rToken,
	}

	return service.common.StatusOk(authResponse, nil, "refresh token successfully")
}

func (service *authService) LogOut(ctx context.Context) responses.Response {
	meta := ctx.Value("metadata").(models.TokenMetaData)
	err := service.redisUtil.DeleteDataFromRedis(fmt.Sprintf("metaat:%v", meta.Id))
	if err != nil {
		log.Println("[authService][LogOut] error delete previous token :", err)
		return service.common.StatusServerError("something went wrong")
	}
	err = service.redisUtil.DeleteDataFromRedis(fmt.Sprintf("metaat:%v:%v", meta.Id, meta.RtId))
	if err != nil {
		log.Println("[authService][LogOut] error delete previous token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	return service.common.StatusOk(nil, nil, "logout successfully")
}
