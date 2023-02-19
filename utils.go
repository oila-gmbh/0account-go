package zeroaccount

import (
	"fmt"
	"net/http"
	"strings"
)

// BearerFromHeader method
// retrieves token from Authorization header
func BearerFromHeader(header http.Header) (string, error) {

	auth := header.Get("Authorization")
	const prefix = "BEARER "

	if !(len(auth) >= len(prefix) && strings.ToUpper(auth[0:len(prefix)]) == prefix) {
		return "", fmt.Errorf("token is not found")
	}

	t := auth[len(prefix):]

	return t, nil
}
