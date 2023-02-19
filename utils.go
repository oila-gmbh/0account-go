package zeroaccount

import (
	"fmt"
	"strings"
)

// BearerFromHeader method
// retrieves token from Authorization header
func BearerFromHeader(headers map[string]string) (string, error) {

	auth := headers["Authorization"]
	const prefix = "BEARER "

	if !(len(auth) >= len(prefix) && strings.ToUpper(auth[0:len(prefix)]) == prefix) {
		return "", fmt.Errorf("token is not found")
	}

	t := auth[len(prefix):]

	return t, nil
}
