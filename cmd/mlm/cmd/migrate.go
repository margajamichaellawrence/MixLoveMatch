package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	migrateDown bool
	migrateStep int
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Run database migrations to update the schema.

Migrations are SQL files in the migration/ directory:
- *.up.sql   = Apply migration
- *.down.sql = Rollback migration

Examples:
  mlm migrate              # Run all pending migrations
  mlm migrate --down       # Rollback last migration
  mlm migrate --step 2     # Run next 2 migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVar(&migrateDown, "down", false, "Rollback migrations")
	migrateCmd.Flags().IntVar(&migrateStep, "step", 0, "Number of migrations to run (0 = all)")
}

func runMigrations() {
	log.Println("🔄 Running database migrations...")

	// Get database configuration
	dbHost := getEnv("MUSICAPP_PG_HOST", "127.0.0.1")
	dbPort := getEnv("MUSICAPP_PG_PORT", "3306")
	dbUser := getEnv("MUSICAPP_PG_USER", "user")
	dbPass := getEnv("MUSICAPP_PG_PASS", "password")
	dbName := getEnv("MUSICAPP_PG_DATABASE", "mlm")

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Get migration files
	migrationDir := "migration"
	if migrateDown {
		runDownMigrations(db, migrationDir)
	} else {
		runUpMigrations(db, migrationDir)
	}
}

func runUpMigrations(db *sql.DB, dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		log.Fatalf("❌ Failed to read migration files: %v", err)
	}

	sort.Strings(files)

	if len(files) == 0 {
		log.Println("ℹ️  No migration files found")
		return
	}

	count := 0
	for _, file := range files {
		if migrateStep > 0 && count >= migrateStep {
			break
		}

		log.Printf("📄 Running: %s", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("❌ Failed to read %s: %v", file, err)
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("❌ Migration failed for %s: %v", file, err)
		}

		log.Printf("✅ Applied: %s", filepath.Base(file))
		count++
	}

	log.Printf("🎉 Successfully applied %d migration(s)", count)
}

func runDownMigrations(db *sql.DB, dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.down.sql"))
	if err != nil {
		log.Fatalf("❌ Failed to read migration files: %v", err)
	}

	// Sort in reverse order for rollback
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	if len(files) == 0 {
		log.Println("ℹ️  No rollback files found")
		return
	}

	count := 0
	limit := migrateStep
	if limit == 0 {
		limit = 1 // Default: rollback only last migration
	}

	for _, file := range files {
		if count >= limit {
			break
		}

		log.Printf("📄 Rolling back: %s", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("❌ Failed to read %s: %v", file, err)
		}

		// Execute rollback
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("❌ Rollback failed for %s: %v", file, err)
		}

		log.Printf("✅ Rolled back: %s", filepath.Base(file))
		count++
	}

	log.Printf("🎉 Successfully rolled back %d migration(s)", count)
}

// Migration tracking table (optional - for future enhancement)
func createMigrationTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationName(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".up.sql")
}
