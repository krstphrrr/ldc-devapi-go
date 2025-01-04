package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"go-api-app/config"
	"go-api-app/internal/repositories"
	"go-api-app/internal/querybuilder"
)

func GenericDynamicDataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Handling request:", r.URL.Path)

	// Extract endpoint and map to table name
	endpoint := r.URL.Path
	tableName, ok := config.EndpointToTableMap[endpoint]
	if !ok {
		log.Printf("Invalid endpoint: %s\n", endpoint)
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}
	log.Printf("Mapped endpoint to table: %s\n", tableName)

	// Fetch column names dynamically
	columns, err := querybuilder.FetchColumns(db, tableName)
	if err != nil {
		log.Printf("Failed to fetch columns for table: %s, error: %v\n", tableName, err)
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched columns: %v\n", columns)

	// Parse query parameters
	queryParams := map[string]interface{}{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}
	log.Printf("Parsed query parameters: %v\n", queryParams)

	// Fetch data from the repository
	data, err := repositories.FetchData(db, tableName, columns, queryParams)
	if err != nil {
		log.Printf("Failed to fetch data for table: %s, error: %v\n", tableName, err)
		http.Error(w, "Internal Server Error: Failed to fetch data", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched data successfully for table: %s\n", tableName)

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
