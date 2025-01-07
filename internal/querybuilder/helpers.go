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
func ParseEncodedQueryFromRaw(rawQuery string, valueIndex int, values *[]interface{}) (string, int, error) {
	var sqlQuery strings.Builder

	// Regular expression to capture column name, operator, and value
	re := regexp.MustCompile(`(?P<key>[a-zA-Z0-9_]+)(?P<operator>%3C%3D|%3E%3D|%21%3D|%3C|%3E|<=|>=|<|>|=)?(?P<value>.+)?`)

	// Manually split raw query string
	params := strings.Split(rawQuery, "&")

	for _, param := range params {
		// Match each parameter with regex
		if strings.Contains(param, "limit") || strings.Contains(param, "offset") {
            // Skip limit and offset parameters
            continue
        }
		matches := re.FindStringSubmatch(param)
		if len(matches) == 0 {
			log.Printf("Invalid parameter format: %s", param)
			continue
		}

		// Map regex groups for easier access
		groupNames := re.SubexpNames()
		groups := map[string]string{}
		for i, name := range groupNames {
			if i != 0 && name != "" { // Ignore the full match at index 0
				groups[name] = matches[i]
			}
		}

		column := groups["key"]
		rawOperator := groups["operator"]
		value := groups["value"]

		log.Printf("Processing parameter: %s, column: %s, operator: %s, value: %s", param, column, rawOperator, value)

		if column == "" || rawOperator == "" || value == "" {
			log.Printf("Missing key, operator, or value in parameter: %s", param)
			continue
		}

		// Parse the operator
		operator, err := parseOperator(rawOperator)
		if err != nil {
			log.Printf("Failed to parse operator: %v", err)
			continue
		}

		// Add SQL fragment
		sqlQuery.WriteString(fmt.Sprintf(` AND "%s" %s $%d`, column, operator, valueIndex))
		*values = append(*values, value)
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
