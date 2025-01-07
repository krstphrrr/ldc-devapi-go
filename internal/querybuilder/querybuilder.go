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
func GenerateQuery(table string, columnTypes map[string]string, rawParams url.Values, rawQuery string) (string, []interface{}, error) {
	if table == "" {
		return "", nil, fmt.Errorf("table name cannot be empty")
	}

	var sqlQuery strings.Builder
	var values []interface{}
	// valueIndex := 1

	columns := make([]string, 0, len(columnTypes))
	for col := range columnTypes {
		columns = append(columns, fmt.Sprintf(`"%s"`, col))
	}
	sqlQuery.WriteString(fmt.Sprintf("SELECT %s FROM public_test.%s WHERE 1 = 1", strings.Join(columns, ", "), table))

	// Parse all query parameters except limit and offset
    queryFragment, _, err := ParseEncodedQueryFromRaw(rawQuery, 1, &values)
    if err != nil {
        return "", nil, fmt.Errorf("failed to parse query parameters: %v", err)
    }
    sqlQuery.WriteString(queryFragment)
    // valueIndex = newIndex
	
	

	// Add ORDER BY clause
	sqlQuery.WriteString(" ORDER BY rid ASC")

	// Handle LIMIT and OFFSET separately
    if limit := rawParams.Get("limit"); limit != "" {
        sqlQuery.WriteString(fmt.Sprintf(" LIMIT %s", limit))
    }
    if offset := rawParams.Get("offset"); offset != "" {
        sqlQuery.WriteString(fmt.Sprintf(" OFFSET %s", offset))
    }

	log.Printf("Final generated query: %s", sqlQuery.String())
	return sqlQuery.String(), values, nil
}



func GenerateQueryFromBody(table string, columnTypes map[string]string, body map[string]interface{}) (string, []interface{}, error) {
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
	sqlQuery.WriteString(fmt.Sprintf("SELECT %s FROM public_test.%s WHERE 1 = 1", strings.Join(columns, ", "), table))

	for key, condition := range body {
		if _, exists := columnTypes[key]; !exists {
			log.Printf("Column %s not found in column metadata, skipping\n", key)
			continue
		}

		switch v := condition.(type) {
		case map[string]interface{}:
			// Handle operator-based conditions
			for operator, value := range v {
				// Validate operator value is not an array
				if _, isArray := value.([]interface{}); isArray {
					log.Printf("Invalid value for operator %s on column %s: Arrays are not allowed\n", operator, key)
					return "", nil, fmt.Errorf("invalid value for operator %s on column %s: arrays are not allowed", operator, key)
				}

				sqlOperator, err := parseBodyOperator(operator)
				if err != nil {
					log.Printf("Unsupported operator %s for column %s\n", operator, key)
					continue
				}
				sqlQuery.WriteString(fmt.Sprintf(` AND "%s" %s $%d`, key, sqlOperator, valueIndex))
				values = append(values, value)
				valueIndex++
			}
		case []interface{}:
			// Handle array-based conditions
			arrayFragment, arrayValues := handleArrayParam(key, v, valueIndex)
			sqlQuery.WriteString(arrayFragment)
			values = append(values, arrayValues...)
			valueIndex += len(arrayValues)
		default:
			// Default equality condition
			sqlQuery.WriteString(fmt.Sprintf(` AND "%s" = $%d`, key, valueIndex))
			values = append(values, v)
			valueIndex++
		}
	}

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