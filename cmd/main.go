package main

import (
	"log"
	"net/http"

	"go-api-app/config"
	"go-api-app/internal/database"
	"go-api-app/internal/routes"
)

func main() {
	// Load the configuration
	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database connection for the default tenant
	defaultTenant := "public" 
	db, err := database.GetTenantDB(defaultTenant)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer db.Close()

	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Register all routes under "/api/v1"
	apiMux := http.NewServeMux()
	routes.RegisterRoutes(apiMux, db)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
