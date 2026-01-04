package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	port string
	host string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the REST API server",
	Long: `Start the Music Dating App REST API server.

The server will:
- Connect to MySQL database
- Initialize stores and handlers
- Register HTTP routes
- Listen on specified host:port

Environment Variables:
  MUSICAPP_PG_HOST      Database host (default: 127.0.0.1)
  MUSICAPP_PG_PORT      Database port (default: 3306)
  MUSICAPP_PG_USER      Database user (default: user)
  MUSICAPP_PG_PASS      Database password (default: password)
  MUSICAPP_PG_DATABASE  Database name (default: mlm)

Examples:
  mlm serve
  mlm serve --port 8080
  mlm serve --host 0.0.0.0 --port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	serveCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "Host to bind to")
}

func runServer() {
	log.Println("🚀 Starting Music Dating App API Server...")

	// Get database configuration from environment
	dbHost := getEnv("MUSICAPP_PG_HOST", "127.0.0.1")
	dbPort := getEnv("MUSICAPP_PG_PORT", "3306")
	dbUser := getEnv("MUSICAPP_PG_USER", "user")
	dbPass := getEnv("MUSICAPP_PG_PASS", "userpass")
	dbName := getEnv("MUSICAPP_PG_DATABASE", "musicapp")

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("📊 Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// Connect to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Database connected successfully")

	// TODO: Initialize stores and handlers following IMAPP pattern
	// userStore := userstore.New()
	// userLogic := users.NewLogic(userStore)
	// userHandler := userhandler.New(userLogic)

	// Setup routes
	log.Println("🛣️  Setting up server...")
	mux := http.NewServeMux()

	// TODO: Register routes when handlers are implemented
	// mux.HandleFunc("GET /users", userHandler.ListUsers)
	// mux.HandleFunc("POST /users", userHandler.CreateUser)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("✅ Server configured")

	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🎵 Server listening on http://%s", addr)
	log.Printf("📝 Available Endpoints:")
	log.Printf("   - GET /health")
	log.Printf("")
	log.Printf("Press Ctrl+C to stop")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
