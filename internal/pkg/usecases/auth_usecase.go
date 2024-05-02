package usecases

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/repositories"
	"claim_api_demo/pkg/auth"
	"claim_api_demo/pkg/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthUsecaseInterface interface {
	Login(cxt context.Context, req dtos.LoginRequest) (string, error)
}

type AuthUsecase struct {
	repo repositories.UserRepositoryInterface
}

func (u *AuthUsecase) Login(ctx context.Context, req dtos.LoginRequest) (string, error) {
	user, err := u.repo.TakeByConditions(ctx, bson.D{{"email", req.Email}})
	if err != nil {
		return "", err
	}
	if utils.EncrtyptPasswords(req.Password) != user.Password {
		return "", fmt.Errorf("INVALID_PASSWORD")
	}

	tokenString, err := auth.GenerateJWT(user.ID)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewAuthUsecase
func NewAuthUsecase(userRepo repositories.UserRepositoryInterface) AuthUsecaseInterface {
	return &AuthUsecase{
		repo: userRepo,
	}
}
