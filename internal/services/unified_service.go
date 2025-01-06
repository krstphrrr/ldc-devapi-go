package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"go-api-app/config"
	"go-api-app/internal/middleware"
	"go-api-app/internal/querybuilder"
	"go-api-app/internal/repositories"
)

func GenericDynamicDataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Handling request:", r.URL.Path)

	// Retrieve tenant from context
	tenant, ok := r.Context().Value(middleware.TenantKey).(string)
	if !ok {
		log.Println("Tenant not found in context. Using default 'public'.")
		tenant = "public"
	}
	log.Printf("Processing request for tenant: %s\n", tenant)

	// Extract endpoint and map to table name
	endpoint := r.URL.Path
	tableName, ok := config.EndpointToTableMap[endpoint]
	if !ok {
		log.Printf("Invalid endpoint: %s\n", endpoint)
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	// Fetch column names and types dynamically
	columnTypes, err := querybuilder.FetchColumns(db, tableName)
	if err != nil {
		log.Printf("Failed to fetch columns for table: %s, error: %v\n", tableName, err)
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
		return
	}


	// Generate the query
	query, params, err := querybuilder.GenerateQuery(tableName, columnTypes, r.URL.Query())
	if err != nil {
		log.Printf("Failed to generate query: %v\n", err)
		http.Error(w, "Internal Server Error: Failed to generate query", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated Query: %s\n", query)
	log.Printf("Query Params: %v\n", params)
	
	// Execute query
	data, err := repositories.FetchDataForTenant(tenant, db, query, params)
if err != nil {
	log.Printf("Failed to fetch data for table: %s, error: %v\n", tableName, err)
	http.Error(w, "Internal Server Error: Failed to fetch data", http.StatusInternalServerError)
	return
}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ONLY FOR POST REQUESTS
func GenericDynamicDataPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Handling POST request:", r.URL.Path)

	// Retrieve tenant from context
	tenant, ok := r.Context().Value(middleware.TenantKey).(string)
	if !ok {
		log.Println("Tenant not found in context. Using default 'public'.")
		tenant = "public"
	}
	log.Printf("Processing request for tenant: %s\n", tenant)

	// Extract endpoint and map to table name
	endpoint := r.URL.Path
	tableName, ok := config.EndpointToTableMap[endpoint]
	if !ok {
		log.Printf("Invalid endpoint: %s\n", endpoint)
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	// Fetch column names and types dynamically
	columnTypes, err := querybuilder.FetchColumns(db, tableName)
	if err != nil {
		log.Printf("Failed to fetch columns for table: %s, error: %v\n", tableName, err)
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
		return
	}

	// Parse request body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed request body: %v\n", body)

	// Generate query from request body
	query, params, err := querybuilder.GenerateQueryFromBody(tableName, columnTypes, body)
	if err != nil {
		log.Printf("Failed to generate query: %v\n", err)
		http.Error(w, "Internal Server Error: Failed to generate query", http.StatusInternalServerError)
		return
	}

	log.Printf("Generated Query: %s\n", query)
	log.Printf("Query Params: %v\n", params)

	// Execute query
	data, err := repositories.FetchDataForTenant(tenant, db, query, params)
	if err != nil {
		log.Printf("Failed to fetch data for table: %s, error: %v\n", tableName, err)
		http.Error(w, "Internal Server Error: Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
