package usecase

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
)

type AuthUserUseCase struct {
	UserRepo         repository.UserRepository
	FirebaseAuthRepo repository.FirebaseAuthRepository
	DB               *sql.DB
}

func NewAuthUserUseCase(firebaseAuthRepo repository.FirebaseAuthRepository, userRepo repository.UserRepository, db *sql.DB) *AuthUserUseCase {
	return &AuthUserUseCase{FirebaseAuthRepo: firebaseAuthRepo, UserRepo: userRepo, DB: db}
}

func (uc *AuthUserUseCase) Exec(ctx context.Context, idToken string) (string, string, *model.User, error) {
	token, err := uc.FirebaseAuthRepo.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", "", nil, err
	}

	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		return "", "", nil, errors.New("email claim is invalid or empty")
	}

	user, err := uc.UserRepo.FindById(ctx, uc.DB, token.UID)
	if err != nil {
		return token.UID, email, nil, err
	}

	return token.UID, email, user, nil
}
