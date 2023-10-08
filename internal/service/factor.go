package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"golang.org/x/crypto/bcrypt"
)

type factorService struct {
	factorRepository domain.FactorRepository
}

func NewFactor(factorRepository domain.FactorRepository) domain.FactorService {
	return &factorService{
		factorRepository: factorRepository,
	}
}

func (f factorService) ValidatePin(ctx context.Context, req dto.ValidatePinReq) error {
	factor, err := f.factorRepository.FindByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if factor == (domain.Factor{}) {
		return domain.ErrPinInvalid
	}

	err = bcrypt.CompareHashAndPassword([]byte(factor.Pin), []byte(req.Pin))
	if err != nil {
		return domain.ErrPinInvalid
	}

	return nil
}

func (f factorService) Insert(ctx context.Context, req dto.CreatePinReq) error {
	res, err := f.factorRepository.FindByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if res != (domain.Factor{}) {
		return domain.ErrPinExist
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	factor := domain.Factor{
		UserID: req.UserID,
		Pin:    string(hashedPin),
	}

	err = f.factorRepository.Insert(ctx, &factor)
	if err != nil {
		return err
	}

	return nil
}
