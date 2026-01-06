package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

// MockAnalysisRepository
type AnalysisRepository struct {
	mock.Mock
}

func (m *AnalysisRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Analysis, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Analysis), args.Error(1)
}

func (m *AnalysisRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Analysis, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Analysis), args.Error(1)
}

func (m *AnalysisRepository) Create(ctx context.Context, analysis *domain.Analysis) (int, error) {
	args := m.Called(ctx, analysis)
	return args.Int(0), args.Error(1)
}

func (m *AnalysisRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

// MockFeatureRepository
type FeatureRepository struct {
	mock.Mock
}

func (m *FeatureRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Feature, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Feature), args.Error(1)
}

func (m *FeatureRepository) Create(ctx context.Context, feature *domain.Feature) (int, error) {
	args := m.Called(ctx, feature)
	return args.Int(0), args.Error(1)
}

func (m *FeatureRepository) BulkCreate(ctx context.Context, features []domain.Feature) error {
	args := m.Called(ctx, features)
	return args.Error(0)
}

// MockTechnologyRepository
type TechnologyRepository struct {
	mock.Mock
}

func (m *TechnologyRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Technology, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Technology), args.Error(1)
}

func (m *TechnologyRepository) BulkCreate(ctx context.Context, techs []domain.Technology) error {
	args := m.Called(ctx, techs)
	return args.Error(0)
}

// MockSuggestionRepository
type SuggestionRepository struct {
	mock.Mock
}

func (m *SuggestionRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Suggestion, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Suggestion), args.Error(1)
}

func (m *SuggestionRepository) GetAllPending(ctx context.Context, userID int) ([]domain.Suggestion, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Suggestion), args.Error(1)
}

func (m *SuggestionRepository) GetByID(ctx context.Context, id int) (*domain.Suggestion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Suggestion), args.Error(1)
}

func (m *SuggestionRepository) Create(ctx context.Context, suggestion *domain.Suggestion) (int, error) {
	args := m.Called(ctx, suggestion)
	return args.Int(0), args.Error(1)
}

func (m *SuggestionRepository) UpdateStatus(ctx context.Context, id int, status domain.SuggestionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockRelationRepository
type RelationRepository struct {
	mock.Mock
}

func (m *RelationRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.RepositoryRelation, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RepositoryRelation), args.Error(1)
}

func (m *RelationRepository) Create(ctx context.Context, relation *domain.RepositoryRelation) (int, error) {
	args := m.Called(ctx, relation)
	return args.Int(0), args.Error(1)
}

// MockUnificationRepository
type UnificationRepository struct {
	mock.Mock
}

func (m *UnificationRepository) Create(ctx context.Context, operation *domain.UnificationOperation) error {
	args := m.Called(ctx, operation)
	return args.Error(0)
}

func (m *UnificationRepository) Update(ctx context.Context, operationID uuid.UUID, updates map[string]interface{}) error {
	args := m.Called(ctx, operationID, updates)
	return args.Error(0)
}

func (m *UnificationRepository) GetByID(ctx context.Context, operationID uuid.UUID) (*domain.UnificationOperation, error) {
	args := m.Called(ctx, operationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UnificationOperation), args.Error(1)
}

func (m *UnificationRepository) GetByUserID(ctx context.Context, userID int) ([]domain.UnificationOperation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UnificationOperation), args.Error(1)
}
