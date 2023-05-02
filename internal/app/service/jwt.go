package service

import (
	"fmt"
	"time"

	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/golang-jwt/jwt/v4"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

const jwtSecretPath = "jwt.secret"

var JWTSecret []byte

type JWTCustomClaims struct {
	Id         string `json:"userId"`
	Username   string `json:"username"`
	Authorized bool   `json:"authorized"`
	jwt.RegisteredClaims
}

// SetupJWTSecret sets the value of the secret from configuration
func SetupJWTSecret(k *koanf.Koanf) {
	JWTSecret = []byte(k.String(jwtSecretPath))
}

// generateJWT generates a new JWT with the given duration in minutes
func generateJWT(user inmemory.UserBean, duration int) (string, error) {
	if duration <= 0 {
		duration = 60 * 24
	}
	claims := &JWTCustomClaims{
		Id:         user.Id.String(),
		Username:   user.Username,
		Authorized: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(duration))),
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

// parseJWT takes tokenStr and parses it into a jwt.Token struct
func parseJWT(tokenStr string) (*jwt.Token, error) {
	claims := &JWTCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid websocket auth token")
		}
		return JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid websocket auth token")
	}

	return token, nil
}
