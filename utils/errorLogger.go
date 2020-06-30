package utils

import "fmt"

// LogError ...
func LogError(message string, err error) {
	fmt.Print(message)
	fmt.Print(err)
}
