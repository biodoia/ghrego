package services

import (
	"context"
	"fmt"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
	"github.com/rs/zerolog/log"
)

type AIAnalysisServiceImpl struct {
	aiClient       ports.AIClient
	repoStore      ports.RepositoryStore
	analysisRepo   ports.AnalysisRepository
	featureRepo    ports.FeatureRepository
	technologyRepo ports.TechnologyRepository
	suggestionRepo ports.SuggestionRepository
}

func NewAIAnalysisService(
	aiClient ports.AIClient,
	repoStore ports.RepositoryStore,
	analysisRepo ports.AnalysisRepository,
	featureRepo ports.FeatureRepository,
	technologyRepo ports.TechnologyRepository,
	suggestionRepo ports.SuggestionRepository,
) ports.AIAnalysisService {
	return &AIAnalysisServiceImpl{
		aiClient:       aiClient,
		repoStore:      repoStore,
		analysisRepo:   analysisRepo,
		featureRepo:    featureRepo,
		technologyRepo: technologyRepo,
		suggestionRepo: suggestionRepo,
	}
}

func (s *AIAnalysisServiceImpl) AnalyzeRepository(ctx context.Context, repoID int, analysisType domain.AnalysisType) (*domain.Analysis, error) {
	// 1. Fetch Repository Details
	repo, err := s.repoStore.GetByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	if repo == nil {
		return nil, fmt.Errorf("repository not found")
	}

	// 2. Build Prompt (Simplified for now - in production would fetch files/readme via GitHubService)
	// Ideally we need GitHubService here too to fetch README and Dependencies.
	// For this refactor, I'll assume we pass enough info or fetch it.
	// To keep "lean", I'll just use the repo metadata we have.
	prompt := fmt.Sprintf(`
Analyze this repository:
Name: %s
Description: %s
Language: %s
Stars: %d
URL: %s
`, repo.FullName, repo.Description.String, repo.Language.String, repo.Stars, repo.URL)

	// 3. Call AI
	log.Info().Str("repo", repo.FullName).Msg("Starting AI analysis...")
	response, err := s.aiClient.AnalyzeRepository(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	// 4. Save Results
	analysis, features, techs, suggestions := response.ToDomain(repoID)
	analysis.AnalysisType = analysisType
	
	// Transaction would be better here, but doing sequential for now
	analysisID, err := s.analysisRepo.Create(ctx, analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}
	analysis.ID = analysisID

	// Save Features
	if len(features) > 0 {
		if err := s.featureRepo.BulkCreate(ctx, features); err != nil {
			log.Error().Err(err).Msg("Failed to save features")
		}
	}

	// Save Technologies
	if len(techs) > 0 {
		if err := s.technologyRepo.BulkCreate(ctx, techs); err != nil {
			log.Error().Err(err).Msg("Failed to save technologies")
		}
	}

	// Save Suggestions
	for _, sugg := range suggestions {
		if _, err := s.suggestionRepo.Create(ctx, &sugg); err != nil {
			log.Error().Err(err).Msg("Failed to save suggestion")
		}
	}

	log.Info().Int("analysis_id", analysisID).Msg("AI Analysis completed and saved")
	return analysis, nil
}

func (s *AIAnalysisServiceImpl) GenerateSuggestions(ctx context.Context, repoID int) ([]domain.Suggestion, error) {
	// Reuse AnalyzeRepository logic or simpler prompt
	// For now, return what we have in DB
	return s.suggestionRepo.GetByRepositoryID(ctx, repoID)
}
