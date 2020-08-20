package jwt

import (
	"errors"
	"fmt"
	guid "github.com/Jarnpher553/gemini/uuid"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"time"
)

var privateKey = "perfect"

type CustomClaims struct {
	UserIdUUID guid.GUID
	UserIdInt  int
	jwt.StandardClaims
}

func New(id interface{}) (string, error) {
	var claims CustomClaims

	t := time.Now()

	switch v := id.(type) {
	case guid.GUID:
		claims = CustomClaims{
			UserIdUUID: v,
		}
	case int:
		claims = CustomClaims{
			UserIdInt: v,
		}
	}
	claims.StandardClaims = jwt.StandardClaims{
		Audience:  "Jarnpher553",
		ExpiresAt: t.Add(time.Hour * 24 * 3).Unix(),
		Id:        uuid.NewV4().String(),
		IssuedAt:  t.Unix(),
		Issuer:    "Jarnpher553",
		Subject:   "auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	return token.SignedString([]byte(privateKey))

}

func Parse(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(privateKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, errors.New("unexpected claims type or token string invalid")
	}
}

func (c *CustomClaims) Valid() error {
	if c.UserIdInt == 0 && c.UserIdUUID == "" {
		return errors.New("invalid user")
	}

	return c.StandardClaims.Valid()
}
