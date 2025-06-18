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
	UseCase *usecase.PostFindByIdUseCase
}

func NewPostFindByIdController(useCase *usecase.PostFindByIdUseCase) *PostFindByIdController {
	return &PostFindByIdController{
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

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
