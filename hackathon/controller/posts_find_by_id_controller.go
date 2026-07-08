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

type PostFindByIdController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.PostFindByIdUseCase
}

func NewPostFindByIdController(authUC *usecase.AuthUserUseCase, useCase *usecase.PostFindByIdUseCase) *PostFindByIdController {
	return &PostFindByIdController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *PostFindByIdController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["postId"]

	if postId == "" {
		http.Error(w, "Missing postId", http.StatusBadRequest)
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

	post, err := c.UseCase.Execute(ctx, userId, postId)
	if err != nil {
		log.Printf("Error while executing post id %s: %v", postId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool                         `json:"success"`
		Post    *model.PostWithUserAndCounts `json:"post"`
		Message string                       `json:"message"`
	}{
		Success: true,
		Post:    post,
		Message: "post fetched successfully",
	})
}
