package querybuilder

import (
	"database/sql"
	"fmt"
	// "go-api-app/config"
	"log"
	"strings"
)

// GenerateQuery dynamically generates an SQL query based on input parameters.
func GenerateQuery(table string, columns []string, queryParams map[string]interface{}) (string, []interface{}, error) {
	
	if table == "" {
		return "", nil, fmt.Errorf("table name cannot be empty")
	}

	var sqlQuery strings.Builder
	var values []interface{}
	valueIndex := 1

	// Wrap each column in double quotes to preserve case sensitivity
	for i, col := range columns {
		columns[i] = fmt.Sprintf(`"%s"`, col)
	}

	// Start building the SQL query
	sqlQuery.WriteString(fmt.Sprintf("SELECT %s FROM public_test.%s WHERE 1 = 1", strings.Join(columns, ", "), table))

	// Process query parameters for filters
	for key, value := range queryParams {
		switch key {
		case "limit", "offset":
			// Skip here; handle separately in AddLimitOffset
		default:
			sqlQuery.WriteString(fmt.Sprintf(" AND %s = $%d", key, valueIndex))
			values = append(values, value)
			valueIndex++
		}
	}

	// Add ORDER BY
	sqlQuery.WriteString(" ORDER BY rid ASC")

	// Append LIMIT and OFFSET
	AddLimitOffsetToBuilder(&sqlQuery, queryParams)
	log.Printf("Final generated query: %s", sqlQuery.String())

	return sqlQuery.String(), values, nil
}

// FetchColumns dynamically fetches column names for the specified table or view.
func FetchColumns(db *sql.DB, tableName string) ([]string, error) {
	log.Printf("Fetching columns for table or view %s", tableName)

	query := fmt.Sprintf(`
		SELECT column_name
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

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			log.Printf("Error scanning column: %v", err)
			return nil, fmt.Errorf("failed to scan column name: %v", err)
		}
		columns = append(columns, column)
	}

	if len(columns) == 0 {
		log.Printf("No columns found for table or view: %s", tableName)
		return nil, fmt.Errorf("no columns found for table or view: %s", tableName)
	}

	log.Printf("Fetched columns for table %s: %v", tableName, columns)
	return columns, nil
}