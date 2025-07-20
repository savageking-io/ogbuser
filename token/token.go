package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Config struct {
	Secret string `yaml:"secret"`
	Expiry int    `yaml:"expiry"`
	Issuer string `yaml:"issuer"`
}

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

var config Config

func SetConfig(inConfig *Config) {
	config = *inConfig
}

func Generate(userId int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.Expiry) * time.Minute)

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
