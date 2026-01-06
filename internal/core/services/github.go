package services

import (
	"context"
	"fmt"
	"time"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
	"github.com/rs/zerolog/log"
)

type GitHubServiceImpl struct {
	ghClient  ports.GitHubClient
	repoStore ports.RepositoryStore
	userRepo  ports.UserRepository
}

func NewGitHubService(ghClient ports.GitHubClient, repoStore ports.RepositoryStore, userRepo ports.UserRepository) ports.GitHubService {
	return &GitHubServiceImpl{
		ghClient:  ghClient,
		repoStore: repoStore,
		userRepo:  userRepo,
	}
}

func (s *GitHubServiceImpl) SyncUserRepositories(ctx context.Context, userID int, openID string) error {
	// 1. Get User to get (potentially) the stored token? 
	// In the original TS code, the token was passed to the constructor or derived.
	// Here we assume the client is initialized with a token, BUT in a multi-user app,
	// the client needs to be created per-request or the token passed here.
	// For now, I'll assume the ghClient passed is generic or the service handles token retrieval.
	// Actually, the original 'GitHubService' was instantiated with a token.
	// So we might need a Factory or pass the token to methods.
	// For simplicity in this step, I'll fetch the user and assume we have a way to get their token.
	// If the token is NOT stored in DB (OAuth), we might need to pass it from the controller.
	
	// Refactoring note: better to pass token or use an OAuth provider service.
	// I will assume for now that we get the repos via the existing client which might use a system token or user token.
	// TO FIX: Pass token or username logic.
	
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if !user.GithubUsername.Valid {
		return fmt.Errorf("user has no github username linked")
	}

	repos, err := s.ghClient.GetUserRepositories(ctx, user.GithubUsername.String)
	if err != nil {
		return err
	}

	log.Info().Int("count", len(repos)).Str("user", user.GithubUsername.String).Msg("Fetched repositories from GitHub")

	for _, repo := range repos {
		repo.UserID = userID
		repo.LastSyncAt.Time = time.Now()
		repo.LastSyncAt.Valid = true
		
		_, err := s.repoStore.Upsert(ctx, repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo.Name).Msg("Failed to upsert repository")
			// Continue with others
		}
	}

	return nil
}

func (s *GitHubServiceImpl) GetRepositoryDetails(ctx context.Context, userID int, repoID int) (*domain.Repository, error) {
	return s.repoStore.GetByID(ctx, repoID)
}

func (s *GitHubServiceImpl) AnalyzeDependencies(ctx context.Context, repoID int) error {
	// 1. Get Repo
	repo, err := s.repoStore.GetByID(ctx, repoID)
	if err != nil {
		return err
	}

	// 2. Fetch file content (package.json, go.mod, etc.)
	// This requires the GitHub client.
	// We need the Owner and RepoName.
	// domain.Repository has "FullName" which is usually "owner/repo".
	
	// Implementation to be added: dependency parsing logic similar to TS.
	// For now, stubbed.
	log.Info().Str("repo", repo.FullName).Msg("Analyzing dependencies...")
	
	return nil
}
