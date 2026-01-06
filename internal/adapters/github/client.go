package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	cache  *cache.Cache
}

func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Initialize cache with 5 minute default expiration and 10 minute cleanup interval
	c := cache.New(5*time.Minute, 10*time.Minute)

	return &Client{
		client: client,
		cache:  c,
	}
}

// GetUserRepositories retrieves all repositories for a user
func (c *Client) GetUserRepositories(ctx context.Context, username string) ([]*domain.Repository, error) {
	cacheKey := fmt.Sprintf("repos:%s", username)
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.([]*domain.Repository), nil
	}

	opts := &github.RepositoryListOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Map to domain.Repository
	var domainRepos []*domain.Repository
	for _, r := range allRepos {
		domainRepos = append(domainRepos, mapGitHubRepoToDomain(r))
	}

	c.cache.Set(cacheKey, domainRepos, cache.DefaultExpiration)
	return domainRepos, nil
}

// GetRepository retrieves detailed information about a repository
func (c *Client) GetRepository(ctx context.Context, owner, repoName string) (*domain.Repository, error) {
	cacheKey := fmt.Sprintf("repo:%s/%s", owner, repoName)
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.(*domain.Repository), nil
	}

	repo, _, err := c.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	domainRepo := mapGitHubRepoToDomain(repo)
	c.cache.Set(cacheKey, domainRepo, cache.DefaultExpiration)
	return domainRepo, nil
}

// GetFileContent retrieves content of a file
func (c *Client) GetFileContent(ctx context.Context, owner, repo, path string) (string, error) {
	fileContent, _, _, err := c.client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		// Handle 404
		if strings.Contains(err.Error(), "404") {
			return "", nil
		}
		return "", err
	}

	if fileContent == nil {
		return "", fmt.Errorf("file content is empty")
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}
	return content, nil
}

// GetLanguages retrieves language statistics
func (c *Client) GetLanguages(ctx context.Context, owner, repo string) (map[string]int, error) {
	langs, _, err := c.client.Repositories.ListLanguages(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return langs, nil
}

// AnalyzeStructure performs a tree analysis (equivalent to analyzeRepositoryStructure)
func (c *Client) AnalyzeStructure(ctx context.Context, owner, repo string) (int, []string, map[string]int, error) {
	repoData, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return 0, nil, nil, err
	}
	branch := repoData.GetDefaultBranch()

	tree, _, err := c.client.Git.GetTree(ctx, owner, repo, branch, true)
	if err != nil {
		return 0, nil, nil, err
	}

	var totalFiles int
	var dirs []string
	fileTypes := make(map[string]int)

	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			totalFiles++
			path := entry.GetPath()
			parts := strings.Split(path, ".")
			ext := "no-extension"
			if len(parts) > 1 {
				ext = parts[len(parts)-1]
			}
			fileTypes[ext]++
		} else if entry.GetType() == "tree" {
			dirs = append(dirs, entry.GetPath())
		}
	}

	return totalFiles, dirs, fileTypes, nil
}

// Helper to map GitHub struct to Domain struct
func mapGitHubRepoToDomain(ghRepo *github.Repository) *domain.Repository {
	// Note: You need to handle sql.Null* types or use helper functions
	// For simplicity in this step, I am focusing on the struct logic.
	// Real implementation needs to handle pointers.
	
	repo := &domain.Repository{
		GithubID:      fmt.Sprintf("%d", ghRepo.GetID()),
		Name:          ghRepo.GetName(),
		FullName:      ghRepo.GetFullName(),
		URL:           ghRepo.GetHTMLURL(),
		IsPrivate:     ghRepo.GetPrivate(),
		Stars:         ghRepo.GetStargazersCount(),
		Forks:         ghRepo.GetForksCount(),
		Size:          ghRepo.GetSize(),
		DefaultBranch: ghRepo.GetDefaultBranch(),
		CreatedAt:     ghRepo.GetCreatedAt().Time,
		UpdatedAt:     ghRepo.GetUpdatedAt().Time,
	}

	if ghRepo.Description != nil {
		repo.Description.String = *ghRepo.Description
		repo.Description.Valid = true
	}
	if ghRepo.Language != nil {
		repo.Language.String = *ghRepo.Language
		repo.Language.Valid = true
	}

	return repo
}
