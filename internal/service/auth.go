package service

import (
	"context"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cnf            *config.Config
	userRepository domain.UserRepository
}

func NewAuth(cnf *config.Config,
	userRepository domain.UserRepository) domain.AuthService {
	return &authService{
		cnf:            cnf,
		userRepository: userRepository,
	}
}

func (a authService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.AuthRes{}, err
	}
	if user.Id == "" {
		return dto.AuthRes{}, domain.ErrInvalidCredential
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.AuthRes{}, domain.ErrInvalidCredential
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(a.cnf.Secret.Jwt))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.AuthRes{}, err
	}

	return dto.AuthRes{
		AccessToken: tokenString,
	}, nil
}

func (a authService) Validate(ctx context.Context, tokenString string) (dto.UserData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cnf.Secret.Jwt), nil
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.UserData{}, err
	}
	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return dto.UserData{
				Id:    claims["id"].(string),
				Name:  claims["name"].(string),
				Email: claims["email"].(string),
			}, nil
		}
	}
	return dto.UserData{}, domain.ErrInvalidCredential
}
