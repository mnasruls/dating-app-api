package services

import (
	"context"
	"dating-app-api/configs"
	"dating-app-api/entities/requests"
	"dating-app-api/entities/responses"
	"dating-app-api/repositories"
	"dating-app-api/utils"
)

/*
  - when current user hit swipe,
    save user_id target to redis, so current user will not saw a same user in 1 day
  - when get list user, check redis first, to add all user_id that user already swipe
    to query NOT IN
*/
type SwipeServiceInterface interface {
	SwipService(ctx context.Context, req *requests.SwipeRequest) responses.Response
}

type swipeService struct {
  swipeRepo repo
	userRepo  repositories.UserRepositoryInterface
	common    responses.CommondResponse
	redisUtil *utils.Redis
	envs      *configs.EnviConfig
}
