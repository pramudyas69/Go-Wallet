package domain

import "github.com/golang-jwt/jwt/v5"

// CustomClaims adalah struktur klaim khusus yang akan digunakan dalam token
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JwtInterface interface {
	GenerateToken(userID int64, email string, expirationMinutes int64) (string, error)
	VerifyToken(tokenString string) (*CustomClaims, error)
}
