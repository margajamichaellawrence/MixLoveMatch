package store_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mlm/internal/musicapp/db/factory"
	"mlm/internal/musicapp/lib/users"
	"mlm/internal/musicapp/lib/users/store"
	"mlm/internal/testsuite"
)

// Test case struct for Users()
type testCaseUsers struct {
	name            string
	setup           func(th *testsuite.Helper) users.UserQueryFilter
	extraAssertions func(th *testsuite.Helper, result []*users.User, err error)
}

// Test cases for Users() method
func usersTestCases() []testCaseUsers {
	return []testCaseUsers{
		{
			name: "success-returns-all-users",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				// Create 2 test users
				factory.User(th.T, th.BackendAppDb(), nil)
				factory.User(th.T, th.BackendAppDb(), nil)

				return users.UserQueryFilter{}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
			},
		},
		{
			name: "success-filters-by-gender-male",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				// Create male user
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Gender: "male",
				})
				// Create female user
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Gender: "female",
				})

				return users.UserQueryFilter{
					Gender: null.StringFrom(string(users.GenderMale)),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 1)
				assert.Equal(th.T, users.GenderMale, result[0].Gender)
			},
		},
		{
			name: "success-filters-by-gender-female",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				// Create male user
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Gender: "male",
				})
				// Create 2 female users
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Gender: "female",
				})
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Gender: "female",
				})

				return users.UserQueryFilter{
					Gender: null.StringFrom(string(users.GenderFemale)),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
				for _, u := range result {
					assert.Equal(th.T, users.GenderFemale, u.Gender)
				}
			},
		},
		{
			name: "success-filters-by-username",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "alice",
				})
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "bob",
				})

				return users.UserQueryFilter{
					Username: null.StringFrom("alice"),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 1)
				assert.Equal(th.T, "alice", result[0].Username)
			},
		},
		{
			name: "success-sorts-by-created-at-desc",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "first",
				})
				time.Sleep(1 * time.Second)
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "second",
				})

				return users.UserQueryFilter{
					OrderBy: null.StringFrom("created_at"),
					Sort:    null.StringFrom("DESC"),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
				// First result should be newer (second user)
				assert.Equal(th.T, "second", result[0].Username)
				assert.Equal(th.T, "first", result[1].Username)
			},
		},
		{
			name: "success-sorts-by-username-asc",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "charlie",
				})
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "alice",
				})
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "bob",
				})

				return users.UserQueryFilter{
					OrderBy: null.StringFrom("username"),
					Sort:    null.StringFrom("ASC"),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 3)
				assert.Equal(th.T, "alice", result[0].Username)
				assert.Equal(th.T, "bob", result[1].Username)
				assert.Equal(th.T, "charlie", result[2].Username)
			},
		},
		{
			name: "success-limits-results",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				factory.Users(th.T, th.BackendAppDb(), 5, nil)

				return users.UserQueryFilter{
					Limit: null.IntFrom(2),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
			},
		},
		{
			name: "success-pagination-with-offset",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				factory.Users(th.T, th.BackendAppDb(), 5, nil)

				return users.UserQueryFilter{
					OrderBy: null.StringFrom("id"),
					Sort:    null.StringFrom("ASC"),
					Limit:   null.IntFrom(2),
					Offset:  null.IntFrom(2),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
				// Should skip first 2 and return next 2
			},
		},
		{
			name: "success-filters-by-ids",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				u1 := factory.User(th.T, th.BackendAppDb(), nil)
				u2 := factory.User(th.T, th.BackendAppDb(), nil)
				factory.User(th.T, th.BackendAppDb(), nil) // Extra user not in filter

				return users.UserQueryFilter{
					IDs: []string{
						fmt.Sprintf("%d", u1.ID),
						fmt.Sprintf("%d", u2.ID),
					},
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
			},
		},
		{
			name: "success-no-users-found",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				// No users created
				return users.UserQueryFilter{
					Username: null.StringFrom("nonexistent"),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 0)
			},
		},
		{
			name: "success-combines-multiple-filters",
			setup: func(th *testsuite.Helper) users.UserQueryFilter {
				// Create male users
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "alice",
					Gender:   "male",
				})
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "bob",
					Gender:   "male",
				})
				// Create female user
				factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
					Username: "charlie",
					Gender:   "female",
				})

				return users.UserQueryFilter{
					Gender:  null.StringFrom(string(users.GenderMale)),
					OrderBy: null.StringFrom("username"),
					Sort:    null.StringFrom("ASC"),
					Limit:   null.IntFrom(10),
				}
			},
			extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
				require.NoError(th.T, err)
				assert.Len(th.T, result, 2)
				assert.Equal(th.T, "alice", result[0].Username)
				assert.Equal(th.T, "bob", result[1].Username)
			},
		},
	}
}

// TestStore_Users - main test function
func TestStore_Users(t *testing.T) {

	for _, tt := range usersTestCases() {
		t.Run(tt.name, func(t *testing.T) {

			testSuite := testsuite.New(t)
			t.Cleanup(testSuite.UseBackendDB())

			store := store.New()
			filter := tt.setup(testSuite)

			result, err := store.Users(
				testSuite.Ctx,
				testSuite.BackendAppDb(),
				filter,
			)

			if tt.extraAssertions != nil {
				tt.extraAssertions(testSuite, result, err)
			}
		})
	}
}

