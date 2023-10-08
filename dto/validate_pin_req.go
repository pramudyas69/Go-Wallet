package dto

type ValidatePinReq struct {
	Pin    string `json:"pin" validate:"required"`
	UserID int64  `json:"-"`
}
