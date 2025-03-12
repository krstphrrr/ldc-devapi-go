package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-api-app/config"
	"go-api-app/internal/database"
	"go-api-app/internal/middleware"
	"go-api-app/internal/routes"
	"go-api-app/internal/version"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "user_id"}, // Include user_id
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response times for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "user_id"}, // Include user_id
	)
)

// Middleware to track requests
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user_id from request context
		userID := "anonymous"
		if ctxUserID, ok := r.Context().Value(middleware.UserIDKey).(string); ok {
			userID = ctxUserID
		}

		// Track request time
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, userID))
		defer timer.ObserveDuration()

		// Increment request counter
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, userID).Inc()

		next.ServeHTTP(w, r)
	})
}

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

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	 // Register the root path endpoint
	 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
				"appVersion": version.Version,
				"lastUpdated": version.BuildDate,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
})

	// Register all routes under "/api/v1"
	apiMux := http.NewServeMux()
	routes.RegisterRoutes(apiMux, db)
	mux.Handle("/api/v1/", corsMiddleware(
		metricsMiddleware(middleware.AuthMiddleware(http.StripPrefix("/api/v1", apiMux))),
	))
	mux.Handle("/metrics", promhttp.Handler())
	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
