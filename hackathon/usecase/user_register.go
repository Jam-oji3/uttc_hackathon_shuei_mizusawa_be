package usecase

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type UserRegisterUseCase struct {
	TxExecutor repository.TransactionExecutor
	UserRepo   repository.UserRepository
	DB         *sql.DB
}

func NewUserRegisterUseCase(txExecutor repository.TransactionExecutor, userRepo repository.UserRepository, db *sql.DB) *UserRegisterUseCase {
	return &UserRegisterUseCase{TxExecutor: txExecutor, UserRepo: userRepo, DB: db}
}

func (uc *UserRegisterUseCase) Execute(ctx context.Context, id string, userName string, displayName string, bio string, iconURL string, email string) (*model.User, error) {
	now := time.Now()
	_, err := uc.UserRepo.FindByUserName(ctx, uc.DB, userName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	//すでにユーザーが存在する場合
	if err == nil {
		return nil, model.ErrUserAlreadyExists
	}

	user := model.User{
		Id:          id,
		UserName:    userName,
		DisplayName: displayName,
		Email:       email,
		Bio:         bio,
		IconURL:     iconURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := validateUserData(&user); err != nil {
		return nil, err
	}

	_, txErr := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.UserRepo.Insert(ctx, tx, &user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if txErr != nil {
		return nil, txErr
	}
	return &user, nil
}

func validateUserData(user *model.User) error {
	if user.Id == "" {
		return errors.New("uid is required")
	}

	// 2. Email形式チェック（シンプルな正規表現）
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	// 3. UserName: 空文字チェック、長さチェック（例: 3～30文字）、半角英数字と_のみ許可
	if len(user.UserName) < 5 || len(user.UserName) > 15 {
		return errors.New("userName must be between 5 and 15 characters")
	}
	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !userNameRegex.MatchString(user.UserName) {
		return errors.New("userName can only contain letters, numbers, and underscore")
	}

	// 4. DisplayName: 空文字禁止、最大長さチェック（例: 50文字以内）
	if strings.TrimSpace(user.DisplayName) == "" {
		return errors.New("displayName cannot be empty")
	}
	if len(user.DisplayName) > 50 {
		return errors.New("displayName is too long (max 50 characters)")
	}

	// 5. Bio: 任意だが長すぎないように（例: 200文字以内）
	if len(user.Bio) > 200 {
		return errors.New("bio is too long (max 200 characters)")
	}

	// 6. IconURL: 空文字は許可、空でなければURLとして妥当かチェック
	if user.IconURL != "" {
		parsedURL, err := url.ParseRequestURI(user.IconURL)
		if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
			return errors.New("iconURL must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}
