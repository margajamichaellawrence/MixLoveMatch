package testsuite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/aarondl/sqlboiler/v4/boil"
)

// Helper provides test infrastructure
type Helper struct {
	T   *testing.T
	Ctx context.Context
	db  *sql.DB
}

// New creates a new test helper
func New(t *testing.T) *Helper {
	return &Helper{
		T:   t,
		Ctx: context.Background(),
	}
}

// UseBackendDB connects to test database and returns cleanup function
func (h *Helper) UseBackendDB() func() {
	// Connect to test database
	// Using environment variables or hardcoded test DB
	dsn := "user:userpass@tcp(127.0.0.1:3306)/musicapp_test?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		h.T.Fatalf("failed to connect to test database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		h.T.Fatalf("failed to ping test database: %v", err)
	}

	h.db = db

	// Clean up test data before each test
	h.cleanDatabase()

	// Return cleanup function
	return func() {
		h.cleanDatabase()
		if err := db.Close(); err != nil {
			h.T.Errorf("failed to close database: %v", err)
		}
	}
}

// BackendAppDb returns the test database connection
func (h *Helper) BackendAppDb() boil.ContextExecutor {
	if h.db == nil {
		h.T.Fatal("database not initialized - call UseBackendDB() first")
	}
	return h.db
}

// cleanDatabase truncates all tables for clean test state
func (h *Helper) cleanDatabase() {
	tables := []string{
		"room_members",
		"rooms",
		"users",
	}

	// Disable foreign key checks
	_, err := h.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		h.T.Logf("warning: failed to disable foreign key checks: %v", err)
	}

	// Truncate each table
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)
		_, err := h.db.Exec(query)
		if err != nil {
			h.T.Logf("warning: failed to truncate table %s: %v", table, err)
		}
	}

	// Re-enable foreign key checks
	_, err = h.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		h.T.Logf("warning: failed to re-enable foreign key checks: %v", err)
	}
}

// BeginTx starts a transaction for testing
func (h *Helper) BeginTx() (*sql.Tx, error) {
	return h.db.BeginTx(h.Ctx, nil)
}
