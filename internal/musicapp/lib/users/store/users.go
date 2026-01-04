package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"mlm/internal/musicapp/lib/users"
	"mlm/models" // SQLBoiler generated models
)

// Store handles user queries
type Store struct {
	// No dependencies - store is pure query logic
}

// New creates a new user store
func New() *Store {
	return &Store{}
}

// Users returns 0 or more users matching the filter
func (s *Store) Users(
	ctx context.Context,
	exec boil.ContextExecutor,
	filter users.UserQueryFilter,
) ([]*users.User, error) {
	mods := []qm.QueryMod{}

	// IDs filter
	if len(filter.IDs) > 0 {
		// Convert string IDs to uint64 for current DB schema
		ids := make([]interface{}, len(filter.IDs))
		for i, id := range filter.IDs {
			idNum, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid user ID %s: %w", id, err)
			}
			ids[i] = idNum
		}
		mods = append(mods, qm.WhereIn("id IN ?", ids...))
	}

	// Username filter
	if filter.Username.Valid {
		mods = append(mods, qm.Where("username = ?", filter.Username.String))
	}

	// Email filter
	if filter.Email.Valid {
		mods = append(mods, qm.Where("email = ?", filter.Email.String))
	}

	// Gender filter
	if filter.Gender.Valid {
		mods = append(mods, qm.Where("gender = ?", string(filter.Gender.String)))
	}

	// Sorting
	if filter.OrderBy.Valid {
		sortDir := "ASC"
		if filter.Sort.Valid {
			sortDir = filter.Sort.String
		}
		mods = append(mods, qm.OrderBy(filter.OrderBy.String+" "+sortDir))
	}

	// Pagination
	if filter.Limit.Valid {
		mods = append(mods, qm.Limit(filter.Limit.Int))
	}
	if filter.Offset.Valid {
		mods = append(mods, qm.Offset(filter.Offset.Int))
	}

	// Execute query
	dbUsers, err := models.Users(mods...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	return dbUsersToUsers(dbUsers), nil
}

// User returns exactly 1 user, errors if 0 or >1 found
func (s *Store) User(
	ctx context.Context,
	exec boil.ContextExecutor,
	filter users.UserQueryFilter,
) (*users.User, error) {
	results, err := s.Users(ctx, exec, filter)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no user found")
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", len(results))
	}

	return results[0], nil
}

// Update performs generic update with nullable fields
func (s *Store) Update(
	ctx context.Context,
	exec boil.ContextExecutor,
	update users.UpdateUser,
) error {
	if len(update.IDs) == 0 {
		return fmt.Errorf("no user IDs provided")
	}

	// Build update columns
	cols := make(map[string]interface{})

	if update.Username.Valid {
		cols["username"] = update.Username.String
	}
	if update.DisplayName.Valid {
		cols["display_name"] = update.DisplayName.String
	}
	if update.Gender.Valid {
		cols["gender"] = string(update.Gender.String)
	}

	if len(cols) == 0 {
		return nil // Nothing to update
	}

	// Convert string IDs to uint64
	ids := make([]interface{}, len(update.IDs))
	for i, id := range update.IDs {
		idNum, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user ID %s: %w", id, err)
		}
		ids[i] = idNum
	}

	// Execute update
	_, err := models.Users(
		qm.WhereIn("id IN ?", ids...),
	).UpdateAll(ctx, exec, cols)

	return err
}

// dbUsersToUsers converts DB models to domain models
func dbUsersToUsers(dbUsers []*models.User) []*users.User {
	result := make([]*users.User, len(dbUsers))
	for i, db := range dbUsers {
		displayName := ""
		if db.DisplayName.Valid {
			displayName = db.DisplayName.String
		}

		gender := users.Gender("")
		if db.Gender.Valid {
			gender = users.Gender(db.Gender.String)
		}

		var createdAt time.Time
		if db.CreatedAt.Valid {
			createdAt = db.CreatedAt.Time
		}

		result[i] = &users.User{
			ID:          fmt.Sprintf("%d", db.ID),
			Username:    db.Username,
			Email:       "", // TODO: Add after migration
			DisplayName: displayName,
			Gender:      gender,
			CreatedAt:   createdAt,
		}
	}
	return result
}

// userToDBUser converts domain model to DB model
func userToDBUser(user *users.User) (*models.User, error) {
	id, err := strconv.ParseUint(user.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return &models.User{
		ID:       id,
		Username: user.Username,
		DisplayName: null.String{
			String: user.DisplayName,
			Valid:  user.DisplayName != "",
		},
		Gender: null.String{
			String: string(user.Gender),
			Valid:  user.Gender != "",
		},
		CreatedAt: null.Time{
			Time:  user.CreatedAt,
			Valid: !user.CreatedAt.IsZero(),
		},
	}, nil
}
