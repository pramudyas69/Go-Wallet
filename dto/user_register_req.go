package dto

type UserRegisterReq struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
