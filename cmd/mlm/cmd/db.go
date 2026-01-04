package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations",
	Long: `Perform database operations like recreate, reset, seed.

Subcommands:
  recreate  Drop and recreate all tables
  reset     Truncate all tables (keep schema)
  seed      Seed database with test data

Examples:
  mlm db recreate
  mlm db reset
  mlm db seed`,
}

var dbRecreateCmd = &cobra.Command{
	Use:   "recreate",
	Short: "Drop and recreate database schema",
	Long: `Drop all tables and recreate from migrations.

WARNING: This will DELETE ALL DATA!

Use this for development only.`,
	Run: func(cmd *cobra.Command, args []string) {
		recreateDatabase()
	},
}

var dbResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Truncate all tables",
	Long: `Remove all data but keep schema intact.

This is faster than recreate and useful for testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		resetDatabase()
	},
}

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with test data",
	Long: `Insert test data into the database.

Useful for development and testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		seedDatabase()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbRecreateCmd)
	dbCmd.AddCommand(dbResetCmd)
	dbCmd.AddCommand(dbSeedCmd)
}

func recreateDatabase() {
	log.Println("🗑️  Recreating database schema...")
	log.Println("⚠️  WARNING: This will DELETE ALL DATA!")

	db := connectDB()
	defer db.Close()

	// Drop all tables
	log.Println("📦 Dropping all tables...")
	tables := []string{"room_members", "rooms", "users"}

	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatalf("❌ Failed to disable foreign key checks: %v", err)
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Printf("⚠️  Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("✅ Dropped table: %s", table)
		}
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatalf("❌ Failed to re-enable foreign key checks: %v", err)
	}

	log.Println("🎉 Database schema dropped")
	log.Println("💡 Run 'mlm migrate' to recreate tables")
}

func resetDatabase() {
	log.Println("🧹 Resetting database (truncating all tables)...")

	db := connectDB()
	defer db.Close()

	tables := []string{"room_members", "rooms", "users"}

	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatalf("❌ Failed to disable foreign key checks: %v", err)
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			log.Printf("⚠️  Failed to truncate table %s: %v", table, err)
		} else {
			log.Printf("✅ Truncated table: %s", table)
		}
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatalf("❌ Failed to re-enable foreign key checks: %v", err)
	}

	log.Println("🎉 Database reset complete")
}

func seedDatabase() {
	log.Println("🌱 Seeding database with test data...")

	db := connectDB()
	defer db.Close()

	// Seed users
	log.Println("👥 Creating test users...")
	users := []struct {
		username    string
		displayName string
		gender      string
	}{
		{"alice", "Alice Smith", "female"},
		{"bob", "Bob Jones", "male"},
		{"charlie", "Charlie Brown", "male"},
		{"diana", "Diana Prince", "female"},
		{"eve", "Eve Anderson", "female"},
	}

	for _, u := range users {
		_, err := db.Exec(`
			INSERT INTO users (username, display_name, gender)
			VALUES (?, ?, ?)
		`, u.username, u.displayName, u.gender)

		if err != nil {
			log.Printf("⚠️  Failed to create user %s: %v", u.username, err)
		} else {
			log.Printf("✅ Created user: %s", u.username)
		}
	}

	log.Println("🎉 Database seeded successfully")
}

func connectDB() *sql.DB {
	dbHost := getEnv("MUSICAPP_PG_HOST", "127.0.0.1")
	dbPort := getEnv("MUSICAPP_PG_PORT", "3306")
	dbUser := getEnv("MUSICAPP_PG_USER", "user")
	dbPass := getEnv("MUSICAPP_PG_PASS", "password")
	dbName := getEnv("MUSICAPP_PG_DATABASE", "mlm")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	return db
}
