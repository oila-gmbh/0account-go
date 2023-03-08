package zeroaccount

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

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
