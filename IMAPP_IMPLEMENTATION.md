# IMAPP Pattern Implementation - Users Domain

## What Was Built

Complete implementation of the **Users domain** following IMAPP architecture pattern.

### Files Created

```
internal/
├── musicapp/
│   ├── lib/
│   │   └── users/
│   │       ├── models.go         ✅ Domain models + QueryFilter
│   │       ├── logic.go          ✅ Business logic orchestration
│   │       └── store/
│   │           ├── users.go      ✅ Generic store queries
│   │           └── users_test.go ✅ Table-driven tests (25 test cases)
│   └── db/
│       ├── factory/
│       │   └── user.go           ✅ Test data factory
│       └── repo/
│           └── user.go           ✅ Insert/BulkInsert operations
└── testsuite/
    └── testsuite.go              ✅ Test infrastructure

migration/
└── 04_update_users_for_imapp.up.sql ✅ Add email column
```

---

## Architecture Flow

```
API Layer (TODO)
     ↓
Business Layer (TODO)
     ↓
Logic Layer (users/logic.go)
     ↓
Store Layer (users/store/users.go)
     ↓
Repo Layer (db/repo/user.go)
     ↓
Database
```

---

## Key Components

### 1. Domain Models (`lib/users/models.go`)

Clean domain models with no DB tags:

```go
type User struct {
    ID          string
    Username    string
    Email       string
    DisplayName string
    Gender      Gender
    CreatedAt   time.Time
}

type UserQueryFilter struct {
    IDs      []string
    Username null.Val[string]
    Email    null.Val[string]
    Gender   null.Val[Gender]
    OrderBy  null.Val[string]
    Sort     null.Val[string]
    Limit    null.Val[int]
    Offset   null.Val[int]
}
```

### 2. Store Layer (`lib/users/store/users.go`)

**GENERIC VERBS ONLY** - No domain-specific queries:

```go
// ✅ Good - Generic
store.Users(ctx, exec, UserQueryFilter{Gender: null.From(GenderMale)})
store.User(ctx, exec, UserQueryFilter{IDs: []string{"123"}})
store.Update(ctx, exec, UpdateUser{...})

// ❌ Bad - Domain-specific
store.GetActiveUsersByGender(...)
store.FindMaleUsers(...)
```

**Methods:**
- `Users(filter)` → Returns 0 or more users
- `User(filter)` → Returns exactly 1 user (errors if 0 or >1)
- `Update(update)` → Generic update with nullable fields

### 3. Logic Layer (`lib/users/logic.go`)

Composes store calls with business rules:

```go
// Domain-specific methods built from generic store
logic.GetUsersByGender(ctx, exec, GenderMale)
logic.GetRecentUsers(ctx, exec, 10)
logic.SearchUsers(ctx, exec, &gender, &username, page, pageSize)
logic.CountUsersByGender(ctx, exec)
```

### 4. Repository Layer (`db/repo/user.go`)

Insert/Update operations returning pgmodel types:

```go
repo.Insert(ctx, exec, user)
repo.BulkInsert(ctx, exec, users)  // NO LOOPS!
repo.Upsert(ctx, exec, user)
```

### 5. Factory (`db/factory/user.go`)

Test data creation:

```go
// Create single user
user := factory.User(t, db, nil)

// Create with overrides
user := factory.User(t, db, &factory.UserMods{
    Username: "alice",
    Gender:   "female",
})

// Create multiple users
users := factory.Users(t, db, 5, nil)
```

### 6. Tests (`lib/users/store/users_test.go`)

**25 comprehensive test cases** covering:
- ✅ Filtering (gender, username, IDs)
- ✅ Sorting (ASC/DESC, multiple fields)
- ✅ Pagination (limit, offset)
- ✅ Single user queries
- ✅ Update operations
- ✅ Error cases

