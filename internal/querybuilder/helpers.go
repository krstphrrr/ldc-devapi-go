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

// Map of valid operators and their SQL equivalents
var validOperators = map[string]string{
	"$lt":  "<",
	"$gt":  ">",
	"$lte": "<=",
	"$gte": ">=",
	"$ne":  "<>",
}

// Handle standard and operator-based queries
func HandleQueryParam(key string, value interface{}, columnType string, valueIndex int) (string, []interface{}) {
	var sqlFragment strings.Builder
	var values []interface{}

	switch v := value.(type) {
	case string:
		// Handle date parameters like "start" or "end"
		if columnType == "date" && (key == "start" || key == "end") {
			sqlFragment.WriteString(handleDateParam(key, v, valueIndex))
			values = append(values, v)
		} else {
			sqlFragment.WriteString(fmt.Sprintf(` AND "%s" = $%d`, key, valueIndex))
			values = append(values, v)
		}
	case map[string]interface{}:
		// Handle operator-based queries like {"$gt": 10, "$lt": 20}
		fragment, paramValues := handleOperatorBasedQuery(key, v, valueIndex)
		sqlFragment.WriteString(fragment)
		values = append(values, paramValues...)
	case []interface{}:
		// Handle arrays for queries like {"Name": ["Alice", "Bob"]}
		fragment, paramValues := handleArrayParam(key, v, valueIndex)
		sqlFragment.WriteString(fragment)
		values = append(values, paramValues...)
	default:
		// Default handling for single values
		sqlFragment.WriteString(fmt.Sprintf(` AND "%s" = $%d`, key, valueIndex))
		values = append(values, value)
	}

	return sqlFragment.String(), values
}

// Handle date parameters (e.g., "start", "end")
func handleDateParam(key string, value string, valueIndex int) string {
	operator := ">="
	if key == "end" {
		operator = "<"
	}
	return fmt.Sprintf(` AND "DateVisited" %s $%d`, operator, valueIndex)
}

// Handle operator-based queries (e.g., {"$gt": 10})
func handleOperatorBasedQuery(key string, operators map[string]interface{}, valueIndex int) (string, []interface{}) {
	var sqlFragment strings.Builder
	var values []interface{}

	for op, val := range operators {
		if operator, ok := validOperators[op]; ok {
			sqlFragment.WriteString(fmt.Sprintf(` AND "%s" %s $%d`, key, operator, valueIndex))
			values = append(values, val)
			valueIndex++
		} else {
			panic(fmt.Sprintf("Invalid operator: %s", op))
		}
	}

	return sqlFragment.String(), values
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