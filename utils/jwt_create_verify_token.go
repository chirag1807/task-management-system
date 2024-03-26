package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/chirag1807/task-management-system/config"
	"github.com/chirag1807/task-management-system/constant"
	"github.com/golang-jwt/jwt"
)

// CreateJWTToken uses golang-jwt package to generate jwt token and return that token as string.
// Here it takes token expire time and user id and add them to token as claims.
func CreateJWTToken(tokenExpiryTime time.Time, userId int64) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.FormatInt(954488202459119617, 10),
		"exp":    tokenExpiryTime.Unix(),
	})

	token, err := jwtToken.SignedString([]byte(config.JWtSecretKey.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyJWTToken takes token as parameter, verifies it and return userId and nil in case of token verified successfully and 0 and error if token verification failed.
func VerifyJWTToken(token string) (int64, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWtSecretKey.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !jwtToken.Valid {
		return 0, errors.New(constant.INVALID_TOKEN)
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New(constant.INVALID_CLAIMS)
	}
	userIdFromClaims, _ := claims["userId"].(string)
	userId, _ := strconv.ParseInt(userIdFromClaims, 10, 64)
	return userId, nil
}
