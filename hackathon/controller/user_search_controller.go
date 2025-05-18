package controller

import (
	"encoding/json"
	"hackathon/usecase"
	"net/http"
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

	users, err := c.UseCase.Execute(name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
