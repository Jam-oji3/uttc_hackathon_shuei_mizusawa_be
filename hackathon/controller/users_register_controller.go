package controller

import (
	"context"
	"encoding/json"
	"errors"
	"hackathon/model"
	"hackathon/usecase"
	"log"
	"net/http"
	"time"
)

type UserRegisterController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.UserRegisterUseCase
}

func NewUserRegisterController(authUC *usecase.AuthUserUseCase, useCase *usecase.UserRegisterUseCase) *UserRegisterController {
	return &UserRegisterController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *UserRegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idToken, err := ExtractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userId, _, _, err := c.AuthUC.Exec(ctx, idToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// JSONデコード用の構造体を定義
	var req struct {
		UserName    string `json:"username"`
		DisplayName string `json:"displayName"`
		Bio         string `json:"bio"`
		Email       string `json:"email"`
		IconURL     string `json:"iconUrl"`
	}

	// Content-Typeの確認（オプションだが安全のため）
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// JSONボディをデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := c.UseCase.Execute(ctx, userId, req.UserName, req.DisplayName, req.Bio, req.IconURL, req.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			log.Printf("user already exists: %v (id=%s, userName=%s)", err, userId, req.UserName)
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			log.Printf("failed to register user: %v (id=%s, userName=%s)", err, userId, req.UserName)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 成功レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool        `json:"success"`
		User    *model.User `json:"user"`
		Message string      `json:"message"`
	}{
		Success: true,
		User:    user,
		Message: "User registration successful",
	})
}
