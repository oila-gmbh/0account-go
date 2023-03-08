package zeroaccount

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

func bearerFromHeader(auth string) (string, error) {
	const prefix = "BEARER "

	if !(len(auth) >= len(prefix) && strings.ToUpper(auth[0:len(prefix)]) == prefix) {
		return "", fmt.Errorf("token is not found")
	}

	t := auth[len(prefix):]

	return t, nil
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
	_, err := buf.ReadFrom(body)
	if err != nil {
		return nil, err
	}
	body = io.NopCloser(&buf)
	return buf.Bytes(), nil
}

func zerror(ctx context.Context, err error) error {
	if errorListener != nil {
		errorListener(ctx, err)
	}
	return err
}
