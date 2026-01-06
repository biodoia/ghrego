package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) ports.UserRepository {
	return &UserRepository{db: db}
}

// Upsert inserts or updates a user
func (r *UserRepository) Upsert(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (
			"openId", name, email, "loginMethod", role, 
			"githubUsername", "githubId", "lastSignedIn", "updatedAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW()
		)
		ON CONFLICT ("openId") DO UPDATE SET
			name = COALESCE(EXCLUDED.name, users.name),
			email = COALESCE(EXCLUDED.email, users.email),
			"loginMethod" = COALESCE(EXCLUDED."loginMethod", users."loginMethod"),
			role = CASE WHEN users.role IS NULL THEN EXCLUDED.role ELSE users.role END,
			"githubUsername" = COALESCE(EXCLUDED."githubUsername", users."githubUsername"),
			"githubId" = COALESCE(EXCLUDED."githubId", users."githubId"),
			"lastSignedIn" = EXCLUDED."lastSignedIn",
			"updatedAt" = NOW()
		RETURNING id, "createdAt"
	`

	// Handle role logic: if role is empty and openID matches owner (env var?), set admin.
	// For now, we trust the domain object has the correct role set by the service layer.

	role := string(user.Role)
	if role == "" {
		role = "user"
	}

	lastSignedIn := user.LastSignedIn
	if lastSignedIn.IsZero() {
		lastSignedIn = time.Now()
	}

	err := r.db.Pool.QueryRow(ctx, query,
		user.OpenID,
		user.Name,
		user.Email,
		user.LoginMethod,
		role,
		user.GithubUsername,
		user.GithubID,
		lastSignedIn,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

// GetByOpenID retrieves a user by their OpenID
func (r *UserRepository) GetByOpenID(ctx context.Context, openID string) (*domain.User, error) {
	const query = `
		SELECT id, "openId", name, email, "loginMethod", role, 
		       "githubUsername", "githubId", "createdAt", "updatedAt", "lastSignedIn"
		FROM users
		WHERE "openId" = $1
	`

	var user domain.User
	var roleStr string

	err := r.db.Pool.QueryRow(ctx, query, openID).Scan(
		&user.ID,
		&user.OpenID,
		&user.Name,
		&user.Email,
		&user.LoginMethod,
		&roleStr,
		&user.GithubUsername,
		&user.GithubID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastSignedIn,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found is not an error in domain logic usually, or handle as specific error
		}
		return nil, fmt.Errorf("failed to get user by openId: %w", err)
	}

	user.Role = domain.UserRole(roleStr)
	return &user, nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	const query = `
		SELECT id, "openId", name, email, "loginMethod", role, 
		       "githubUsername", "githubId", "createdAt", "updatedAt", "lastSignedIn"
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var roleStr string

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.OpenID,
		&user.Name,
		&user.Email,
		&user.LoginMethod,
		&roleStr,
		&user.GithubUsername,
		&user.GithubID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastSignedIn,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user.Role = domain.UserRole(roleStr)
	return &user, nil
}
