package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// sqlboilerCmd generates SQLBoiler models from database schema
var sqlboilerCmd = &cobra.Command{
	Use:   "sqlboiler",
	Short: "Generate SQLBoiler models from database",
	Long: `Generate Go models from your database schema using SQLBoiler.

Reads configuration from sqlboiler.toml

Examples:
  mlm sqlboiler`,
	Run: func(cmd *cobra.Command, args []string) {
		runSQLBoiler()
	},
}

func init() {
	rootCmd.AddCommand(sqlboilerCmd)
}

func runSQLBoiler() {
	log.Println("🔧 Generating SQLBoiler models...")

	// Check if sqlboiler.toml exists
	if _, err := os.Stat("sqlboiler.toml"); os.IsNotExist(err) {
		log.Println("⚠️  sqlboiler.toml not found")
		log.Println("💡 Create sqlboiler.toml with your database configuration")
		return
	}

	// Run sqlboiler command
	cmd := exec.Command("sqlboiler", "mysql", "--wipe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ SQLBoiler failed: %v", err)
	}

	log.Println("✅ SQLBoiler models generated successfully")
	log.Println("💡 Models are in the models/ directory")
}
