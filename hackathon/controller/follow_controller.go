package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"hackathon/usecase"
	"log"
	"net/http"
	"time"
)

type FollowController struct {
	AuthUC   *usecase.AuthUserUseCase
	CreateUC *usecase.FollowCreateUseCase
	DeleteUC *usecase.FollowDeleteUseCase
}

func NewFollowController(authUC *usecase.AuthUserUseCase, createUC *usecase.FollowCreateUseCase, deleteUC *usecase.FollowDeleteUseCase) *FollowController {
	return &FollowController{
		AuthUC:   authUC,
		CreateUC: createUC,
		DeleteUC: deleteUC,
	}
}

func (c *FollowController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	followedId := vars["followed"]

	if followedId == "" {
		http.Error(w, "followed id is required", http.StatusBadRequest)
		return
	}

	idToken, err := ExtractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	followerId, _, _, err := c.AuthUC.Exec(ctx, idToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {

		err = c.CreateUC.Execute(ctx, followerId, followedId)
		if err != nil {
			log.Printf("failed to follow %s : %v", followedId, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: true,
			Message: fmt.Sprintf("Followed %s", followedId),
		})
	} else if r.Method == http.MethodDelete {
		err = c.DeleteUC.Execute(ctx, followerId, followedId)
		if err != nil {
			log.Printf("failed to unfollow %s : %v", followedId, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
