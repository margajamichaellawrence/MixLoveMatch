package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd shows current configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long: `Display current configuration from environment variables.

Shows database connection settings and other environment vars.

Examples:
  mlm config`,
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func showConfig() {
	fmt.Println("📝 Current Configuration:")
	fmt.Println()

	// Database config
	fmt.Println("🗄️  Database:")
	fmt.Printf("  Host:     %s\n", getEnv("MUSICAPP_PG_HOST", "127.0.0.1"))
	fmt.Printf("  Port:     %s\n", getEnv("MUSICAPP_PG_PORT", "3306"))
	fmt.Printf("  User:     %s\n", getEnv("MUSICAPP_PG_USER", "user"))
	fmt.Printf("  Password: %s\n", maskPassword(getEnv("MUSICAPP_PG_PASS", "password")))
	fmt.Printf("  Database: %s\n", getEnv("MUSICAPP_PG_DATABASE", "mlm"))
	fmt.Println()

	// Config file
	if viper.ConfigFileUsed() != "" {
		fmt.Printf("📄 Config file: %s\n", viper.ConfigFileUsed())
	} else {
		fmt.Println("📄 Config file: Not found (using defaults)")
	}
}

func maskPassword(pass string) string {
	if len(pass) <= 2 {
		return "***"
	}
	return pass[:2] + strings.Repeat("*", len(pass)-2)
}

// ViperPG returns viper with PG env vars loaded (for compatibility)
func ViperPG() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("MUSICAPP_PG")
	return v
}
