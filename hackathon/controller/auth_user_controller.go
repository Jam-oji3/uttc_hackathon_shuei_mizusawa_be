package controller

import (
	"context"
	"encoding/json"
	"errors"
	"hackathon/model"
	"hackathon/usecase"
	"net/http"
	"strings"
	"time"
)

type AuthUserController struct {
	UseCase *usecase.AuthUserUseCase
}

func NewAuthUserController(useCase *usecase.AuthUserUseCase) *AuthUserController {
	return &AuthUserController{UseCase: useCase}
}

func (c *AuthUserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "no authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		http.Error(w, "authorization header must be in format 'Bearer <token>'", http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := c.UseCase.Exec(ctx, parts[1])
	if errors.Is(err, model.ErrUserNotFound) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
