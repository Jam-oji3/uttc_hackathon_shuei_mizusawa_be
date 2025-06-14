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
		writeUnifiedResponse(w, false, "method not allowed", "", "", nil, http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeUnifiedResponse(w, false, "no authorization header", "", "", nil, http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		writeUnifiedResponse(w, false, "authorization header must be in format 'Bearer <token>'", "", "", nil, http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	uid, email, user, err := c.UseCase.Exec(ctx, parts[1])
	if errors.Is(err, model.ErrUserNotFound) {
		writeUnifiedResponse(w, false, "user not found", uid, email, nil, http.StatusNotFound)
		return
	}
	if err != nil {
		writeUnifiedResponse(w, false, err.Error(), uid, email, nil, http.StatusInternalServerError)
		return
	}

	writeUnifiedResponse(w, true, "user authenticated", uid, email, user, http.StatusOK)
}

func writeUnifiedResponse(w http.ResponseWriter, success bool, message string, uid string, email string, user interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": success,
		"message": message,
		"uid":     uid,
		"email":   email,
	}
	if user != nil {
		response["user"] = user
	}

	json.NewEncoder(w).Encode(response)
}
