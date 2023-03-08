package zeroaccount

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

const appSecretKey = "appSecret"
const x0accountUUIDKey = "x-0account-uuid"
const x0accountAuthKey = "x-0account-auth"

var (
	appSecret     string
	setter        Setter
	getter        Getter
	errorListener ErrorListener
)

func init() {
	SetEngine(NewInMemoryEngine())
}

// Auth handles the authentication
func Auth[T Header](ctx context.Context, headers map[string]T, body io.Reader) ([]byte, error) {
	if appSecret == "" {
		return nil, zerror(ctx, fmt.Errorf("app secret is not set"))
	}
	if setter == nil || getter == nil {
		return nil, zerror(ctx, fmt.Errorf("engine is not set and/or the library is not initialised"))
	}
	uuid := headersToString(headers[x0accountUUIDKey])
	authenticating := headersToString(headers[x0accountAuthKey]) == "true"

	if !authenticating {
		bytes, err := copyBody(body)
		if bytes == nil || err != nil {
			return nil, zerror(ctx, fmt.Errorf("error getting data from body"))
		}
		data := map[string]string{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			return nil, zerror(ctx, fmt.Errorf("cannot unmarshal data"))
		}
		if data[appSecretKey] != appSecret {
			return nil, zerror(ctx, fmt.Errorf("incorrect app secret"))
		}
		if err = save(ctx, uuid, bytes); err != nil {
			return nil, zerror(ctx, err)
		}
		return nil, nil
	}

	data, err := authorize(ctx, uuid)
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

func authorize(ctx context.Context, uuid string) ([]byte, error) {
	if uuid == "" {
		return nil, fmt.Errorf("uuid is not provided")
	}

	v, err := getter(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("engine error: key is not found, err: %v, key: %s, value:%s", err, uuid, string(v))
	}
	return v, nil
}
