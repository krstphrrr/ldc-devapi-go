package querybuilder

import (
	// "database/sql"
	"fmt"
	"net/url"

	// "regexp"

	"go-api-app/internal/schemas"
	"log"
	"strings"
)

// GenerateQuery dynamically generates an SQL query based on input parameters.
func GenerateQuery(table string, rawParams url.Values, columnTypes map[string]string, rawQuery string) (string, []interface{}, error) {
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
 	sqlQuery.WriteString(fmt.Sprintf("SELECT %s FROM %s WHERE 1 = 1", strings.Join(columns, ", "), table))
	//  sqlQuery.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE 1 = 1", table))

	// Parse all query parameters except limit and offset
	queryFragment, _, err := ParseEncodedQueryFromRaw(rawQuery, 1, &values)
	if err != nil {
			return "", nil, fmt.Errorf("failed to parse query parameters: %v", err)
	}
	sqlQuery.WriteString(queryFragment)
	// valueIndex = newIndex
	
	

	// Add ORDER BY clause if table name does not contain "tblProject"
	if !strings.Contains(table, "tblproject") {
		sqlQuery.WriteString(" ORDER BY rid ASC")
	}

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
	// sqlQuery.WriteString(fmt.Sprintf("SELECT * FROM public_test.%s WHERE 1 = 1", table))

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

	// Add ORDER BY clause if table name does not contain "tblProject"
	if !strings.Contains(table, "tblproject") {
		sqlQuery.WriteString(" ORDER BY rid ASC")
	}

	log.Printf("Final generated query: %s", sqlQuery.String())
	return sqlQuery.String(), values, nil
}


// FetchColumns dynamically fetches column names and types for the specified table or view.
func FetchColumns(fullTableName string) (map[string]string, error) {
	log.Printf("Fetching columns and types for table or view %s", fullTableName)

	// Check if the table exists in the `TableSchemas` map
	columnTypes, exists := schemas.TableSchemas[fullTableName]
	if !exists {
		log.Printf("Table schema not found for: %s", fullTableName)
		return nil, fmt.Errorf("table schema not found for: %s", fullTableName)
	}

	if len(columnTypes) == 0 {
		log.Printf("No columns found for table or view: %s", fullTableName)
		return nil, fmt.Errorf("no columns found for table or view: %s", fullTableName)
	}

	log.Printf("Fetched columns and types for table %s: %v", fullTableName, columnTypes)
	return columnTypes, nil
}