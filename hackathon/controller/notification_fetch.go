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

type NotificationGetController struct {
	AuthUC  *usecase.AuthUserUseCase
	UseCase *usecase.NotificationFetchUseCase
}

func NewNotificationGetController(authUC *usecase.AuthUserUseCase, useCase *usecase.NotificationFetchUseCase) *NotificationGetController {
	return &NotificationGetController{
		AuthUC:  authUC,
		UseCase: useCase,
	}
}

func (c *NotificationGetController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータから userId と limit を取得
	limitStr := r.URL.Query().Get("limit")

	idToken, err := ExtractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userId, _, _, err := c.AuthUC.Exec(ctx, idToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userId == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20 // デフォルト値
	}

	notifications, err := c.UseCase.Execute(ctx, userId, limit)
	if err != nil {
		log.Printf("failed to fetch notifications: %v", err)
		http.Error(w, "failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Success       bool                  `json:"success"`
		Notifications []*model.Notification `json:"notifications"`
		Message       string                `json:"message"`
	}{
		Success:       true,
		Notifications: notifications,
		Message:       "notifications fetched successfully",
	})
}
