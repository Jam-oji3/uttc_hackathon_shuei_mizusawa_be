package controller

import (
	"context"
	"encoding/json"
	"errors"
	"hackathon/usecase"
	"net/http"
	"net/mail"
	"regexp"
	"time"
)

type UserRegisterController struct {
	UseCase *usecase.UserRegisterUseCase
}
type ReqBodyForHTTPPost struct {
	UserName    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	IconURL     string `json:"icon_url"`
	Email       string `json:"email"`
}

type ResBodyForHTTPPost struct {
	Id string `json:"id"`
}

func (c *UserRegisterController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req ReqBodyForHTTPPost

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := isValidInput(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := c.UseCase.Execute(ctx, req.UserName, req.DisplayName, req.Bio, req.IconURL, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResBodyForHTTPPost{Id: id})
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func isValidInput(req ReqBodyForHTTPPost) error {
	if req.UserName == "" {
		return errors.New("username is required")
	}
	if req.DisplayName == "" {
		return errors.New("display_name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if len(req.UserName) > 50 {
		return errors.New("username is too long")
	}
	if !usernameRegex.MatchString(req.UserName) {
		return errors.New("username can only contain letters, numbers, and underscore")
	}
	if len(req.DisplayName) > 50 {
		return errors.New("display_name is too long")
	}
	if len(req.Bio) > 200 {
		return errors.New("bio is too long")
	}
	if len(req.IconURL) > 2083 {
		return errors.New("icon_url is too long")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
