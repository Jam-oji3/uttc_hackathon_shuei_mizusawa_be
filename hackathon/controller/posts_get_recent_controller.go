package controller

import (
	"context"
	"encoding/json"
	"hackathon/model"
	"hackathon/usecase"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PostGetRecentController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.PostGetRecentUseCase
}

func NewPostGetRecentController(authUC *usecase.AuthUserUseCase, useCase *usecase.PostGetRecentUseCase) *PostGetRecentController {
	return &PostGetRecentController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *PostGetRecentController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	posts, err := c.UseCase.Execute(ctx, userId, limit, offset)
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
