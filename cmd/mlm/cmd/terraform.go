package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// terraformCmd combines db recreate + migrate + seed
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Drop, recreate, and seed database (all-in-one)",
	Long: `Terraform combines multiple operations:
1. Drop all tables (db recreate)
2. Run migrations
3. Seed test data

WARNING: This will DELETE ALL DATA!

Use this for quick development reset.

Examples:
  mlm terraform`,
	Run: func(cmd *cobra.Command, args []string) {
		runTerraform()
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}

func runTerraform() {
	log.Println("🏗️  Terraforming database (recreate + migrate + seed)...")
	log.Println("⚠️  WARNING: This will DELETE ALL DATA!")

	// Step 1: Recreate (drop tables)
	log.Println("\n📦 Step 1/3: Dropping tables...")
	recreateDatabase()

	// Step 2: Run migrations
	log.Println("\n🔄 Step 2/3: Running migrations...")
	runUpMigrations(connectDB(), "migration")

	// Step 3: Seed data
	log.Println("\n🌱 Step 3/3: Seeding test data...")
	seedDatabase()

	log.Println("\n🎉 Terraform complete! Database is ready for development.")
}
