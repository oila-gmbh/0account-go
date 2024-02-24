package zeroaccount

import (
	"context"
)

var (
	authHeaders = []string{"x-0account-auth", "X-0account-Auth", "X-0account-AUTH"}
	uuidHeaders = []string{"x-0account-uuid", "X-0account-Uuid", "X-0account-UUID"}
)

func getAuthHeader(headers map[string][]string) string {
	return getFromHeader(authHeaders, headers)
}

func getUUIDHeader(headers map[string][]string) string {
	return getFromHeader(uuidHeaders, headers)
}

func getFromHeader(keys []string, headers map[string][]string) string {
	for _, key := range keys {
		if result := firstFromHeaders(headers[key]); result != "" {
			return result
		}
	}
	return ""
}

func firstFromHeaders(value []string) string {
	if len(value) == 0 {
		return ""
	}
	return value[0]
}

func zerror(ctx context.Context, err error) error {
	if errorListener != nil {
		errorListener(ctx, err)
	}
	return err
}
