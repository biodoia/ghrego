package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
)

// Relation Repository
type RelationRepository struct {
	db *DB
}

func NewRelationRepository(db *DB) ports.RelationRepository {
	return &RelationRepository{db: db}
}

func (r *RelationRepository) GetByRepositoryID(ctx context.Context, repoID int) ([]domain.RepositoryRelation, error) {
	const query = `SELECT id, "sourceRepositoryId", "targetRepositoryId", "relationType", similarity, description, "createdAt" FROM "repositoryRelations" WHERE "sourceRepositoryId" = $1 ORDER BY similarity DESC`
	rows, err := r.db.Pool.Query(ctx, query, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.RepositoryRelation
	for rows.Next() {
		var i domain.RepositoryRelation
		if err := rows.Scan(&i.ID, &i.SourceRepositoryID, &i.TargetRepositoryID, &i.RelationType, &i.Similarity, &i.Description, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *RelationRepository) Create(ctx context.Context, relation *domain.RepositoryRelation) (int, error) {
	const query = `INSERT INTO "repositoryRelations" ("sourceRepositoryId", "targetRepositoryId", "relationType", similarity, description, "createdAt") VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	var id int
	err := r.db.Pool.QueryRow(ctx, query, relation.SourceRepositoryID, relation.TargetRepositoryID, relation.RelationType, relation.Similarity, relation.Description).Scan(&id)
	return id, err
}

// Unification Repository
type UnificationRepository struct {
	db *DB
}

func NewUnificationRepository(db *DB) ports.UnificationRepository {
	return &UnificationRepository{db: db}
}

func (r *UnificationRepository) Create(ctx context.Context, operation *domain.UnificationOperation) error {
	const query = `
		INSERT INTO "unificationOperations" (
			"userId", "operationId", "sourceRepositoryIds", "targetRepositoryName", "targetRepositoryUrl", 
			visibility, status, progress, "currentStep", "filesProcessed", "totalFiles", errors, "createdAt", "updatedAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()
		)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		operation.UserID, operation.OperationID, operation.SourceRepositoryIDs, operation.TargetRepositoryName,
		operation.TargetRepositoryURL, operation.Visibility, operation.Status, operation.Progress,
		operation.CurrentStep, operation.FilesProcessed, operation.TotalFiles, operation.Errors,
	)
	return err
}

func (r *UnificationRepository) Update(ctx context.Context, operationID uuid.UUID, updates map[string]interface{}) error {
	// Simple update logic (production should use dynamic builder)
	// For MVP, assuming specific updates or implementing generic builder again.
	// Implementing generic builder for completeness.
	if len(updates) == 0 {
		return nil
	}
	query := `UPDATE "unificationOperations" SET `
	args := []interface{}{}
	i := 1
	for k, v := range updates {
		query += fmt.Sprintf(`"%s" = $%d, `, k, i)
		args = append(args, v)
		i++
	}
	query += fmt.Sprintf(`"updatedAt" = NOW() WHERE "operationId" = $%d`, i)
	args = append(args, operationID)
	_, err := r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *UnificationRepository) GetByID(ctx context.Context, operationID uuid.UUID) (*domain.UnificationOperation, error) {
	const query = `SELECT id, "userId", "operationId", "sourceRepositoryIds", "targetRepositoryName", "targetRepositoryUrl", visibility, status, progress, "currentStep", "filesProcessed", "totalFiles", errors, "createdAt", "updatedAt", "completedAt" FROM "unificationOperations" WHERE "operationId" = $1`
	var op domain.UnificationOperation
	err := r.db.Pool.QueryRow(ctx, query, operationID).Scan(
		&op.ID, &op.UserID, &op.OperationID, &op.SourceRepositoryIDs, &op.TargetRepositoryName, &op.TargetRepositoryURL,
		&op.Visibility, &op.Status, &op.Progress, &op.CurrentStep, &op.FilesProcessed, &op.TotalFiles, &op.Errors,
		&op.CreatedAt, &op.UpdatedAt, &op.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &op, nil
}

func (r *UnificationRepository) GetByUserID(ctx context.Context, userID int) ([]domain.UnificationOperation, error) {
	const query = `SELECT id, "userId", "operationId", "sourceRepositoryIds", "targetRepositoryName", "targetRepositoryUrl", visibility, status, progress, "currentStep", "filesProcessed", "totalFiles", errors, "createdAt", "updatedAt", "completedAt" FROM "unificationOperations" WHERE "userId" = $1 ORDER BY "createdAt" DESC`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.UnificationOperation
	for rows.Next() {
		var op domain.UnificationOperation
		if err := rows.Scan(
			&op.ID, &op.UserID, &op.OperationID, &op.SourceRepositoryIDs, &op.TargetRepositoryName, &op.TargetRepositoryURL,
			&op.Visibility, &op.Status, &op.Progress, &op.CurrentStep, &op.FilesProcessed, &op.TotalFiles, &op.Errors,
			&op.CreatedAt, &op.UpdatedAt, &op.CompletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, op)
	}
	return items, nil
}
