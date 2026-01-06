package postgres

import (
	"context"
	"fmt"

	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
)

type AnalysisRepository struct {
	db *DB
}

func NewAnalysisRepository(db *DB) ports.AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.Analysis, error) {
	const query = `
		SELECT id, "repositoryId", "analysisType", status, result, summary, score, "errorMessage", "createdAt", "completedAt"
		FROM analyses
		WHERE "repositoryId" = $1
		ORDER BY "createdAt" DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []domain.Analysis
	for rows.Next() {
		var a domain.Analysis
		// Note: Enums are strings in DB usually, need casting if driver doesn't handle custom types automatically.
		// pgx handles string -> custom string type well.
		if err := rows.Scan(
			&a.ID, &a.RepositoryID, &a.AnalysisType, &a.Status, &a.Result, &a.Summary,
			&a.Score, &a.ErrorMessage, &a.CreatedAt, &a.CompletedAt,
		); err != nil {
			return nil, err
		}
		analyses = append(analyses, a)
	}
	return analyses, nil
}

func (r *AnalysisRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Analysis, error) {
	const query = `
		SELECT a.id, a."repositoryId", a."analysisType", a.status, a.result, a.summary, a.score, a."errorMessage", a."createdAt", a."completedAt"
		FROM analyses a
		JOIN repositories r ON a."repositoryId" = r.id
		WHERE r."userId" = $1
		ORDER BY a."createdAt" DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []domain.Analysis
	for rows.Next() {
		var a domain.Analysis
		if err := rows.Scan(
			&a.ID, &a.RepositoryID, &a.AnalysisType, &a.Status, &a.Result, &a.Summary,
			&a.Score, &a.ErrorMessage, &a.CreatedAt, &a.CompletedAt,
		); err != nil {
			return nil, err
		}
		analyses = append(analyses, a)
	}
	return analyses, nil
}

func (r *AnalysisRepository) Create(ctx context.Context, analysis *domain.Analysis) (int, error) {
	const query = `
		INSERT INTO analyses (
			"repositoryId", "analysisType", status, result, summary, score, "errorMessage", "createdAt", "completedAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), $8
		)
		RETURNING id, "createdAt"
	`
	err := r.db.Pool.QueryRow(ctx, query,
		analysis.RepositoryID,
		analysis.AnalysisType,
		analysis.Status,
		analysis.Result,
		analysis.Summary,
		analysis.Score,
		analysis.ErrorMessage,
		analysis.CompletedAt,
	).Scan(&analysis.ID, &analysis.CreatedAt)

	if err != nil {
		return 0, fmt.Errorf("failed to create analysis: %w", err)
	}
	return analysis.ID, nil
}

func (r *AnalysisRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	// Dynamic update query builder
	// Warning: Input map keys must be safe (no sql injection). 
	// In strict porting, we should likely have a concrete struct or fixed fields.
	// For now, I'll implement a simple status update as fallback or use a safe builder.
	
	if len(updates) == 0 {
		return nil
	}
	
	query := "UPDATE analyses SET "
	args := []interface{}{}
	i := 1
	for k, v := range updates {
		// Whitelist columns for safety
		switch k {
		case "status", "result", "summary", "score", "errorMessage", "completedAt":
			query += fmt.Sprintf(`"%s" = $%d, `, k, i)
			args = append(args, v)
			i++
		}
	}
	query = query[:len(query)-2] // remove trailing comma
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	_, err := r.db.Pool.Exec(ctx, query, args...)
	return err
}
