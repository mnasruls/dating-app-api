package configs

import (
	"dating-app-api/utils"
	"errors"
	"os"
	"strconv"
)

type EnviConfig struct {
	AppVersion   string
	AppHost      string
	AppPort      string
	AppEnv       string
	DbHost       string
	DbPort       string
	DbUsername   string
	DbPassword   string
	DbName       string
	JwtKey       string
	JwtRKey      string
	JwtAtExpTime int
	JwtRtExpTime int
	Redis        *utils.Redis
	ApiKey       string
}

func InitEnv() (EnviConfig, []error) {
	var env EnviConfig
	var errs []error

	env.AppVersion = os.Getenv("APP_VERSION")
	if env.AppVersion == "" {
		errs = append(errs, errors.New("app version env not found"))
	}

	env.AppPort = os.Getenv("APP_HOST")
	if env.AppPort == "" {
		errs = append(errs, errors.New("app version env not found"))
	}

	env.AppPort = os.Getenv("APP_PORT")
	if env.AppPort == "" {
		errs = append(errs, errors.New("app port env not found"))
	}

	env.AppEnv = os.Getenv("APP_ENV")
	if env.AppEnv == "" {
		errs = append(errs, errors.New("app env not found"))
	}

	env.DbHost = os.Getenv("DB_HOST")
	if env.DbHost == "" {
		errs = append(errs, errors.New("db host env not found"))
	}

	env.DbPort = os.Getenv("DB_PORT")
	if env.DbPort == "" {
		errs = append(errs, errors.New("db port env not found"))
	}

	env.DbUsername = os.Getenv("DB_USERNAME")
	if env.DbUsername == "" {
		errs = append(errs, errors.New("db username env not found"))
	}

	env.DbPassword = os.Getenv("DB_PASSWORD")
	if env.DbPassword == "" {
		errs = append(errs, errors.New("db password env not found"))
	}

	env.DbName = os.Getenv("DB_NAME")
	if env.DbName == "" {
		errs = append(errs, errors.New("db name env not found"))
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		errs = append(errs, errors.New("redis port env not found or invalid"))
	}

	confRedis := utils.ConfRedis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redisPort,
		User:     os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	if confRedis.Host == "" {
		errs = append(errs, errors.New("redis host env not found"))
	}

	if env.AppEnv != "local" {
		if confRedis.User == "" {
			errs = append(errs, errors.New("redis user env not found"))
		}

		if confRedis.Password == "" {
			errs = append(errs, errors.New("redis password env not found"))
		}
	}

	env.JwtKey = os.Getenv("JWT_KEY")
	if env.JwtKey == "" {
		errs = append(errs, errors.New("jwt key env not found"))
	}

	env.JwtRKey = os.Getenv("REFRESH_KEY")
	if env.JwtKey == "" {
		errs = append(errs, errors.New("jwt refresh key env not found"))
	}

	env.JwtAtExpTime, err = strconv.Atoi(os.Getenv("JWT_AT_EXP"))
	if err != nil {
		errs = append(errs, errors.New("jwt expired env not found or invalid"))
	}

	env.JwtRtExpTime, err = strconv.Atoi(os.Getenv("JWT_RT_EXP"))
	if err != nil {
		errs = append(errs, errors.New("jwt refresh expired env not found or invalid"))
	}

	env.ApiKey = os.Getenv("API_KEY")
	if env.ApiKey == "" {
		errs = append(errs, errors.New("api key env not found"))
	}

	if len(errs) > 0 {
		return env, errs
	} else {
		redisClient, err := utils.NewRedis(confRedis)
		if err != nil {
			errs = append(errs, err)
			return env, errs
		}
		env.Redis = redisClient
	}

	return env, nil
}
