package token

import (
	"errors"
	"file-server/db/redis"
	"file-server/setting"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// TokenExpiredTime 的单位是以秒计算的
	AccessTokenExpiredTime = time.Minute * 5
	FreshTokenExpiredTime  = time.Minute * 2
	TokenExpiredError      = "token is expired"
	TokenIllegalError      = "token is illegal"
	TokenHandleError       = "can not handle this token, maybe it is not a valid token"
)

type TokenClainms struct {
	Data interface{} `json:"data"`
	jwt.RegisteredClaims
}

func AccessCreate(data interface{}) (string, error) {
	return Create(data, AccessTokenExpiredTime)
}

func FreshCreat(data interface{}) (string, error) {
	token, err := Create(data, FreshTokenExpiredTime)
	if err != nil {
		return "", err
	}
	//将refreshToke放入redis
	err = PushToken(token, data.(string))
	if err != nil {
		return "", err
	}
	return token, nil
}

func Create(data interface{}, expiredTime time.Duration) (string, error) {
	claim := TokenClainms{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredTime)),
			Issuer:    setting.Conf.TokenConfig.Issuer,
		},
	}

	originToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return originToken.SignedString([]byte(setting.Conf.TokenConfig.Secret))
}

func Parse(token string) (*TokenClainms, error) {
	tokenStruct, err := jwt.ParseWithClaims(token, &TokenClainms{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.Conf.TokenConfig.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenStruct.Claims.(*TokenClainms); ok && tokenStruct.Valid {
		return claims, nil
	}
	return nil, errors.New(TokenHandleError)
}

func PushToken(tokenString, userId string) error {
	return redis.DB.Set(redis.Ctx, tokenString, userId, FreshTokenExpiredTime).Err()
}
