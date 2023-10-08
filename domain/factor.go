package domain

import (
	"context"
	"e-wallet/dto"
)

type Factor struct {
	ID     int64  `db:"id"`
	UserID int64  `db:"user_id"`
	Pin    string `db:"pin"`
}

type FactorRepository interface {
	FindByUserID(ctx context.Context, userID int64) (Factor, error)
	Insert(ctx context.Context, factor *Factor) error
}

type FactorService interface {
	ValidatePin(ctx context.Context, req dto.ValidatePinReq) error
	Insert(ctx context.Context, req dto.CreatePinReq) error
}
