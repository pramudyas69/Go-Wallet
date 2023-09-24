package util

import (
	"e-wallet/domain"
	"e-wallet/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtToken struct {
	cnf *config.Config
}

func NewJwt(cnf *config.Config) domain.JwtInterface {
	return &jwtToken{
		cnf: cnf,
	}
}

func (j jwtToken) GenerateToken(userID int64, email string, expirationMinutes int64) (string, error) {
	exp := time.Now().Add(time.Minute * time.Duration(expirationMinutes))

	claims := &domain.CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.cnf.Jwt.AccessTokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j jwtToken) VerifyToken(tokenString string) (*domain.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.cnf.Jwt.AccessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
