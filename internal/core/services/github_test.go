package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGitHubServiceImpl_SyncUserRepositories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockRepoStore := new(mocks.RepositoryStore)
		mockGHClient := new(mocks.GitHubClient)
		svc := NewGitHubService(mockGHClient, mockRepoStore, mockUserRepo)

		user := &domain.User{
			ID:             1,
			OpenID:         "test-openid",
			GithubUsername: domain.SQLNullString("testuser"),
		}
		
		ghRepos := []*domain.Repository{
			{Name: "repo1", GithubID: "101"},
		}

		mockUserRepo.On("GetByID", mock.Anything, 1).Return(user, nil)
		mockGHClient.On("GetUserRepositories", mock.Anything, "testuser").Return(ghRepos, nil)
		mockRepoStore.On("Upsert", mock.Anything, mock.MatchedBy(func(r *domain.Repository) bool {
			return r.Name == "repo1"
		})).Return(1, nil)

		err := svc.SyncUserRepositories(context.Background(), 1, "test-openid")
		assert.NoError(t, err)
		
		mockUserRepo.AssertExpectations(t)
		mockGHClient.AssertExpectations(t)
		mockRepoStore.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockRepoStore := new(mocks.RepositoryStore)
		mockGHClient := new(mocks.GitHubClient)
		svc := NewGitHubService(mockGHClient, mockRepoStore, mockUserRepo)

		mockUserRepo.On("GetByID", mock.Anything, 99).Return(nil, errors.New("not found"))
		
		err := svc.SyncUserRepositories(context.Background(), 99, "")
		assert.Error(t, err)
	})
	
	t.Run("github client error", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockRepoStore := new(mocks.RepositoryStore)
		mockGHClient := new(mocks.GitHubClient)
		svc := NewGitHubService(mockGHClient, mockRepoStore, mockUserRepo)

		user := &domain.User{
			ID:             1,
			OpenID:         "test-openid",
			GithubUsername: domain.SQLNullString("testuser"),
		}
		mockUserRepo.On("GetByID", mock.Anything, 1).Return(user, nil)
		mockGHClient.On("GetUserRepositories", mock.Anything, "testuser").Return(nil, errors.New("api error"))
		
		err := svc.SyncUserRepositories(context.Background(), 1, "")
		assert.Error(t, err)
		assert.Equal(t, "api error", err.Error())
	})
	
	t.Run("no github username linked", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockRepoStore := new(mocks.RepositoryStore)
		mockGHClient := new(mocks.GitHubClient)
		svc := NewGitHubService(mockGHClient, mockRepoStore, mockUserRepo)

		user := &domain.User{
			ID:             1,
			GithubUsername: sql.NullString{Valid: false},
		}
		mockUserRepo.On("GetByID", mock.Anything, 1).Return(user, nil)
		
		err := svc.SyncUserRepositories(context.Background(), 1, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no github username linked")
	})
}
