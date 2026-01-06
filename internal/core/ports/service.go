package ports

import (
	"context"

	"github.com/biodoia/ghrego/internal/core/domain"
)

// External Adapter Interfaces
type GitHubClient interface {
	GetUserRepositories(ctx context.Context, username string) ([]*domain.Repository, error)
	GetRepository(ctx context.Context, owner, repoName string) (*domain.Repository, error)
	GetFileContent(ctx context.Context, owner, repo, path string) (string, error)
	GetLanguages(ctx context.Context, owner, repo string) (map[string]int, error)
	AnalyzeStructure(ctx context.Context, owner, repo string) (int, []string, map[string]int, error)
}

type AIClient interface {
	AnalyzeRepository(ctx context.Context, prompt string) (*domain.RepositoryAnalysisResponse, error)
}

// Service Interfaces
type GitHubService interface {
	SyncUserRepositories(ctx context.Context, userID int, openID string) error
	GetRepositoryDetails(ctx context.Context, userID int, repoID int) (*domain.Repository, error)
	AnalyzeDependencies(ctx context.Context, repoID int) error
}

type AIAnalysisService interface {
	AnalyzeRepository(ctx context.Context, repoID int, analysisType domain.AnalysisType) (*domain.Analysis, error)
	GenerateSuggestions(ctx context.Context, repoID int) ([]domain.Suggestion, error)
}