package controller

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"hackathon/model"
	"hackathon/usecase"
	"log"
	"net/http"
	"time"
)

type UserFindProfileController struct {
	UseCase *usecase.UserFindProfileUseCase
}

func NewUserFindProfileController(useCase *usecase.UserFindProfileUseCase) *UserFindProfileController {
	return &UserFindProfileController{
		UseCase: useCase,
	}
}

func (c *UserFindProfileController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prof, err := c.UseCase.Execute(ctx, username)
	if err != nil {
		log.Printf("failed to fetch user profile (username: %v): %v", username, err)
		http.Error(w, "failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool              `json:"success"`
		Profile model.UserProfile `json:"profile"`
		Message string            `json:"message"`
	}{
		Success: true,
		Profile: *prof,
		Message: "user profile fetched successfully",
	})

}
