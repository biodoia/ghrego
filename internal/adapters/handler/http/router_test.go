package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biodoia/ghrego/internal/config"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_handleGetRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepoStore := new(mocks.RepositoryStore)
		mockUserRepo := new(mocks.UserRepository)
		server := NewServer(&config.Config{Port: "8080"}, nil, nil, mockRepoStore, mockUserRepo, nil)

		repo := &domain.Repository{ID: 10, Name: "my-repo"}
		mockRepoStore.On("GetByID", mock.Anything, 10).Return(repo, nil)

		req := httptest.NewRequest("GET", "/api/repositories/10", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var respRepo domain.Repository
		json.Unmarshal(rr.Body.Bytes(), &respRepo)
		assert.Equal(t, "my-repo", respRepo.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepoStore := new(mocks.RepositoryStore)
		mockUserRepo := new(mocks.UserRepository)
		server := NewServer(&config.Config{Port: "8080"}, nil, nil, mockRepoStore, mockUserRepo, nil)

		mockRepoStore.On("GetByID", mock.Anything, 99).Return(nil, nil)

		req := httptest.NewRequest("GET", "/api/repositories/99", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestServer_handleSyncRepositories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockGHService := new(mocks.GitHubService)
		mockRepoStore := new(mocks.RepositoryStore)
		mockUserRepo := new(mocks.UserRepository)
		server := NewServer(&config.Config{Port: "8080"}, mockGHService, nil, mockRepoStore, mockUserRepo, nil)

		user := &domain.User{ID: 1, OpenID: "open-123"}
		mockUserRepo.On("GetByID", mock.Anything, 1).Return(user, nil)
		mockGHService.On("SyncUserRepositories", mock.Anything, 1, "open-123").Return(nil)
		mockRepoStore.On("GetByUserID", mock.Anything, 1).Return([]domain.Repository{{ID: 1}}, nil)

		req := httptest.NewRequest("POST", "/api/repositories/sync", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"success":true`)
	})
	
	t.Run("auth error", func(t *testing.T) {
		mockGHService := new(mocks.GitHubService)
		mockRepoStore := new(mocks.RepositoryStore)
		mockUserRepo := new(mocks.UserRepository)
		server := NewServer(&config.Config{Port: "8080"}, mockGHService, nil, mockRepoStore, mockUserRepo, nil)

		// Mock GetByID failing (e.g. user not found)
		mockUserRepo.On("GetByID", mock.Anything, 1).Return(nil, errors.New("db err"))

		req := httptest.NewRequest("POST", "/api/repositories/sync", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}