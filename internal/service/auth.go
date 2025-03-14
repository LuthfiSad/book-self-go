package service

import (
	"context"
	"go-rest-api/domain"
	"go-rest-api/dto"
	"go-rest-api/internal/config"
	"go-rest-api/internal/constants"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (a authService) Register(ctx context.Context, req dto.RegisterReq) (dto.UserData, error) {
	existingUser, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != uuid.Nil {
		return dto.UserData{}, constants.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.UserData{}, err
	}

	newUser := domain.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = a.userRepository.Save(ctx, &newUser)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.UserData{}, err
	}

	return dto.UserData{
		Id:    newUser.ID.String(),
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}, nil
}

func (a authService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return dto.AuthRes{}, err
	}
	if user.ID == uuid.Nil {
		return dto.AuthRes{}, constants.ErrInvalidCredential
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.AuthRes{}, constants.ErrInvalidCredential
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
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
				Role:  claims["role"].(string),
			}, nil
		}
	}
	return dto.UserData{}, constants.ErrInvalidCredential
}
