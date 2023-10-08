package dto

type CreatePinReq struct {
	PIN    string `json:"pin"`
	UserID int64  `json:"-"`
}
