package controller

import (
	"context"
	"encoding/json"
	"hackathon/usecase"
	"net/http"
	"time"
)

type SearchUserController struct {
	UseCase *usecase.UserSearchUseCase
}

func (c *SearchUserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := c.UseCase.Execute(ctx, name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
