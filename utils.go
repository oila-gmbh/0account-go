package zeroaccount

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

var (
	authHeaders = []string{"x-0account-auth", "X-0account-Auth", "X-0account-AUTH"}
	uuidHeaders = []string{"x-0account-uuid", "X-0account-Uuid", "X-0account-UUID"}
)

func getAuthHeader[T Header](headers map[string]T) string {
	return getFromHeader(authHeaders, headers)
}

func getUUIDHeader[T Header](headers map[string]T) string {
	return getFromHeader(uuidHeaders, headers)
}

func getFromHeader[T Header](keys []string, headers map[string]T) string {
	for _, key := range keys {
		if result := headersToString(headers[key]); result != "" {
			return result
		}
	}
	return ""
}

func headersToString[T Header](value T) string {
	switch s := any(value).(type) {
	case []string:
		return s[0]
	case string:
		return s
	default:
		return fmt.Sprint(s)
	}
}

func copyBody(body io.Reader) ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}
	data := bytes.Clone(buf.Bytes())
	body = io.NopCloser(&buf)
	return data, nil
}

func zerror(ctx context.Context, err error) error {
	if errorListener != nil {
		errorListener(ctx, err)
	}
	return err
}
