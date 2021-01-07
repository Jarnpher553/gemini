package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"time"
)

var privateKey = "1234!@#$qwer"

func New(custom map[string]interface{}, duration time.Duration) (string, error) {
	claims := jwt.MapClaims(custom)

	t := time.Now()

	claims["aud"] = "Jarnpher553"
	claims["exp"] = t.Add(duration).Unix()
	claims["jti"] = uuid.NewV4().String()
	claims["iat"] = t.Unix()
	claims["iss"] = "Jarnpher553"
	claims["sub"] = "auth"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(privateKey + "_" + claims["jti"].(string)))

}

func Parse(tokenStr string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("unexpected claims type")
		}

		return []byte(privateKey + "_" + claims["jti"].(string)), nil
	})

	if err != nil {
		return nil, err
	}

	if t.Valid {
		return t.Claims.(jwt.MapClaims), nil
	} else {
		return nil, errors.New("token string invalid")
	}
}
