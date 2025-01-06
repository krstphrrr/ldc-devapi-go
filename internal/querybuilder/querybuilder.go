package querybuilder

import (
	"database/sql"
	"fmt"
	"net/url"

	// "regexp"

	// "go-api-app/config"
	"log"
	"strings"
)

// GenerateQuery dynamically generates an SQL query based on input parameters.
func GenerateQuery(table string, columnTypes map[string]string, rawParams url.Values) (string, []interface{}, error) {
	if table == "" {
		return "", nil, fmt.Errorf("table name cannot be empty")
	}

	var sqlQuery strings.Builder
	var values []interface{}
	valueIndex := 1

	columns := make([]string, 0, len(columnTypes))
	for col := range columnTypes {
		columns = append(columns, fmt.Sprintf(`"%s"`, col))
	}
	// Start building the SQL query
	sqlQuery.WriteString(fmt.Sprintf("SELECT %s FROM public_test.%s WHERE 1 = 1", strings.Join(columns, ", "), table))

	// Parse and append encoded query parameters
	queryFragment, newIndex, err := ParseEncodedQuery(rawParams, valueIndex, &values)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse query parameters: %v", err)
	}
	sqlQuery.WriteString(queryFragment)
	valueIndex = newIndex

	// Add ORDER BY clause
	sqlQuery.WriteString(" ORDER BY rid ASC")

	log.Printf("Final generated query: %s", sqlQuery.String())
	return sqlQuery.String(), values, nil
}

// FetchColumns dynamically fetches column names and types for the specified table or view.
func FetchColumns(db *sql.DB, tableName string) (map[string]string, error) {
	log.Printf("Fetching columns and types for table or view %s", tableName)

	query := fmt.Sprintf(`
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = '%s'
		  AND table_schema = 'public_test'
		  ORDER BY ordinal_position
	`, tableName)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching columns: %v", err)
		return nil, fmt.Errorf("failed to fetch columns for table or view %s: %v", tableName, err)
	}
	defer rows.Close()

	columnTypes := make(map[string]string)
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			log.Printf("Error scanning column: %v", err)
			return nil, fmt.Errorf("failed to scan column name and type: %v", err)
		}
		columnTypes[columnName] = dataType
	}

	if len(columnTypes) == 0 {
		log.Printf("No columns found for table or view: %s", tableName)
		return nil, fmt.Errorf("no columns found for table or view: %s", tableName)
	}

	log.Printf("Fetched columns and types for table %s: %v", tableName, columnTypes)
	return columnTypes, nil
}