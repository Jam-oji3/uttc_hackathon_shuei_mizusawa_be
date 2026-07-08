package controller

import (
	"context"
	"encoding/json"
	"hackathon/model"
	"hackathon/usecase"
	"log"
	"net/http"
	"time"
)

type PostCreateController struct {
	AuthUC   *usecase.AuthUserUseCase
	CreateUC *usecase.PostCreateUseCase
}

func NewPostCreateController(authUC *usecase.AuthUserUseCase, createUC *usecase.PostCreateUseCase) *PostCreateController {
	return &PostCreateController{
		AuthUC:   authUC,
		CreateUC: createUC,
	}
}

func (c *PostCreateController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text      string  `json:"text"`
		ReplyTo   *string `json:"replyTo"`
		RepostRef *string `json:"repostRef"`
		MediaType *string `json:"mediaType"`
		MediaURL  *string `json:"mediaUrl"`
	}

	// Content-Typeの確認（オプションだが安全のため）
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	idToken, err := ExtractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	userId, _, _, err := c.AuthUC.Exec(ctx, idToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// JSONボディをデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	post, err := c.CreateUC.Execute(ctx, userId, req.Text, req.ReplyTo, req.RepostRef, req.MediaType, req.MediaURL)
	if err != nil {
		log.Printf("failed to create a post: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool        `json:"success"`
		Post    *model.Post `json:"post"`
		Message string      `json:"message"`
	}{
		Success: true,
		Post:    post,
		Message: "Post created successfully",
	})
}
