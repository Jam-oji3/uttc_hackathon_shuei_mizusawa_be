package controller

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"hackathon/model"
	"hackathon/usecase"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PostGetByUserController struct {
	UseCase *usecase.PostGetByUserUseCase
}

func NewPostGetByUserController(UseCase *usecase.PostGetByUserUseCase) *PostGetByUserController {
	return &PostGetByUserController{
		UseCase: UseCase,
	}
}

func (c *PostGetByUserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetId := vars["target"]

	if targetId == "" {
		http.Error(w, "target required", http.StatusBadRequest)
		return
	}

	viewerId := r.URL.Query().Get("viewer")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if viewerId == "" {
		http.Error(w, "viewer required", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // デフォルト値 or エラーにしても良い
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	posts, err := c.UseCase.Execute(ctx, targetId, viewerId, limit, offset)
	if err != nil {
		log.Printf("failed to fetch posts: %v", err)
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool                          `json:"success"`
		Posts   []model.PostWithUserAndCounts `json:"posts"`
		Message string                        `json:"message"`
	}{
		Success: true,
		Posts:   *posts,
		Message: "posts fetched successfully",
	})
}
