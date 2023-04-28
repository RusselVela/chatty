package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/golang-jwt/jwt/v4"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
	"time"
)

const jwtSecretPath = "jwt.secret"

var JWTSecret []byte

type JWTCustomClaims struct {
	Id         string `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.RegisteredClaims
}

func SetupJWTSecret(k *koanf.Koanf) {
	JWTSecret = []byte(k.String(jwtSecretPath))
}

func generateJWT(user inmemory.UserBean) (string, error) {
	claims := &JWTCustomClaims{
		Id:         user.Id.String(),
		Username:   user.Username,
		Authorized: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	zap.L().Info(fmt.Sprintf("token generated for user: %s", user.Id.String()))
	return tokenString, nil
}
