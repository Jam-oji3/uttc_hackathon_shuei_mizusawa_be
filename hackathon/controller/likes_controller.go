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

type LikeController struct {
	AuthUC   *usecase.AuthUserUseCase
	CreateUC *usecase.LikeCreateUseCase
	DeleteUC *usecase.LikeDeleteUseCase
}

func NewLikeController(authUC *usecase.AuthUserUseCase, createUC *usecase.LikeCreateUseCase, deleteUC *usecase.LikeDeleteUseCase) *LikeController {
	return &LikeController{
		AuthUC:   authUC,
		CreateUC: createUC,
		DeleteUC: deleteUC,
	}
}

func (c *LikeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == http.MethodPost {
		var req struct {
			PostId string `json:"postId"`
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		like, err := c.CreateUC.Execute(ctx, userId, req.PostId)
		if err != nil {
			log.Printf("like creation failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Success bool        `json:"success"`
			Like    *model.Like `json:"like"`
			Message string      `json:"message"`
		}{
			Success: true,
			Like:    like,
			Message: "Like creation successful",
		})
		return

	} else if r.Method == http.MethodDelete {

		postId := r.URL.Query().Get("postId")

		err := c.DeleteUC.Execute(ctx, userId, postId)
		if err != nil {
			log.Printf("like deletion failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
