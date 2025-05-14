package controller

import (
	"db/usecase"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReqBodyForHTTPPost struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ResBodyForHTTPPost struct {
	Id string `json:"id"`
}

type RegisterUserController struct {
	UseCase *usecase.RegisterUserUseCase
}

func (c *RegisterUserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var req ReqBodyForHTTPPost
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || len(req.Name) > 50 || req.Age < 20 || req.Age > 80 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	id, err := c.UseCase.Execute(req.Name, req.Age)
	if err != nil {
		fmt.Printf("fail: RegisterUserController, %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ResBodyForHTTPPost{Id: id})
}
