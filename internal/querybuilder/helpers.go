package querybuilder

import (
	"fmt"
	"log"
	"strings"
)

// // Handle date parameters (e.g., start, end)
// func handleDateParam(param string, _ interface{}, valueIndex int) string {
// 	if param == "start" {
// 		return fmt.Sprintf(` AND "DateVisited" >= $%d`, valueIndex)
// 	} else if param == "end" {
// 		return fmt.Sprintf(` AND "DateVisited" < $%d`, valueIndex)
// 	}
// 	return ""
// }

// // Handle standard key-value pairs
// func handleStandardParam(key string, _ interface{}, valueIndex int) string {
// 	return fmt.Sprintf(` AND "%s" = $%d`, key, valueIndex)
// }
 

func AddLimitOffsetToBuilder(sqlQuery *strings.Builder, queryParams map[string]interface{}) {
    // Handle LIMIT
    if limit, ok := queryParams["limit"]; ok {
        log.Printf("Adding LIMIT: %v", limit)
        sqlQuery.WriteString(fmt.Sprintf(" LIMIT %v", limit))
    }

    // Handle OFFSET
    if offset, ok := queryParams["offset"]; ok {
        log.Printf("Adding OFFSET: %v", offset)
        sqlQuery.WriteString(fmt.Sprintf(" OFFSET %v", offset))
    }
}