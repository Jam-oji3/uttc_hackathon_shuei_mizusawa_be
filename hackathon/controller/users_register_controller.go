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
	UseCase *usecase.UserRegisterUseCase
}

func NewUserRegisterController(useCase *usecase.UserRegisterUseCase) *UserRegisterController {
	return &UserRegisterController{UseCase: useCase}
}

func (c *UserRegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// JSONデコード用の構造体を定義
	var req struct {
		ID          string `json:"id"`
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := c.UseCase.Execute(ctx, req.ID, req.UserName, req.DisplayName, req.Bio, req.IconURL, req.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			log.Printf("user already exists: %v (id=%s, userName=%s)", err, req.ID, req.UserName)
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			log.Printf("failed to register user: %v (id=%s, userName=%s)", err, req.ID, req.UserName)
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
