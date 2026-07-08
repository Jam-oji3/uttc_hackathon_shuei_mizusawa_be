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

type PostGetRepliesController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.PostGetRepliesUseCase
}

func NewPostGetRepliesController(authUC *usecase.AuthUserUseCase, useCase *usecase.PostGetRepliesUseCase) *PostGetRepliesController {
	return &PostGetRepliesController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *PostGetRepliesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // デフォルト値 or エラーにしても良い
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := c.UseCase.Execute(ctx, userId, postId, limit, offset)
	if err != nil {
		log.Printf("failed to fetch replies (postid=%v): %v", postId, err)
		http.Error(w, "failed to fetch replies", http.StatusInternalServerError)
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
