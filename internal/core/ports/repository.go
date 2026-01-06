package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/biodoia/ghrego/internal/core/domain"
)

// UserRepository defines operations for user management
type UserRepository interface {
	Upsert(ctx context.Context, user *domain.User) error
	GetByOpenID(ctx context.Context, openID string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

// RepositoryStore defines operations for GitHub repository management
type RepositoryStore interface {
	GetByUserID(ctx context.Context, userID int) ([]domain.Repository, error)
	GetByID(ctx context.Context, id int) (*domain.Repository, error)
	GetByIDs(ctx context.Context, ids []int) ([]domain.Repository, error)
	Upsert(ctx context.Context, repo *domain.Repository) (int, error)
	Delete(ctx context.Context, id int) error
	GetStats(ctx context.Context, userID int) (map[string]interface{}, error)
}

// AnalysisRepository defines operations for analysis results
type AnalysisRepository interface {
	GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Analysis, error)
	GetByUserID(ctx context.Context, userID int) ([]domain.Analysis, error)
	Create(ctx context.Context, analysis *domain.Analysis) (int, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
}

// FeatureRepository defines operations for detected features
type FeatureRepository interface {
	GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Feature, error)
	Create(ctx context.Context, feature *domain.Feature) (int, error)
	BulkCreate(ctx context.Context, features []domain.Feature) error
}

// RelationRepository defines operations for repository relationships
type RelationRepository interface {
	GetByRepositoryID(ctx context.Context, repoID int) ([]domain.RepositoryRelation, error)
	Create(ctx context.Context, relation *domain.RepositoryRelation) (int, error)
}

// SuggestionRepository defines operations for AI suggestions
type SuggestionRepository interface {
	GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Suggestion, error)
	GetAllPending(ctx context.Context, userID int) ([]domain.Suggestion, error)
	GetByID(ctx context.Context, id int) (*domain.Suggestion, error)
	Create(ctx context.Context, suggestion *domain.Suggestion) (int, error)
	UpdateStatus(ctx context.Context, id int, status domain.SuggestionStatus) error
}

// TechnologyRepository defines operations for technology stacks
type TechnologyRepository interface {
	GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Technology, error)
	BulkCreate(ctx context.Context, techs []domain.Technology) error
}

// UnificationRepository defines operations for repo unification
type UnificationRepository interface {
	Create(ctx context.Context, operation *domain.UnificationOperation) error
	Update(ctx context.Context, operationID uuid.UUID, updates map[string]interface{}) error
	GetByID(ctx context.Context, operationID uuid.UUID) (*domain.UnificationOperation, error)
	GetByUserID(ctx context.Context, userID int) ([]domain.UnificationOperation, error)
}
