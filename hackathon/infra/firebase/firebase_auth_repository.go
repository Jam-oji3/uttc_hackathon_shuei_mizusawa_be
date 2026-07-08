package firebase

import (
	"context"
	firebaseAdmin "firebase.google.com/go/v4" // ← ここで別名をつける
	"firebase.google.com/go/v4/auth"
	"fmt"
)

type FirebaseAuthRepository struct {
	client *auth.Client
}

func NewFirebaseAuthRepository(ctx context.Context) (*FirebaseAuthRepository, error) {
	app, err := firebaseAdmin.NewApp(ctx, nil) // ← 別名で呼び出す
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}

	return &FirebaseAuthRepository{client: client}, nil
}

func (r *FirebaseAuthRepository) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := r.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid ID token: %w", err)
	}
	return token, nil
}
