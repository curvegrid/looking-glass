package customError

import "fmt"

// NewAPICallError returns a new error representing a failed API call
func NewAPICallError(endpoint string, statusCode int, message string) error {
	return fmt.Errorf("API call to %s failed with status: %d, message: %s", endpoint, statusCode, message)
}
