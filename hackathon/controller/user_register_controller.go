package controller

import (
	"context"
	"encoding/json"
	"errors"
	"hackathon/model"
	"hackathon/usecase"
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	id := r.PostFormValue("id")
	userName := r.FormValue("userName")
	displayName := r.FormValue("displayName")
	bio := r.FormValue("bio")
	email := r.FormValue("email")

	// ファイルは受け取るだけで何もしない
	file, _, err := r.FormFile("iconFile")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "failed to get icon file", http.StatusBadRequest)
		return
	}
	if file != nil {
		defer file.Close()
	}

	// 仮のサンプルアイコンURLを使用
	iconURL := "https://example.com/sample-icon.png"

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := c.UseCase.Execute(ctx, id, userName, displayName, bio, iconURL, email)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409
		} else {
			// その他の予期せぬエラー
			http.Error(w, "internal server error", http.StatusInternalServerError) // 500
		}
		return
	}

	// JSONレスポンスを構築して返す
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
