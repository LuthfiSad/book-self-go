package dto

type AuthReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
