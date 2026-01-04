package factory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"

	"mlm/models"
)

// UserMods - optional overrides for user creation
type UserMods struct {
	ID          *uint64
	Username    string
	Email       string
	DisplayName string
	Gender      string
}

// User creates a test user with optional overrides
func User(
	t *testing.T,
	exec boil.ContextExecutor,
	mods *UserMods,
) *models.User {
	if mods == nil {
		mods = &UserMods{}
	}

	// Generate unique defaults
	timestamp := time.Now().UnixNano()

	// Username
	if mods.Username == "" {
		mods.Username = fmt.Sprintf("user_%d", timestamp)
	}

	// Email
	if mods.Email == "" {
		mods.Email = fmt.Sprintf("%s@example.com", mods.Username)
	}

	// Display name
	if mods.DisplayName == "" {
		mods.DisplayName = fmt.Sprintf("User %d", timestamp)
	}

	// Gender
	if mods.Gender == "" {
		mods.Gender = "male"
	}

	user := &models.User{
		Username:    mods.Username,
		DisplayName: null.StringFrom(mods.DisplayName),
		Gender:      null.StringFrom(mods.Gender),
		CreatedAt:   null.TimeFrom(time.Now()),
	}

	// If ID is provided, set it (for specific test cases)
	if mods.ID != nil {
		user.ID = *mods.ID
	}

	err := user.Insert(context.Background(), exec, boil.Infer())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	return user
}

// Users creates multiple test users
func Users(
	t *testing.T,
	exec boil.ContextExecutor,
	count int,
	baseMods *UserMods,
) []*models.User {
	users := make([]*models.User, count)

	for i := 0; i < count; i++ {
		var mods *UserMods
		if baseMods != nil {
			// Copy base mods and make unique
			mods = &UserMods{
				Gender:      baseMods.Gender,
				DisplayName: baseMods.DisplayName,
			}
			if baseMods.Username != "" {
				mods.Username = fmt.Sprintf("%s_%d", baseMods.Username, i)
			}
		}

		users[i] = User(t, exec, mods)
		// Small delay to ensure unique timestamps
		time.Sleep(1 * time.Millisecond)
	}

	return users
}
