# Tonight's Work - Store Layer, Migrations & Tests

## Goal
✅ Create store layer functions
✅ Create migrations
✅ Create and run tests

---

## ✅ What You Already Have

### 1. Users Domain (100% Complete)
- ✅ Models: `internal/musicapp/lib/users/models.go`
- ✅ Store: `internal/musicapp/lib/users/store/users.go`
- ✅ Tests: `internal/musicapp/lib/users/store/users_test.go` (25 tests)
- ✅ Factory: `internal/musicapp/db/factory/user.go`
- ✅ Repo: `internal/musicapp/db/repo/user.go`
- ✅ Logic: `internal/musicapp/lib/users/logic.go`

### 2. Migrations
- ✅ `01_create_users.up.sql`
- ✅ `02_create_rooms.up.sql`
- ✅ `03_create_room_members.up.sql`
- ✅ `04_update_users_for_imapp.up.sql`

### 3. Test Infrastructure
- ✅ `internal/testsuite/testsuite.go`
- ✅ `setup_test_db.sh`
- ✅ `run_tests.sh`

---

## 🎯 Tonight's Tasks

### Task 1: Setup Test Database (5 mins)

```bash
# Option 1: Use the script
./setup_test_db.sh

# Option 2: Manual
mysql -u root -p <<EOF
CREATE DATABASE mlm_test;
GRANT ALL PRIVILEGES ON mlm_test.* TO 'user'@'localhost';
FLUSH PRIVILEGES;
EOF

# Run migrations on test DB
mysql -u user -ppassword mlm_test < migration/01_create_users.up.sql
mysql -u user -ppassword mlm_test < migration/02_create_rooms.up.sql
mysql -u user -ppassword mlm_test < migration/03_create_room_members.up.sql
mysql -u user -ppassword mlm_test < migration/04_update_users_for_imapp.up.sql
```

### Task 2: Run Existing Tests (2 mins)

```bash
# Quick run
./run_tests.sh

# Or manually
go test -v ./internal/musicapp/lib/users/store/...
```

**Expected:** All 25 tests should PASS ✅

### Task 3: Create Rooms Store (30 mins)

Following the **exact same pattern** as users:

#### 3a. Create Models

**File:** `internal/musicapp/lib/rooms/models.go`

```go
package rooms

import (
    "time"
    "github.com/aarondl/null/v8"
)

type Room struct {
    ID          string
    Name        string
    CreatedBy   string
    IsActive    bool
    CreatedAt   time.Time
}

type RoomQueryFilter struct {
    IDs        []string
    Name       null.Val[string]
    CreatedBy  null.Val[string]
    IsActive   null.Val[bool]
    OrderBy    null.Val[string]
    Sort       null.Val[string]
    Limit      null.Val[int]
    Offset     null.Val[int]
}

type UpdateRoom struct {
    IDs      []string
    Name     null.Val[string]
    IsActive null.Val[bool]
}
```

#### 3b. Create Store

**File:** `internal/musicapp/lib/rooms/store/rooms.go`

Copy the structure from `users/store/users.go` and adapt:

```go
package store

import (
    "context"
    "fmt"
    // ... imports
    "mlm/internal/musicapp/lib/rooms"
    "mlm/models"
)

type Store struct{}

func New() *Store {
    return &Store{}
}

// Rooms - returns 0 or more
func (s *Store) Rooms(ctx, exec, filter) ([]*rooms.Room, error) {
    // Build query with filters
    // Execute
    // Convert and return
}

// Room - returns exactly 1
func (s *Store) Room(ctx, exec, filter) (*rooms.Room, error) {
    results, err := s.Rooms(ctx, exec, filter)
    if len(results) == 0 {
        return nil, fmt.Errorf("no room found")
    }
    if len(results) > 1 {
        return nil, fmt.Errorf("expected 1 room, got %d", len(results))
    }
    return results[0], nil
}

// Update - generic update
func (s *Store) Update(ctx, exec, update) error {
    // Build update
    // Execute
}
```

#### 3c. Create Factory

**File:** `internal/musicapp/db/factory/room.go`

```go
package factory

import (
    "testing"
    "mlm/models"
    // ...
)

type RoomMods struct {
    ID        *uint64
    Name      string
    CreatedBy uint64
    IsActive  bool
}

func Room(t *testing.T, exec boil.ContextExecutor, mods *RoomMods) *models.Room {
    if mods == nil {
        mods = &RoomMods{}
    }

    // Set defaults
    if mods.Name == "" {
        mods.Name = fmt.Sprintf("room_%d", time.Now().UnixNano())
    }

    if mods.CreatedBy == 0 {
        // Auto-create a user
        user := User(t, exec, nil)
        mods.CreatedBy = user.ID
    }

    room := &models.Room{
        Name:      mods.Name,
        CreatedBy: mods.CreatedBy,
        IsActive:  mods.IsActive,
        CreatedAt: time.Now(),
    }

    err := room.Insert(context.Background(), exec, boil.Infer())
    if err != nil {
        t.Fatalf("failed to create room: %v", err)
    }

    return room
}
```

#### 3d. Create Tests

**File:** `internal/musicapp/lib/rooms/store/rooms_test.go`

