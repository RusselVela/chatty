package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
	"time"
)

const jwtSecretPath = "jwt.secret"

var jwtSecret []byte

func SetupJWTSecret(k *koanf.Koanf) {
	jwtSecret = []byte(k.String(jwtSecretPath))
}

func generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour)
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	zap.L().Info(fmt.Sprintf("token generated for user: %s", username))
	return tokenString, nil
}
