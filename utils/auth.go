package utils

import (
	"encoding/base64"
	"fmt"
)

// CreateBasicAuth creates a Basic Authentication header
func CreateBasicAuth(username, password string) string {
	credentials := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}
