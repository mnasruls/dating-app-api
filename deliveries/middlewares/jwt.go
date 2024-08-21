package middlewares

import (
	"dating-app-api/configs"
	"dating-app-api/entities/models"
	"dating-app-api/entities/responses"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func UserVerify(conf *configs.EnviConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authToken := string(c.Request().Header.Peek("Authorization"))
		if authToken == "" {
			return AuthFailedHandler(c, "token is required")
		}

		splitToken := strings.Split(strings.TrimSpace(authToken), " ")
		if len(splitToken) < 2 {
			return AuthFailedHandler(c, "no token is headers")
		}

		token, err := verifyToken(splitToken[1], conf.JwtKey)
		if err != nil || !token.Valid {
			return AuthFailedHandler(c, "invalid token")
		}

		metadata := extractTokenMetadata(token)
		if metadata.Id == "" {
			return AuthFailedHandler(c, "invalid meta data")
		}

		var userTokenData models.TokenMetaData
		err = conf.Redis.RetrieveDataFromRedis(fmt.Sprintf("metaat:%v", metadata.Id), &userTokenData)
		if err != nil {
			return AuthFailedHandler(c, "invalid metadata or token expired")
		}

		c.Locals("metadata", userTokenData)
		return c.Next()
	}
}

func verifyToken(token, key string) (*jwt.Token, error) {
	tokenJWT, err := jwt.Parse(token, func(tokenJWT *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := tokenJWT.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tokenJWT.Header["alg"])
		}

		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	return tokenJWT, nil
}

func extractTokenMetadata(jwtToken *jwt.Token) (res models.TokenMetaData) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok {
		ID, _ := claims["id"].(string)
		verify, _ := claims["verify"].(bool)

		exp, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 0, 64)
		if err != nil {
			return res
		}

		res.Id = ID
		if claims["rt_id"] != nil {
			res.RtId = claims["rt_id"].(string)
		}

		res.Exp = exp
		res.Verify = verify
		return res
	}

	return res
}

func AuthFailedHandler(c *fiber.Ctx, message string) error {

	return c.Status(http.StatusUnauthorized).JSON(responses.Response{
		StatusCode: 401,
		Message:    message,
	})
}

func GenerateToken(conf *configs.EnviConfig, data models.TokenMetaData, isRefresh bool) (token string, err error) {

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	acClaims := jwtToken.Claims.(jwt.MapClaims)

	jwtKey := conf.JwtKey
	expTime := conf.JwtAtExpTime
	if isRefresh {
		jwtKey = conf.JwtRKey
		expTime = conf.JwtRtExpTime
	}

	jwtAcExpiredAt := time.Minute * time.Duration(expTime)
	jwtRtExpiredAt := time.Minute * time.Duration(expTime)
	// Set claims
	acClaims["id"] = data.Id
	acClaims["exp"] = time.Now().Add(jwtAcExpiredAt).Unix()
	acClaims["verify"] = data.Verify
	if isRefresh {
		acClaims["rt_id"] = data.RtId
		acClaims["exp"] = time.Now().Add(jwtRtExpiredAt).Unix()
		acClaims["verify"] = data.Verify
	}

	token, err = jwtToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	if isRefresh {
		err = conf.Redis.SaveDataToRedis(fmt.Sprintf("metart:%v:%v", data.Id, acClaims["rt_id"]), data, jwtRtExpiredAt)
		if err != nil {
			return "", err
		}
	} else {
		err = conf.Redis.SaveDataToRedis(fmt.Sprintf("metaat:%v", data.Id), data, jwtAcExpiredAt)
		if err != nil {
			return "", err
		}
	}

	return token, nil
}

func RefreshTokenVerify(conf *configs.EnviConfig) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authToken := string(c.Request().Header.Peek("Authorization"))
		if authToken == "" {
			return AuthFailedHandler(c, "token is required")
		}

		splitToken := strings.Split(strings.TrimSpace(authToken), " ")
		if len(splitToken) < 2 {
			return AuthFailedHandler(c, "no token is headers")
		}

		token, err := verifyToken(splitToken[1], conf.JwtRKey)
		if err != nil || !token.Valid {
			return AuthFailedHandler(c, "invalid token")
		}

		metaData := extractTokenMetadata(token)
		if metaData.Id == "" {
			return AuthFailedHandler(c, "invalid meta data")
		}

		var newTokenMetaData models.TokenMetaData
		redisKey := fmt.Sprintf("metart:%v:%v", metaData.Id, metaData.RtId)
		if err := conf.Redis.RetrieveDataFromRedis(redisKey, &newTokenMetaData); err != nil {
			return AuthFailedHandler(c, "invalid metadata or token expired")
		}

		newTokenMetaData.RtId = metaData.RtId
		c.Locals("metadata", newTokenMetaData)
		return c.Next()
	}
}
