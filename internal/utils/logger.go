// logger.go
package utils

import "log"

func LogError(err error) {
    if err != nil {
        log.Println("Error:", err)
    }
}