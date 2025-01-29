package main

import (
	"log"
	"net/http"

	"go-api-app/config"
	"go-api-app/internal/database"
	"go-api-app/internal/routes"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"https://api.landscapedatacommons.org",
			"https://devapi.landscapedatacommons.org",
		}

		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	mux.Handle("/api/v1/", corsMiddleware(http.StripPrefix("/api/v1", apiMux)))

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
