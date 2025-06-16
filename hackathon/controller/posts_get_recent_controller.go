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
	UseCase *usecase.PostGetRecentUseCase
}

func NewPostGetRecentController(useCase *usecase.PostGetRecentUseCase) *PostGetRecentController {
	return &PostGetRecentController{
		UseCase: useCase,
	}
}

func (c *PostGetRecentController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータから取得
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	posts, err := c.UseCase.Execute(ctx, limit, offset)
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
