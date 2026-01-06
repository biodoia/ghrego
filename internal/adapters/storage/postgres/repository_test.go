package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &UserRepository{
		db: &DB{Pool: mock},
	}

	t.Run("success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "openId", "name", "email", "loginMethod", "role", "githubUsername", "githubId", "createdAt", "updatedAt", "lastSignedIn"}).
			AddRow(1, "open-123", "John", "john@test.com", "google", "user", "johngh", "1001", time.Now(), time.Now(), time.Now())

		mock.ExpectQuery(`SELECT id, "openId", name, email, "loginMethod", role`).
			WithArgs(1).
			WillReturnRows(rows)

		user, err := repo.GetByID(context.Background(), 1)
		
		assert.NoError(t, err)
		assert.Equal(t, "John", user.Name.String)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id`).
			WithArgs(99).
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetByID(context.Background(), 99)
		
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestRepositoryStore_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repoStore := &RepositoryStore{
		db: &DB{Pool: mock},
	}
	
	t.Run("success", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "userId", "githubId", "name", "fullName", "description", "url", "language", 
			"isPrivate", "stars", "forks", "size", "defaultBranch", "lastCommitAt", "lastSyncAt", 
			"createdAt", "updatedAt",
		}).
		AddRow(10, 1, "gh-10", "my-repo", "owner/my-repo", "desc", "url", "Go", false, 5, 1, 100, "main", nil, nil, time.Now(), time.Now())
		
		mock.ExpectQuery(`SELECT .* FROM repositories WHERE "userId" = \$1`).
			WithArgs(1).
			WillReturnRows(rows)
			
		repos, err := repoStore.GetByUserID(context.Background(), 1)
		
		assert.NoError(t, err)
		assert.Len(t, repos, 1)
		assert.Equal(t, "my-repo", repos[0].Name)
	})
}
