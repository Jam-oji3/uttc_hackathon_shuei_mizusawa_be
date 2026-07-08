package repository

import "context"

type GeminiClient interface {
	GenerateContent(ctx context.Context, prompt string) (*string, error)
}
