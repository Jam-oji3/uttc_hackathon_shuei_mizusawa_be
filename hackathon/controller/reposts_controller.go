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

type RepostController struct {
	CreateUC *usecase.RepostCreateUseCase
	DeleteUC *usecase.RepostDeleteUseCase
}

func NewRepostController(createUC *usecase.RepostCreateUseCase, deleteUC *usecase.RepostDeleteUseCase) *RepostController {
	return &RepostController{
		CreateUC: createUC,
		DeleteUC: deleteUC,
	}
}

func (c *RepostController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			UserId string `json:"userId"`
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

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		repost, err := c.CreateUC.Execute(ctx, req.UserId, req.PostId)
		if err != nil {
			log.Printf("repost creation failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Success bool          `json:"success"`
			Repost  *model.Repost `json:"repost"`
			Message string        `json:"message"`
		}{
			Success: true,
			Repost:  repost,
			Message: "Repost creation successful",
		})

	} else if r.Method == http.MethodDelete {

		userId := r.URL.Query().Get("userId")
		postId := r.URL.Query().Get("postId")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := c.DeleteUC.Execute(ctx, userId, postId)
		if err != nil {
			log.Printf("repost deletion failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
