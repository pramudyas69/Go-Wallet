package domain

import (
	"math/rand"
	"time"
)

type UtilInterface interface {
	GetTokenGenerator(length int) string
	GenerateRandomNumber(length int) string
}

type utilInterface struct {
}

func NewUtil() UtilInterface {
	return &utilInterface{}
}

func (u utilInterface) GetTokenGenerator(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func (u utilInterface) GenerateRandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
