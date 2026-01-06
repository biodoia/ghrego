package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	// Use Gemini 1.5 Flash for speed and cost, or Pro for complex analysis
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.2) // Low temperature for deterministic analysis

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (c *GeminiClient) Close() {
	c.client.Close()
}

func (c *GeminiClient) AnalyzeRepository(ctx context.Context, prompt string) (*domain.RepositoryAnalysisResponse, error) {
	// Configure JSON mode
	c.model.ResponseMIMEType = "application/json"
	c.model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(`Sei un esperto analista di codice e architetture software. 
Analizza il repository fornito e identifica:
1. Architettura e pattern utilizzati
2. Feature implementate con dettagli
3. Tecnologie e stack tecnologico
4. Pattern di design e best practices
5. Qualità del codice e aree di miglioramento
6. Suggerimenti concreti per ottimizzazione

Rispondi in formato JSON strutturato con i campi: architecture, features, technologies, patterns, quality, suggestions.`),
		},
	}

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("empty response from model")
	}

	var analysis domain.RepositoryAnalysisResponse
	var rawJSON string

	// Extract text from parts
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			rawJSON += string(txt)
		}
	}

	// Sanitize markdown code blocks if present (Gemini sometimes adds ```json ... ``` even in JSON mode)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	
	if err := json.Unmarshal([]byte(rawJSON), &analysis); err != nil {
		log.Error().Err(err).Str("raw", rawJSON).Msg("Failed to unmarshal JSON from LLM")
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return &analysis, nil
}

// Helper to handle JSON serialization for DB fields
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
