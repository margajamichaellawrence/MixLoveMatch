# How to Run Store Layer Tests

## Quick Start

### 1. Create Test Database

```bash
# Connect to MySQL
mysql -u root -p

# Create test database
CREATE DATABASE mlm_test;
GRANT ALL PRIVILEGES ON mlm_test.* TO 'user'@'localhost';
FLUSH PRIVILEGES;
exit;
```

**Or use the setup script:**
```bash
./setup_test_db.sh
```

### 2. Run Migrations on Test DB

```bash
# Export test database env var
export MUSICAPP_PG_DATABASE=mlm_test

# Run migrations
mysql -u user -p mlm_test < migration/01_create_users.up.sql
mysql -u user -p mlm_test < migration/02_create_rooms.up.sql
mysql -u user -p mlm_test < migration/03_create_room_members.up.sql
mysql -u user -p mlm_test < migration/04_update_users_for_imapp.up.sql
```

### 3. Run the Tests

```bash
# Run all user store tests
go test -v ./internal/musicapp/lib/users/store/...

# Run with race detection
go test -v -race ./internal/musicapp/lib/users/store/...

# Run specific test
go test -v -run TestStore_Users/success-filters-by-gender ./internal/musicapp/lib/users/store/...

# Run with coverage
go test -v -cover ./internal/musicapp/lib/users/store/...
```

---

## What Tests Exist?

### Users Store Tests (25 test cases)

**File:** `internal/musicapp/lib/users/store/users_test.go`

#### TestStore_Users (12 cases)
- ✅ success-returns-all-users
- ✅ success-filters-by-gender-male
- ✅ success-filters-by-gender-female
- ✅ success-filters-by-username
- ✅ success-sorts-by-created-at-desc
- ✅ success-sorts-by-username-asc
- ✅ success-limits-results
- ✅ success-pagination-with-offset
- ✅ success-filters-by-ids
- ✅ success-no-users-found
- ✅ success-combines-multiple-filters

#### TestStore_User (3 cases)
- ✅ success-returns-single-user
- ✅ error-no-user-found
- ✅ error-multiple-users-found

#### TestStore_Update (10 cases)
- ✅ success-updates-username
- ✅ success-updates-multiple-fields
- ✅ success-updates-multiple-users
- ✅ error-no-ids-provided
- ✅ success-nothing-to-update

---

## Expected Output

```bash
$ go test -v ./internal/musicapp/lib/users/store/...

=== RUN   TestStore_Users
=== RUN   TestStore_Users/success-returns-all-users
=== PAUSE TestStore_Users/success-returns-all-users
=== RUN   TestStore_Users/success-filters-by-gender-male
=== PAUSE TestStore_Users/success-filters-by-gender-male
...
=== CONT  TestStore_Users/success-returns-all-users
=== CONT  TestStore_Users/success-filters-by-gender-male
--- PASS: TestStore_Users (0.15s)
    --- PASS: TestStore_Users/success-returns-all-users (0.03s)
    --- PASS: TestStore_Users/success-filters-by-gender-male (0.02s)
    ...

=== RUN   TestStore_User
--- PASS: TestStore_User (0.08s)

=== RUN   TestStore_Update
--- PASS: TestStore_Update (0.12s)

PASS
ok      mlm/internal/musicapp/lib/users/store   0.350s
```

---

## Troubleshooting

### Issue: "database not initialized"

**Solution:** Update testsuite.go with correct test DB credentials:

```go
// internal/testsuite/testsuite.go
dsn := "user:password@tcp(127.0.0.1:3306)/mlm_test?parseTime=true"
```

### Issue: "table doesn't exist"

**Solution:** Run migrations on test database:
```bash
mysql -u user -p mlm_test < migration/01_create_users.up.sql
```

### Issue: Tests are slow

**Solution:** Use t.Parallel() - already implemented!

---

## Creating More Tests

### Test Template

```go
func TestStore_YourMethod(t *testing.T) {
    type testCase struct {
        name            string
        setup           func(th *testsuite.Helper) YourFilter
        extraAssertions func(th *testsuite.Helper, result []*YourModel, err error)
    }

    testCases := []testCase{
        {
            name: "success-your-scenario",
            setup: func(th *testsuite.Helper) YourFilter {
                // Create test data
                factory.YourEntity(th.T, th.BackendAppDb(), nil)

                return YourFilter{
                    // Your filter params
                }
            },
            extraAssertions: func(th *testsuite.Helper, result []*YourModel, err error) {
                require.NoError(th.T, err)
                assert.Len(th.T, result, 1)
            },
        },
    }

    for _, tt := range testCases {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            testSuite := testsuite.New(t)
            t.Cleanup(testSuite.UseBackendDB())

            store := store.New()
            filter := tt.setup(testSuite)

            result, err := store.YourMethod(
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

---

## Test Checklist

Before running tests:
- ✅ Test database created (`mlm_test`)
- ✅ Migrations run on test DB
- ✅ Correct credentials in testsuite.go
- ✅ Dependencies installed (`go get`)

When writing tests:
- ✅ Use table-driven pattern
- ✅ Use factories for test data
- ✅ Add `t.Parallel()` in subtests
- ✅ Test happy paths AND error cases
- ✅ Clean assertions with testify

---

## Next Steps

1. **Run existing tests** to verify they work
2. **Create rooms store tests** following same pattern
3. **Add more test cases** as needed
4. **Keep test coverage high** (aim for >80%)

---

**You have 25 comprehensive tests ready to run!** 🎉
