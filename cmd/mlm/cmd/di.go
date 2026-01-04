package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// diCmd handles dependency injection operations
var diCmd = &cobra.Command{
	Use:   "di",
	Short: "Dependency injection operations",
	Long: `Manage dependency injection container.

Subcommands:
  generate  Generate DI container code

Examples:
  mlm di generate`,
}

var diGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate dependency injection container",
	Long: `Generate dependency injection container code.

This will scan your code and generate DI wiring.

Examples:
  mlm di generate`,
	Run: func(cmd *cobra.Command, args []string) {
		generateDI()
	},
}

func init() {
	rootCmd.AddCommand(diCmd)
	diCmd.AddCommand(diGenerateCmd)
}

func generateDI() {
	log.Println("🔧 Generating dependency injection container...")
	log.Println("ℹ️  DI generation not yet implemented")
	log.Println("💡 For now, manually wire dependencies in main.go or cmd/serve.go")
	log.Println()
	log.Println("Planned features:")
	log.Println("  - Auto-generate DI container from type signatures")
	log.Println("  - Support for dingo or wire")
	log.Println("  - Lifecycle management")
}
