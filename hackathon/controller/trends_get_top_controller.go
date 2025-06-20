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

type TrendGetTopController struct {
	UseCase *usecase.TrendGetTopUseCase
}

func NewTrendGetTopController(useCase *usecase.TrendGetTopUseCase) *TrendGetTopController {
	return &TrendGetTopController{
		UseCase: useCase,
	}
}

func (c *TrendGetTopController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sinceHoursStr := r.URL.Query().Get("since")
	limitStr := r.URL.Query().Get("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	sinceHours, err := strconv.Atoi(sinceHoursStr)
	if err != nil || sinceHours <= 0 {
		sinceHours = 24
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	trends, err := c.UseCase.Execute(ctx, sinceHours, limit)
	if err != nil {
		log.Printf("failed to fetch trends: %v", err)
		http.Error(w, "failed to fetch trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success bool                 `json:"success"`
		Trends  []model.TrendSummary `json:"trends"`
		Message string               `json:"message"`
	}{
		Success: true,
		Trends:  trends,
		Message: "trends successfully fetched",
	})
}