// TestStore_User - test User() method (singular)
func TestStore_User(t *testing.T) {
	t.Run("success-returns-single-user", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		dbUser := factory.User(testSuite.T, testSuite.BackendAppDb(), &factory.UserMods{
			Username: "testuser",
		})

		store := store.New()
		result, err := store.User(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				IDs: []string{fmt.Sprintf("%d", dbUser.ID)},
			},
		)

		require.NoError(testSuite.T, err)
		assert.Equal(testSuite.T, fmt.Sprintf("%d", dbUser.ID), result.ID)
		assert.Equal(testSuite.T, "testuser", result.Username)
	})

	t.Run("error-no-user-found", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		store := store.New()
		_, err := store.User(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				Username: null.StringFrom("nonexistent"),
			},
		)

		require.Error(testSuite.T, err)
		assert.Contains(testSuite.T, err.Error(), "no user found")
	})

	t.Run("error-multiple-users-found", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		// Create 2 users with same gender
		factory.User(testSuite.T, testSuite.BackendAppDb(), &factory.UserMods{
			Gender: "male",
		})
		factory.User(testSuite.T, testSuite.BackendAppDb(), &factory.UserMods{
			Gender: "male",
		})

		store := store.New()
		_, err := store.User(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				Gender: null.StringFrom(string(users.GenderMale)),
			},
		)

		require.Error(testSuite.T, err)
		assert.Contains(testSuite.T, err.Error(), "expected 1 user, got 2")
	})
}

// TestStore_Update - test Update() method
func TestStore_Update(t *testing.T) {
	t.Run("success-updates-username", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		dbUser := factory.User(testSuite.T, testSuite.BackendAppDb(), &factory.UserMods{
			Username: "oldname",
		})

		store := store.New()
		err := store.Update(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UpdateUser{
				IDs:      []string{fmt.Sprintf("%d", dbUser.ID)},
				Username: null.StringFrom("newname"),
			},
		)

		require.NoError(testSuite.T, err)

		// Verify update
		updated, err := store.User(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				IDs: []string{fmt.Sprintf("%d", dbUser.ID)},
			},
		)

		require.NoError(testSuite.T, err)
		assert.Equal(testSuite.T, "newname", updated.Username)
	})

	t.Run("success-updates-multiple-fields", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		dbUser := factory.User(testSuite.T, testSuite.BackendAppDb(), &factory.UserMods{
			Username:    "oldname",
			Gender:      "male",
			DisplayName: "Old Display",
		})

		store := store.New()
		err := store.Update(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UpdateUser{
				IDs:         []string{fmt.Sprintf("%d", dbUser.ID)},
				Username:    null.StringFrom("newname"),
				Gender:      null.StringFrom(string(users.GenderFemale)),
				DisplayName: null.StringFrom("New Display"),
			},
		)

		require.NoError(testSuite.T, err)

		// Verify all updates
		updated, err := store.User(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				IDs: []string{fmt.Sprintf("%d", dbUser.ID)},
			},
		)

		require.NoError(testSuite.T, err)
		assert.Equal(testSuite.T, "newname", updated.Username)
		assert.Equal(testSuite.T, users.GenderFemale, updated.Gender)
		assert.Equal(testSuite.T, "New Display", updated.DisplayName)
	})

	t.Run("success-updates-multiple-users", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		u1 := factory.User(testSuite.T, testSuite.BackendAppDb(), nil)
		u2 := factory.User(testSuite.T, testSuite.BackendAppDb(), nil)

		store := store.New()
		err := store.Update(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UpdateUser{
				IDs: []string{
					fmt.Sprintf("%d", u1.ID),
					fmt.Sprintf("%d", u2.ID),
				},
				Gender: null.StringFrom(string(users.GenderOther)),
			},
		)

		require.NoError(testSuite.T, err)

		// Verify both updated
		updated, err := store.Users(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UserQueryFilter{
				Gender: null.StringFrom(string(users.GenderOther)),
			},
		)

		require.NoError(testSuite.T, err)
		assert.Len(testSuite.T, updated, 2)
	})

	t.Run("error-no-ids-provided", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		store := store.New()
		err := store.Update(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UpdateUser{
				IDs:      []string{},
				Username: null.StringFrom("newname"),
			},
		)

		require.Error(testSuite.T, err)
		assert.Contains(testSuite.T, err.Error(), "no user IDs provided")
	})

	t.Run("success-nothing-to-update", func(t *testing.T) {

		testSuite := testsuite.New(t)
		t.Cleanup(testSuite.UseBackendDB())

		dbUser := factory.User(testSuite.T, testSuite.BackendAppDb(), nil)

		store := store.New()
		err := store.Update(
			testSuite.Ctx,
			testSuite.BackendAppDb(),
			users.UpdateUser{
				IDs: []string{fmt.Sprintf("%d", dbUser.ID)},
				// No fields to update
			},
		)

		require.NoError(testSuite.T, err) // Should succeed but do nothing
	})
}
