package zeroaccount

import (
	"context"
	"fmt"
	"io"
)

var (
	setter        Setter
	getter        Getter
	errorListener ErrorListener
)

func init() {
	SetEngine(NewInMemoryEngine())
}

// Auth handles the authentication
func Auth[T Header](ctx context.Context, headers map[string]T, body io.Reader) ([]byte, error) {
	if setter == nil || getter == nil {
		return nil, zerror(ctx, fmt.Errorf("engine is not set and/or the library is not initialised"))
	}
	uuid := headersToString(headers["x-0account-uuid"])
	token, err := bearerFromHeader(headersToString(headers["Authorization"]))

	if err != nil || token == "" {
		bytes := copyBody(body)
		if bytes == nil {
			return nil, zerror(ctx, fmt.Errorf("body cannot be nil"))
		}
		err := save(ctx, uuid, bytes)
		if err != nil {
			return nil, zerror(ctx, err)
		}
		return nil, nil
	}

	data, err := authorize(ctx, token, uuid)
	if err != nil {
		return nil, zerror(ctx, fmt.Errorf("cannot authorise: %v", err))
	}
	return data, nil
}

func save(ctx context.Context, uuid string, data []byte) error {
	err := setter(ctx, uuid, data)
	if err != nil {
		return fmt.Errorf("engine error: cannot set. err: %v, key: %s, value: %s", err, uuid, string(data))
	}
	return nil
}

func authorize(ctx context.Context, token, uuid string) ([]byte, error) {
	if token == "" {
		return nil, fmt.Errorf("empty or wrong bearer token")
	}
	if uuid == "" {
		return nil, fmt.Errorf("uuid is not provided")
	}

	v, err := getter(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("engine error: key is not found, err: %v, key: %s, value:%s", err, uuid, string(v))
	}
	return v, nil
}