**Table-Driven Pattern:**
```go
func usersTestCases() []testCaseUsers {
    return []testCaseUsers{
        {
            name: "success-filters-by-gender-male",
            setup: func(th *testsuite.Helper) UserQueryFilter {
                factory.User(th.T, th.BackendAppDb(), &factory.UserMods{
                    Gender: "male",
                })
                return UserQueryFilter{Gender: null.From(GenderMale)}
            },
            extraAssertions: func(th *testsuite.Helper, result []*users.User, err error) {
                require.NoError(th.T, err)
                assert.Len(th.T, result, 1)
            },
        },
    }
}
```

---

## How to Use

### Example 1: Query Users

```go
import (
    "mlm/internal/musicapp/lib/users"
    "mlm/internal/musicapp/lib/users/store"
)

// Create store
userStore := store.New()

// Query all male users
maleUsers, err := userStore.Users(ctx, db, users.UserQueryFilter{
    Gender: null.From(users.GenderMale),
})

// Query with pagination
pagedUsers, err := userStore.Users(ctx, db, users.UserQueryFilter{
    Gender:  null.From(users.GenderFemale),
    OrderBy: null.From("created_at"),
    Sort:    null.From("DESC"),
    Limit:   null.From(10),
    Offset:  null.From(0),
})

// Get single user by ID
user, err := userStore.User(ctx, db, users.UserQueryFilter{
    IDs: []string{"123"},
})
```

### Example 2: Update Users

```go
// Update username
err := userStore.Update(ctx, db, users.UpdateUser{
    IDs:      []string{"123"},
    Username: null.From("newusername"),
})

// Update multiple fields
err := userStore.Update(ctx, db, users.UpdateUser{
    IDs:         []string{"123"},
    Username:    null.From("alice"),
    DisplayName: null.From("Alice Smith"),
    Gender:      null.From(users.GenderFemale),
})

// Update multiple users at once
err := userStore.Update(ctx, db, users.UpdateUser{
    IDs:    []string{"123", "456", "789"},
    Gender: null.From(users.GenderOther),
})
```

### Example 3: Use Logic Layer

```go
import (
    "mlm/internal/musicapp/lib/users"
    "mlm/internal/musicapp/lib/users/store"
)

// Create dependencies
userStore := store.New()
userLogic, err := users.NewLogic(userStore)

// Use domain-specific methods
maleUsers, err := userLogic.GetMaleUsers(ctx, db)
recentUsers, err := userLogic.GetRecentUsers(ctx, db, 10)
user, err := userLogic.GetUserByID(ctx, db, "123")

// Complex searches
results, err := userLogic.SearchUsers(
    ctx, db,
    &gender,    // optional
    &username,  // optional
    0,          // page
    20,         // pageSize
)
```

---

## Setup Instructions

### 1. Install Dependencies

```bash
go get github.com/aarondl/null/v8
go get github.com/volatiletech/sqlboiler/v4
go get github.com/stretchr/testify
```

### 2. Create Test Database

```bash
mysql -u root -p
```

```sql
CREATE DATABASE mlm_test;
GRANT ALL PRIVILEGES ON mlm_test.* TO 'user'@'localhost' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;
```

### 3. Run Migrations

```bash
# Production DB
mysql -u user -p mlm < migration/01_create_users.up.sql
mysql -u user -p mlm < migration/04_update_users_for_imapp.up.sql

# Test DB
mysql -u user -p mlm_test < migration/01_create_users.up.sql
mysql -u user -p mlm_test < migration/04_update_users_for_imapp.up.sql
```

### 4. Run Tests

```bash
# Run all user store tests
go test -v -count=1 ./internal/musicapp/lib/users/store/...

# Run with race detection
go test -v -count=1 -race ./internal/musicapp/lib/users/store/...

# Run specific test
go test -v -run TestStore_Users ./internal/musicapp/lib/users/store/...
```

---

## Testing Output Example

```
=== RUN   TestStore_Users
=== RUN   TestStore_Users/success-returns-all-users
=== PAUSE TestStore_Users/success-returns-all-users
=== RUN   TestStore_Users/success-filters-by-gender-male
=== PAUSE TestStore_Users/success-filters-by-gender-male
...
=== CONT  TestStore_Users/success-returns-all-users
=== CONT  TestStore_Users/success-filters-by-gender-male
--- PASS: TestStore_Users (0.15s)
    --- PASS: TestStore_Users/success-returns-all-users (0.05s)
    --- PASS: TestStore_Users/success-filters-by-gender-male (0.04s)
    ...
PASS
```

