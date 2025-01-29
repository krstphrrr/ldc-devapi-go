package repositories

import (
	"database/sql"
	"fmt"
	// "go-api-app/internal/querybuilder"
	"log"
	// "strings"
	"go-api-app/internal/database"
)

func FetchDataForTenant(tenant string, db *sql.DB, query string, params []interface{}) ([]map[string]interface{}, error) {
	log.Printf("Fetching data for tenant: %s\n", tenant)

	// Get tenant-specific database connection
	tenantDB, err := database.GetTenantDB(tenant)
	if err != nil {
		log.Printf("Failed to get database connection for tenant: %s, error: %v\n", tenant, err)
		return nil, fmt.Errorf("failed to get database for tenant: %v", err)
	}
	defer tenantDB.Close()

	log.Printf("Executing query: %s with params: %v\n", query, params)

	// Execute the query
	rows, err := tenantDB.Query(query, params...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Parse results into a map
	var results []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Failed to fetch column names: %v\n", err)
		return nil, fmt.Errorf("failed to fetch column names: %v", err)
	}

	for rows.Next() {
		row := make(map[string]interface{})
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			var value interface{}
			columnPointers[i] = &value
		}

		if err := rows.Scan(columnPointers...); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, fmt.Errorf("row scan failed: %v", err)
		}

		// for i, col := range columns {
		// 	row[col] = *(columnPointers[i].(*interface{}))
		// }
		for i, col := range columns {
			rawValue := *(columnPointers[i].(*interface{}))
		
			if col == "delay_range" {
				switch v := rawValue.(type) {
				case []uint8:  // postgreSQL ENUMs often return as []uint8 (byte slices)
					row[col] = string(v)  // convert byte slice to string
				case string:
					row[col] = v  // I]if it's already a string, keep it
				default:
					row[col] = "UNKNOWN"  // unexpected cases
				}
			} else {
				row[col] = rawValue
			}
		}
		results = append(results, row)
	}

	log.Printf("Fetched %d rows for tenant: %s", len(results), tenant)
	return results, nil
}


