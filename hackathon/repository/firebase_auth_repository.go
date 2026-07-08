package repository

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthRepository interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}
