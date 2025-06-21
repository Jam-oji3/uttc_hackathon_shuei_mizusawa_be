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

type PostSearchController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.PostSearchUseCase
}

func NewPostSearchController(authUC *usecase.AuthUserUseCase, useCase *usecase.PostSearchUseCase) *PostSearchController {
	return &PostSearchController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *PostSearchController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idToken, err := ExtractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	viewerUserId, _, _, err := c.AuthUC.Exec(ctx, idToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "Missing keyword parameter", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := c.UseCase.SearchPostsByKeyword(ctx, viewerUserId, keyword, limit, offset)
	if err != nil {
		log.Printf("failed to search posts: %v", err)
		http.Error(w, "failed to search posts", http.StatusInternalServerError)
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
