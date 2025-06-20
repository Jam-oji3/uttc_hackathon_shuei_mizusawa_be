package gemini

import (
	"context"
	"fmt"
	"hackathon/repository"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

type GeminiGateway struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

var _ repository.GeminiClient = (*GeminiGateway)(nil)

func NewGeminiGateway(ctx context.Context) (*GeminiGateway, error) {
	projectID := os.Getenv("VERTEXAI_PROJECT_ID")
	location := os.Getenv("VERTEXAI_LOCATION")
	modelName := os.Getenv("VERTEXAI_MODEL_NAME")

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	model := client.GenerativeModel(modelName)

	return &GeminiGateway{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiGateway) GenerateContent(ctx context.Context, prompt string) (*string, error) {
	promptPart := genai.Text(prompt)
	resp, err := g.model.GenerateContent(ctx, promptPart)
	if err != nil {
		return nil, fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	generatedText := fmt.Sprint(resp.Candidates[0].Content.Parts[0])

	return &generatedText, nil
}

// アプリケーション終了時にClientをクローズするためのメソッド
func (g *GeminiGateway) Close() {
	g.client.Close()
}
