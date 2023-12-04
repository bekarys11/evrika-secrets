package common

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetUserIdFromToken(r *http.Request) (string, error) {
	claims, err := GetTokenClaims(r)
	if err != nil {
		log.Printf("error getting token claims: %v", err)
		return "", err
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("user id is not provided")
	}

	return strconv.Itoa(int(id)), nil
}

func GetRoleFromToken(r *http.Request) (role string, err error) {
	claims, err := GetTokenClaims(r)
	if err != nil {
		log.Printf("error getting token claims: %v", err)
		return "", err
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("user id is not provided")
	}

	return role, nil
}

func GetTokenClaims(r *http.Request) (jwt.MapClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	_, tokenStr, _ := strings.Cut(r.Header.Get("Authorization"), "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}
