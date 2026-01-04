package repo

import (
	"context"
	"fmt"
	"strings"

	"mlm/models"

	"github.com/aarondl/sqlboiler/v4/boil"
)

// UserRepo handles Insert/Update operations (returns pgmodel types)
type UserRepo struct {
	// Dependencies if needed
}

// NewUserRepo creates a new user repository
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// Insert creates a new user in the database
func (r *UserRepo) Insert(
	ctx context.Context,
	exec boil.ContextExecutor,
	user *models.User,
) (*models.User, error) {
	err := user.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

// BulkInsert - REQUIRED for multiple records (NO LOOPS)
// Inserts multiple users in a single query for performance
func (r *UserRepo) BulkInsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	users []*models.User,
) error {
	if len(users) == 0 {
		return nil
	}

	// Build placeholders for values
	// For each user: (?, ?, ?, ?, ?)
	placeholders := make([]string, len(users))
	args := make([]interface{}, 0, len(users)*5)

	for i, user := range users {
		placeholders[i] = "(?, ?, ?, ?, ?)"
		args = append(args,
			user.ID,
			user.Username,
			user.DisplayName,
			user.Gender,
			user.CreatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO users (id, username, display_name, gender, created_at)
		VALUES %s
	`, strings.Join(placeholders, ", "))

	_, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("bulk insert users: %w", err)
	}

	return nil
}

// Upsert inserts or updates a user
func (r *UserRepo) Upsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	user *models.User,
) (*models.User, error) {
	err := user.Upsert(
		ctx,
		exec,
		boil.Infer(), // update columns
		boil.Infer(), // insert columns
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	return user, nil
}
