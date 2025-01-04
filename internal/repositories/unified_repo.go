package repositories

import (
	"database/sql"
	"fmt"
	"go-api-app/internal/querybuilder"
	"log"
	"strings"
)

func FetchData(db *sql.DB, table string, columns []string, queryParams map[string]interface{}) ([]map[string]interface{}, error) {
	log.Printf("Generating query for table: %s, columns: %v, params: %v\n", table, columns, queryParams)

	// Generate the dynamic SQL query
	sqlQuery, values, err := querybuilder.GenerateQuery(table, columns, queryParams)
	if err != nil {
		log.Printf("Query generation failed for table: %s, error: %v\n", table, err)
		return nil, fmt.Errorf("query generation failed: %v", err)
	}
	log.Printf("Generated query: %s, values: %v\n", sqlQuery, values)

	// Execute the query
	rows, err := db.Query(sqlQuery, values...)
	if err != nil {
		log.Printf("Query execution failed for table: %s, error: %v\n", table, err)
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// Parse results into a map
	var results []map[string]interface{}
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
		for i, col := range columns {
			// Strip quotes from column names
			cleanCol := strings.Trim(col, `"`)
			row[cleanCol] = *(columnPointers[i].(*interface{}))
		}
		results = append(results, row)
	}
	log.Printf("Fetched %d rows from table %s", len(results), table)

	return results, nil
}
