package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
)

// Feature Repository
type FeatureRepository struct {
	db *DB
}

func NewFeatureRepository(db *DB) ports.FeatureRepository {
	return &FeatureRepository{db: db}
}

func (r *FeatureRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Feature, error) {
	const query = `SELECT id, "repositoryId", name, description, category, "filePaths", "codeSnippet", confidence, "createdAt" FROM features WHERE "repositoryId" = $1`
	rows, err := r.db.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Feature
	for rows.Next() {
		var i domain.Feature
		if err := rows.Scan(&i.ID, &i.RepositoryID, &i.Name, &i.Description, &i.Category, &i.FilePaths, &i.CodeSnippet, &i.Confidence, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *FeatureRepository) Create(ctx context.Context, feature *domain.Feature) (int, error) {
	const query = `INSERT INTO features ("repositoryId", name, description, category, "filePaths", "codeSnippet", confidence, "createdAt") VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`
	var id int
	err := r.db.Pool.QueryRow(ctx, query, feature.RepositoryID, feature.Name, feature.Description, feature.Category, feature.FilePaths, feature.CodeSnippet, feature.Confidence).Scan(&id)
	return id, err
}

func (r *FeatureRepository) BulkCreate(ctx context.Context, features []domain.Feature) error {
	if len(features) == 0 {
		return nil
	}
	// Using CopyFrom for bulk insert is most efficient in pgx
	rows := [][]interface{}{}
	for _, f := range features {
		rows = append(rows, []interface{}{f.RepositoryID, f.Name, f.Description, f.Category, f.FilePaths, f.CodeSnippet, f.Confidence, f.CreatedAt}) // assuming CreatedAt is set or default
	}

	_, err := r.db.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"features"},
		[]string{"repositoryId", "name", "description", "category", "filePaths", "codeSnippet", "confidence", "createdAt"},
		pgx.CopyFromRows(rows),
	)
	return err
}

// Technology Repository
type TechnologyRepository struct {
	db *DB
}

func NewTechnologyRepository(db *DB) ports.TechnologyRepository {
	return &TechnologyRepository{db: db}
}

func (r *TechnologyRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Technology, error) {
	const query = `SELECT id, "repositoryId", name, version, type, "packageManager", "createdAt" FROM technologies WHERE "repositoryId" = $1`
	rows, err := r.db.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Technology
	for rows.Next() {
		var i domain.Technology
		if err := rows.Scan(&i.ID, &i.RepositoryID, &i.Name, &i.Version, &i.Type, &i.PackageManager, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *TechnologyRepository) BulkCreate(ctx context.Context, techs []domain.Technology) error {
	if len(techs) == 0 {
		return nil
	}
	rows := [][]interface{}{}
	for _, t := range techs {
		rows = append(rows, []interface{}{t.RepositoryID, t.Name, t.Version, t.Type, t.PackageManager, t.CreatedAt})
	}
	_, err := r.db.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"technologies"},
		[]string{"repositoryId", "name", "version", "type", "packageManager", "createdAt"},
		pgx.CopyFromRows(rows),
	)
	return err
}

// Suggestion Repository
type SuggestionRepository struct {
	db *DB
}

func NewSuggestionRepository(db *DB) ports.SuggestionRepository {
	return &SuggestionRepository{db: db}
}

func (r *SuggestionRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Suggestion, error) {
	const query = `SELECT id, "repositoryId", "suggestionType", title, description, "sourceRepositoryId", priority, status, "createdAt", "updatedAt" FROM suggestions WHERE "repositoryId" = $1`
	rows, err := r.db.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Suggestion
	for rows.Next() {
		var i domain.Suggestion
		if err := rows.Scan(&i.ID, &i.RepositoryID, &i.SuggestionType, &i.Title, &i.Description, &i.SourceRepositoryID, &i.Priority, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *SuggestionRepository) GetAllPending(ctx context.Context, userID int) ([]domain.Suggestion, error) {
	const query = `
		SELECT s.id, s."repositoryId", s."suggestionType", s.title, s.description, s."sourceRepositoryId", s.priority, s.status, s."createdAt", s."updatedAt"
		FROM suggestions s
		JOIN repositories r ON s."repositoryId" = r.id
		WHERE r."userId" = $1 AND s.status = 'pending'
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Suggestion
	for rows.Next() {
		var i domain.Suggestion
		if err := rows.Scan(&i.ID, &i.RepositoryID, &i.SuggestionType, &i.Title, &i.Description, &i.SourceRepositoryID, &i.Priority, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *SuggestionRepository) GetByID(ctx context.Context, id int) (*domain.Suggestion, error) {
	const query = `SELECT id, "repositoryId", "suggestionType", title, description, "sourceRepositoryId", priority, status, "createdAt", "updatedAt" FROM suggestions WHERE id = $1`
	var i domain.Suggestion
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&i.ID, &i.RepositoryID, &i.SuggestionType, &i.Title, &i.Description, &i.SourceRepositoryID, &i.Priority, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

func (r *SuggestionRepository) Create(ctx context.Context, suggestion *domain.Suggestion) (int, error) {
	const query = `
		INSERT INTO suggestions (
			"repositoryId", "suggestionType", title, description, "sourceRepositoryId", priority, status, "createdAt", "updatedAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
		) RETURNING id
	`
	var id int
	err := r.db.Pool.QueryRow(ctx, query, 
		suggestion.RepositoryID, suggestion.SuggestionType, suggestion.Title, suggestion.Description, 
		suggestion.SourceRepositoryID, suggestion.Priority, suggestion.Status,
	).Scan(&id)
	return id, err
}

func (r *SuggestionRepository) UpdateStatus(ctx context.Context, id int, status domain.SuggestionStatus) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE suggestions SET status = $1, "updatedAt" = NOW() WHERE id = $2`, status, id)
	return err
}
