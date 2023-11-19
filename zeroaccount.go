package zeroaccount

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

var (
	appSecret     string
	setter        Setter
	getter        Getter
	errorListener ErrorListener
)

func init() {
	SetEngine(NewInMemoryEngine())
}

type Data struct {
	Metadata struct {
		AppSecret string `json:"appSecret,omitempty"`
	} `json:"metadata"`
}

// Auth handles the authentication
func Auth[T Header](ctx context.Context, headers map[string]T, body io.Reader) ([]byte, error) {
	fmt.Println("---------- AUTH CALLED")
	if appSecret == "" {
		return nil, zerror(ctx, fmt.Errorf("app secret is not set"))
	}
	if setter == nil || getter == nil {
		return nil, zerror(ctx, fmt.Errorf("engine is not set and/or the library is not initialised"))
	}
	uuid := getUUIDHeader(headers)
	authenticating := getAuthHeader(headers) == "true"
	fmt.Println("uuid: ", uuid)
	fmt.Println("authenticating: ", authenticating)

	if !authenticating {
		fmt.Println("---------- SAVING 0")
		bytes, err := io.ReadAll(body)
		if err != nil {
			return nil, zerror(ctx, fmt.Errorf("error getting data from body: %w", err))
		}
		if bytes == nil || len(bytes) == 0 {
			return nil, zerror(ctx, fmt.Errorf("error getting data from body"))
		}
		data := Data{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			return nil, zerror(ctx, fmt.Errorf("cannot unmarshal data %w", err))
		}
		if data.Metadata.AppSecret != appSecret {
			return nil, zerror(ctx, fmt.Errorf("incorrect app secret"))
		}
		fmt.Println("---------- SAVING 1")
		if err = save(ctx, uuid, bytes); err != nil {
			fmt.Println("---------- SAVING ERROR: ", err)
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
	fmt.Println("---------- SAVING 2")
	if uuid == "" {
		return fmt.Errorf("uuid is not provided")
	}

	fmt.Println("---------- SAVING 3")
	if err := setter(ctx, uuid, data); err != nil {
		fmt.Println("---------- SAVING ERROR 2: ", err)
		return fmt.Errorf("engine error: cannot set. err: %v, key: %s, value: %s", err, uuid, string(data))
	}
	fmt.Println("---------- SAVED")
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
