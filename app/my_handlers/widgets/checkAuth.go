package widgets

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-chat-app/config"
)

func CheckAuth(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, nil, fmt.Errorf("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SECRET_KEY), nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token")
	}

	return token, claims, nil
}