```go
package store_test

import (
    "testing"
    "github.com/aarondl/null/v8"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "mlm/internal/musicapp/db/factory"
    "mlm/internal/musicapp/lib/rooms"
    "mlm/internal/musicapp/lib/rooms/store"
    "mlm/internal/testsuite"
)

type testCaseRooms struct {
    name            string
    setup           func(th *testsuite.Helper) rooms.RoomQueryFilter
    extraAssertions func(th *testsuite.Helper, result []*rooms.Room, err error)
}

func roomsTestCases() []testCaseRooms {
    return []testCaseRooms{
        {
            name: "success-returns-all-rooms",
            setup: func(th *testsuite.Helper) rooms.RoomQueryFilter {
                factory.Room(th.T, th.BackendAppDb(), nil)
                factory.Room(th.T, th.BackendAppDb(), nil)

                return rooms.RoomQueryFilter{}
            },
            extraAssertions: func(th *testsuite.Helper, result []*rooms.Room, err error) {
                require.NoError(th.T, err)
                assert.Len(th.T, result, 2)
            },
        },
        {
            name: "success-filters-by-active",
            setup: func(th *testsuite.Helper) rooms.RoomQueryFilter {
                factory.Room(th.T, th.BackendAppDb(), &factory.RoomMods{
                    IsActive: true,
                })
                factory.Room(th.T, th.BackendAppDb(), &factory.RoomMods{
                    IsActive: false,
                })

                return rooms.RoomQueryFilter{
                    IsActive: null.From(true),
                }
            },
            extraAssertions: func(th *testsuite.Helper, result []*rooms.Room, err error) {
                require.NoError(th.T, err)
                assert.Len(th.T, result, 1)
                assert.True(th.T, result[0].IsActive)
            },
        },
        // Add more test cases...
    }
}

func TestStore_Rooms(t *testing.T) {
    for _, tt := range roomsTestCases() {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            testSuite := testsuite.New(t)
            t.Cleanup(testSuite.UseBackendDB())

            store := store.New()
            filter := tt.setup(testSuite)

            result, err := store.Rooms(
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
```

### Task 4: Run Rooms Tests (5 mins)

```bash
go test -v ./internal/musicapp/lib/rooms/store/...
```

### Task 5: Add More Migrations (Optional, 15 mins)

If you want to add artists, genres, etc:

**File:** `migration/05_create_artists.up.sql`

```sql
CREATE TABLE artists (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    spotify_artist_id VARCHAR(100),
    genre_id BIGINT UNSIGNED,
    image_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Then repeat the pattern:
1. Models
2. Store
3. Factory
4. Tests

---

## 📊 Progress Tracker

```
✅ Users Domain
   ✅ Models
   ✅ Store
   ✅ Tests (25 cases)
   ✅ Factory
   ✅ Repo
   ✅ Logic

⬜ Rooms Domain
   ⬜ Models
   ⬜ Store
   ⬜ Tests (10+ cases)
   ⬜ Factory

⬜ Artists Domain (optional tonight)
   ⬜ Migration
   ⬜ Models
   ⬜ Store
   ⬜ Tests
```

---

## 🎯 Success Criteria for Tonight

Minimum:
- ✅ User tests run and pass
- ✅ Rooms store implemented
- ✅ Rooms tests written (at least 10 test cases)
- ✅ All tests passing

Stretch:
- ✅ Artists domain started
- ✅ Playlists domain started

---

## 💡 Tips

### Copy-Paste Pattern

1. **Copy users → rooms:**
   ```bash
   cp -r internal/musicapp/lib/users internal/musicapp/lib/rooms
   # Then find-replace "user" → "room" in all files
   ```

2. **Use factory pattern:**
   - Always create FK dependencies automatically
   - Use unique names with timestamps

3. **Test patterns:**
   - success-* for happy paths
   - error-* for failure cases
   - Test filters, sorting, pagination

### Common Test Cases

For every store, test:
- ✅ Returns all
- ✅ Filters by each field
- ✅ Sorting ASC/DESC
- ✅ Pagination (limit/offset)
- ✅ Filter by IDs
- ✅ No results found
- ✅ Single entity (error if 0 or >1)
- ✅ Update operations

---

## ⚡ Quick Commands

```bash
# Setup
./setup_test_db.sh

# Run all tests
go test -v ./internal/musicapp/lib/.../store/...

# Run specific domain
go test -v ./internal/musicapp/lib/users/store/...
go test -v ./internal/musicapp/lib/rooms/store/...

# With coverage
go test -v -cover ./internal/musicapp/lib/.../store/...

# With race detection
go test -v -race ./internal/musicapp/lib/.../store/...
```

---

## 📝 Test Template (Save This)

```go
{
    name: "success-your-test-name",
    setup: func(th *testsuite.Helper) YourFilter {
        // Create test data
        entity := factory.YourEntity(th.T, th.BackendAppDb(), &factory.YourMods{
            Field: value,
        })

        return YourFilter{
            Field: null.From(value),
        }
    },
    extraAssertions: func(th *testsuite.Helper, result []*YourModel, err error) {
        require.NoError(th.T, err)
        assert.Len(th.T, result, 1)
        assert.Equal(th.T, expectedValue, result[0].Field)
    },
},
```

---

## 🎉 End Goal for Tonight

By end of tonight you should have:

1. ✅ All user tests passing (25 tests)
2. ✅ Rooms store fully implemented
3. ✅ Rooms tests written (10-15 tests)
4. ✅ All tests green

**Total test count goal:** 35-40 passing tests

---

**You already have 25 tests ready! Just need to set up the DB and run them.** 🚀

Good luck tonight!
