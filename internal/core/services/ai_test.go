package services

import (
	"context"
	"errors"
	"testing"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAIAnalysisServiceImpl_AnalyzeRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockAIClient := new(mocks.AIClient)
		mockRepoStore := new(mocks.RepositoryStore)
		mockAnalysisRepo := new(mocks.AnalysisRepository)
		mockFeatureRepo := new(mocks.FeatureRepository)
		mockTechRepo := new(mocks.TechnologyRepository)
		mockSuggRepo := new(mocks.SuggestionRepository)

		svc := NewAIAnalysisService(mockAIClient, mockRepoStore, mockAnalysisRepo, mockFeatureRepo, mockTechRepo, mockSuggRepo)

		// Setup Data
		repo := &domain.Repository{
			ID: 1, Name: "repo1", Description: domain.SQLNullString("desc"), Language: domain.SQLNullString("Go"),
		}
		
		analysisResponse := &domain.RepositoryAnalysisResponse{
			Architecture: "Microservices",
			Features: []struct{Name string `json:"name"`; Description string `json:"description"`; Category string `json:"category"`; Confidence int `json:"confidence"`; FilePaths []string `json:"filePaths"`}{
				{Name: "Auth", Description: "Login", Confidence: 90},
			},
		}

		// Expectations
		mockRepoStore.On("GetByID", mock.Anything, 1).Return(repo, nil)
		mockAIClient.On("AnalyzeRepository", mock.Anything, mock.AnythingOfType("string")).Return(analysisResponse, nil)
		
		// Expect saves
		mockAnalysisRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Analysis")).Return(100, nil)
		mockFeatureRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]domain.Feature")).Return(nil)
		mockTechRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]domain.Technology")).Return(nil)
		mockSuggRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Suggestion")).Return(1, nil).Maybe() // Depends if suggestions are present

		// Execute
		res, err := svc.AnalyzeRepository(context.Background(), 1, domain.AnalysisTypeArchitecture)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 100, res.ID)
		mockAIClient.AssertExpectations(t)
	})

	t.Run("repo not found", func(t *testing.T) {
		mockRepoStore := new(mocks.RepositoryStore)
		svc := NewAIAnalysisService(nil, mockRepoStore, nil, nil, nil, nil)
		
		mockRepoStore.On("GetByID", mock.Anything, 99).Return(nil, nil) // or error
		// Note: implementation checks if repo == nil -> error "repository not found"
		
		_, err := svc.AnalyzeRepository(context.Background(), 99, domain.AnalysisTypeArchitecture)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ai client error", func(t *testing.T) {
		mockAIClient := new(mocks.AIClient)
		mockRepoStore := new(mocks.RepositoryStore)
		svc := NewAIAnalysisService(mockAIClient, mockRepoStore, nil, nil, nil, nil)

		repo := &domain.Repository{ID: 1}
		mockRepoStore.On("GetByID", mock.Anything, 1).Return(repo, nil)
		mockAIClient.On("AnalyzeRepository", mock.Anything, mock.Anything).Return(nil, errors.New("ai error"))
		
		_, err := svc.AnalyzeRepository(context.Background(), 1, domain.AnalysisTypeArchitecture)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AI analysis failed")
	})
}