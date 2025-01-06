package querybuilder

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func parseOperator(operator string) (string, error) {
	switch operator {
	case "%3C", "<": // <
		return "<", nil
	case "%3E", ">": // >
		return ">", nil
	case "%3C%3D", "<=": // <=
		return "<=", nil
	case "%3E%3D", ">=": // >=
		return ">=", nil
	case "%21%3D", "!=": // !=
		return "<>", nil
	case "=": // =
		return "=", nil
	default:
		return "", fmt.Errorf("unsupported operator: %s", operator)
	}
}

func parseBodyOperator(operator string) (string, error) {
	switch operator {
	case "$lt": // Less than
		return "<", nil
	case "$gt": // Greater than
		return ">", nil
	case "$lte": // Less than or equal
		return "<=", nil
	case "$gte": // Greater than or equal
		return ">=", nil
	case "$ne": // Not equal
		return "<>", nil
	default:
		return "", fmt.Errorf("unsupported operator: %s", operator)
	}
}
func ParseEncodedQuery(queryParams map[string][]string, valueIndex int, values *[]interface{}) (string, int, error) {
	var sqlQuery strings.Builder

	// Regular expression to capture column name and operator
	re := regexp.MustCompile(`^([a-zA-Z0-9_]+)(%[0-9A-Za-z]{2}|[<>!=]+)?$`)

	for key, vals := range queryParams {
		if len(vals) == 0 {
			continue
		}

		// Match the key to extract the column and operator
		matches := re.FindStringSubmatch(key)
		if len(matches) < 2 {
			log.Printf("Invalid parameter format: %s", key)
			continue
		}

		column := matches[1]
		rawOperator := matches[2]

		var operator string
		var err error

		if rawOperator == "" {
			// Default to equality if no operator is specified
			operator = "="
		} else {
			operator, err = parseOperator(rawOperator)
			if err != nil {
				log.Printf("Failed to parse operator: %v", err)
				continue
			}
		}

		// Add SQL fragment
		sqlQuery.WriteString(fmt.Sprintf(` AND "%s" %s $%d`, column, operator, valueIndex))
		*values = append(*values, vals[0]) // Assuming single value per key
		valueIndex++
	}

	return sqlQuery.String(), valueIndex, nil
}


// Handle array parameters (e.g., {"Name": ["Alice", "Bob"]})
func handleArrayParam(key string, valuesArray []interface{}, valueIndex int) (string, []interface{}) {
	var sqlFragment strings.Builder
	placeholders := make([]string, len(valuesArray))

	for i := range valuesArray {
		placeholders[i] = fmt.Sprintf("$%d", valueIndex)
		valueIndex++
	}

	sqlFragment.WriteString(fmt.Sprintf(` AND "%s" IN (%s)`, key, strings.Join(placeholders, ", ")))
	return sqlFragment.String(), valuesArray
}


func AddLimitOffsetToBuilder(sqlQuery *strings.Builder, parsedParams map[string]map[string]interface{}) {
    // Handle LIMIT
    if limit, ok := parsedParams["limit"]; ok {
        log.Printf("Adding LIMIT: %v", limit)
        sqlQuery.WriteString(fmt.Sprintf(" LIMIT %v", limit))
    }

    // Handle OFFSET
    if offset, ok := parsedParams["offset"]; ok {
        log.Printf("Adding OFFSET: %v", offset)
        sqlQuery.WriteString(fmt.Sprintf(" OFFSET %v", offset))
    }
}