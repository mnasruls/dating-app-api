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
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	RegisterUser(request *requests.CreateUserRequest, tx *gorm.DB) responses.Response
	UpdateUser(request *requests.UpdateUserRequest, id string, tx *gorm.DB) responses.Response
	GetDetail(id string) responses.Response
	GetList(meta *requests.MetaPaginationRequest) responses.Response
	DeleteUser(id string, tx *gorm.DB) responses.Response
	CheckUsername(username string) responses.Response
	ChangePassword(ctx context.Context, password string, tx *gorm.DB) responses.Response
	VerifyUser(request *requests.VerifyUser, tx *gorm.DB) responses.Response
}

type userService struct {
	userRepo  repositories.UserRepositoryInterface
	common    responses.CommondResponse
	redisUtil *utils.Redis
	envs      *configs.EnviConfig
}

func NewUserService(userRepo repositories.UserRepositoryInterface, common responses.CommondResponse, redisUtil *utils.Redis, envs *configs.EnviConfig) UserServiceInterface {
	return &userService{
		userRepo:  userRepo,
		common:    common,
		redisUtil: redisUtil,
		envs:      envs,
	}
}

func (service *userService) RegisterUser(request *requests.CreateUserRequest, tx *gorm.DB) responses.Response {

	// store how many times user request
	var count requests.CountRequest
	err := service.redisUtil.RetrieveDataFromRedis(request.Username+"-attempt", &count)
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			count.Count = 1
			err = service.redisUtil.SaveDataToRedis(request.Username+"-attempt", count, time.Duration(helpers.GetTimeToMidnight())*time.Hour)
			if err != nil {
				log.Println("[userService][RegisterUser] error save attemps to redis :", err)
				return service.common.StatusServerError("something went wrong")
			}
		} else {
			log.Println("[userService][RegisterUser] error save get attempts to redis :", err)
			return service.common.StatusServerError("something went wrong")
		}
	} else {
		if count.Count > 3 {
			log.Println("[userService][RegisterUser] maximum register attempts in a day")
			return service.common.StatusBadRequest(nil, "maximum register attempts in day")
		} else {
			count.Count += 1
			err = service.redisUtil.SaveDataToRedis(request.Username+"-attempt", count, time.Duration(helpers.GetTimeToMidnight())*time.Hour)
			if err != nil {
				log.Println("[userService][RegisterUser] error save attemps to redis :", err)
				return service.common.StatusServerError("something went wrong")
			}
		}
	}

	whereClause := map[string]interface{}{
		"username": request.Username,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[userService][RegisterUser] error get detail user :", err)
		return service.common.StatusServerError(err.Error())
	}

	if user != nil {
		log.Println("[userService][RegisterUser] username already exist")
		return service.common.StatusBadRequest(nil, "username already exist")
	}

	var userModel *models.UserModel
	err = helpers.Unmarshal(request, &userModel)
	if err != nil {
		log.Println("[userService][RegisterUser] error unmarshal user request to model :", err)
		return service.common.StatusServerError(err.Error())
	}
	userModel.Verified = false

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[userService][RegisterUser] error bcrypt password :", err)
		return service.common.StatusServerError("something went wrong")
	}

	userModel.Password = string(hashedPass)

	user, err = service.userRepo.CreateUser(userModel, tx)
	if err != nil {
		log.Println("[userService][RegisterUser] error insert user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	meta := models.TokenMetaData{
		Id:     user.Id,
		Verify: false,
		RtId:   "rt",
	}
	token, err := middlewares.GenerateToken(service.envs, meta, false)
	if err != nil {
		log.Println("[userService][RegisterUser] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	rToken, err := middlewares.GenerateToken(service.envs, meta, true)
	if err != nil {
		log.Println("[userService][RegisterUser] error generate token :", err)
		return service.common.StatusServerError("something went wrong")
	}

	return service.common.StatusCreated(map[string]interface{}{
		"access_token":  token,
		"refresh_token": rToken,
	}, "register user successfully")
}

func (service *userService) UpdateUser(request *requests.UpdateUserRequest, id string, tx *gorm.DB) responses.Response {
	whereClause := map[string]interface{}{
		"id": id,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[userService][UpdateUser] error get detail user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	if user == nil {
		log.Println("[userService][UpdateUser] user not found with id", id)
		return service.common.StatusNotFound("user not found")
	}

	return service.common.StatusOk(nil, nil, "update user successfully")
}

func (service *userService) GetDetail(id string) responses.Response {
	whereClause := map[string]interface{}{
		"id": id,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[userService][CheckEmail] error get detail user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	if user == nil {
		log.Println("[userService][CheckEmail] user not found with id", id)
		return service.common.StatusNotFound("user not found")
	}

	var userResponse responses.UserResponse
	err = helpers.Unmarshal(user, &userResponse)
	if err != nil {
		log.Println("[authService][RefreshToken] error unmarshal user model to responses :", err)
		return service.common.StatusServerError("something went wrong")
	}

	return service.common.StatusOk(userResponse, nil, "get detail user successfully")
}

func (service *userService) CheckUsername(username string) responses.Response {

	// check username
	whereClause := map[string]interface{}{
		"username": username,
	}

	user, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[userService][CheckUsername] error get detail user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	if user != nil {
		log.Println("[userService][CheckUsername] username already exist")
		return service.common.StatusBadRequest(nil, "username already exist")
	}

	return service.common.StatusOk(nil, nil, "username valid")
}

func (service *userService) GetList(meta *requests.MetaPaginationRequest) responses.Response {

	// users, count, err := service.userRepo.GetListUser(meta,)

	return service.common.StatusOk(nil, nil, "get list user successfully")
}

func (service *userService) DeleteUser(id string, tx *gorm.DB) responses.Response {

	return service.common.StatusOk(nil, nil, "delete user successfully")
}

func (service *userService) VerifyUser(request *requests.VerifyUser, tx *gorm.DB) responses.Response {

	authResponse := responses.AuthResponse{
		// User:         &userResponse,
		// AccessToken:  token,
		// RefreshToken: rToken,
	}

	return service.common.StatusCreated(authResponse, "verify user successfully")
}

func (service *userService) ChangePassword(ctx context.Context, password string, tx *gorm.DB) responses.Response {

	meta := ctx.Value("metadata").(models.TokenMetaData)

	whereClause := map[string]interface{}{
		"id": meta.Id,
	}

	var userModel *models.UserModel
	userModel, err := service.userRepo.GetDetailUser(whereClause, nil, nil, nil)
	if err != nil {
		log.Println("[userService][ChangePassword] error get detail user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	if userModel == nil {
		log.Println("[userService][ChangePassword] something wrong in your request")
		return service.common.StatusBadRequest(nil, "something wrong in your request")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[userService][RegisterUser] error bcrypt password :", err)
		return service.common.StatusServerError("something went wrong")
	}

	userModel.Password = string(hashedPass)

	_, err = service.userRepo.UpdateUser(userModel, tx)
	if err != nil {
		log.Println("[userService][ChangePassword] error update user :", err)
		return service.common.StatusServerError("something went wrong")
	}

	return service.common.StatusOk(nil, nil, "change password user successfully")
}
