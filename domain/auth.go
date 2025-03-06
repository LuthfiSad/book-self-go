package domain

import (
	"context"
	"go-rest-api/dto"
)

type AuthService interface {
	Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error)
	Validate(ctx context.Context, tokenString string) (dto.UserData, error)
}
