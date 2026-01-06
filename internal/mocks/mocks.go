package mocks

import (
	"context"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Upsert(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) GetByOpenID(ctx context.Context, openID string) (*domain.User, error) {
	args := m.Called(ctx, openID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockRepositoryStore
type RepositoryStore struct {
	mock.Mock
}

func (m *RepositoryStore) GetByUserID(ctx context.Context, userID int) ([]domain.Repository, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Repository), args.Error(1)
}

func (m *RepositoryStore) GetByID(ctx context.Context, id int) (*domain.Repository, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Repository), args.Error(1)
}

func (m *RepositoryStore) GetByIDs(ctx context.Context, ids []int) ([]domain.Repository, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Repository), args.Error(1)
}

func (m *RepositoryStore) Upsert(ctx context.Context, repo *domain.Repository) (int, error) {
	args := m.Called(ctx, repo)
	return args.Int(0), args.Error(1)
}

func (m *RepositoryStore) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryStore) GetStats(ctx context.Context, userID int) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// MockGitHubService
type GitHubService struct {
	mock.Mock
}

func (m *GitHubService) SyncUserRepositories(ctx context.Context, userID int, openID string) error {
	args := m.Called(ctx, userID, openID)
	return args.Error(0)
}

func (m *GitHubService) GetRepositoryDetails(ctx context.Context, userID int, repoID int) (*domain.Repository, error) {
	args := m.Called(ctx, userID, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Repository), args.Error(1)
}

func (m *GitHubService) AnalyzeDependencies(ctx context.Context, repoID int) error {
	args := m.Called(ctx, repoID)
	return args.Error(0)
}

// MockGitHubClient
type GitHubClient struct {
	mock.Mock
}

func (m *GitHubClient) GetUserRepositories(ctx context.Context, username string) ([]*domain.Repository, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Repository), args.Error(1)
}

func (m *GitHubClient) GetRepository(ctx context.Context, owner, repoName string) (*domain.Repository, error) {
	args := m.Called(ctx, owner, repoName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Repository), args.Error(1)
}

func (m *GitHubClient) GetFileContent(ctx context.Context, owner, repo, path string) (string, error) {
	args := m.Called(ctx, owner, repo, path)
	return args.String(0), args.Error(1)
}

func (m *GitHubClient) GetLanguages(ctx context.Context, owner, repo string) (map[string]int, error) {
	args := m.Called(ctx, owner, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *GitHubClient) AnalyzeStructure(ctx context.Context, owner, repo string) (int, []string, map[string]int, error) {
	args := m.Called(ctx, owner, repo)
	return args.Int(0), args.Get(1).([]string), args.Get(2).(map[string]int), args.Error(3)
}

// MockAIClient
type AIClient struct {
	mock.Mock
}

func (m *AIClient) AnalyzeRepository(ctx context.Context, prompt string) (*domain.RepositoryAnalysisResponse, error) {
	args := m.Called(ctx, prompt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RepositoryAnalysisResponse), args.Error(1)
}

type AIAnalysisService struct {
	mock.Mock
}

func (m *AIAnalysisService) AnalyzeRepository(ctx context.Context, repoID int, analysisType domain.AnalysisType) (*domain.Analysis, error) {
	args := m.Called(ctx, repoID, analysisType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Analysis), args.Error(1)
}

func (m *AIAnalysisService) GenerateSuggestions(ctx context.Context, repoID int) ([]domain.Suggestion, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Suggestion), args.Error(1)
}