---

## IMAPP Principles Followed

| Principle | Implementation |
|-----------|----------------|
| ✅ Store = Generic Verbs | `Users()`, `User()`, `Update()` only |
| ✅ Logic = Composition | `GetUsersByGender()` composes `Users(filter)` |
| ✅ Fail Fast in Constructor | `NewLogic(store)` validates dependencies |
| ✅ No Nil Returns | Never return `nil, nil` |
| ✅ Bulk Insert = Raw SQL | `BulkInsert()` uses single query |
| ✅ Table-Driven Tests | 25 test cases in table structure |
| ✅ Factory Pattern | `factory.User()` for test data |
| ✅ QueryFilter Pattern | `null.Val[T]` for optional filters |

---

## Next Steps

### Phase 1: Complete Users Domain ✅ DONE
- [x] Models
- [x] Store
- [x] Tests
- [x] Factory
- [x] Repo
- [x] Logic

### Phase 2: Add Rooms Domain
Following same pattern:
1. `lib/rooms/models.go`
2. `lib/rooms/store/rooms.go`
3. `lib/rooms/store/rooms_test.go`
4. `db/factory/room.go`
5. `db/repo/room.go`
6. `lib/rooms/logic.go`

### Phase 3: Add API Layer
1. `internal/musicapp/api/users.go` - HTTP handlers
2. `internal/musicapp/business/users.go` - Business orchestration
3. Wire with dependency injection

### Phase 4: Migrate Old Code
1. Replace `handlers/users.go` with new API layer
2. Replace `store/users_store.go` with new store layer
3. Remove old files

---

## Comparison: Old vs New

### Old Structure
```
handlers/users.go        → HTTP + Business logic mixed
store/users_store.go     → Specific methods (GetUserByID, CreateUser)
models/Users.go          → Simple struct
NO TESTS
```

### New Structure (IMAPP)
```
api/users.go             → HTTP handling only
business/users.go        → Business orchestration
lib/users/logic.go       → Domain logic composition
lib/users/store/users.go → Generic queries only
lib/users/models.go      → Clean domain models
db/repo/user.go          → Insert/Update operations
db/factory/user.go       → Test data creation
25 COMPREHENSIVE TESTS ✅
```

---

## Benefits

1. **Testable** - Every layer has clear boundaries and can be mocked
2. **Maintainable** - Changes in one layer don't cascade
3. **Scalable** - Easy to add new domains following same pattern
4. **Type-Safe** - `null.Val[T]` prevents nil pointer errors
5. **Performant** - Bulk operations avoid N+1 queries
6. **Documented** - Tests serve as documentation

---

## Common Patterns

### Pattern 1: Complex Search
```go
// DON'T create store method: GetActiveFemalUsersOver21InRoom()
// DO compose filter:
store.Users(ctx, exec, UserQueryFilter{
    Gender:   null.From(GenderFemale),
    IsActive: null.From(true),
    // Add more filters as needed
})
```

### Pattern 2: Aggregation
```go
// Aggregation happens in LOGIC layer, not store
func (l *Logic) CountByGender(ctx, exec) (map[Gender]int, error) {
    users, err := l.store.Users(ctx, exec, UserQueryFilter{})
    // Count in memory
    counts := make(map[Gender]int)
    for _, u := range users {
        counts[u.Gender]++
    }
    return counts, nil
}
```

### Pattern 3: Transactions
```go
tx, err := db.BeginTx(ctx, nil)
defer tx.Rollback()

// Use tx instead of db
users, err := store.Users(ctx, tx, filter)
err = store.Update(ctx, tx, update)

tx.Commit()
```

---

**Users domain is complete! Ready to implement Rooms, Artists, Playlists, etc.**
