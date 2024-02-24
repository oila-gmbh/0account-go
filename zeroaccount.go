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
	if setter == nil && getter == nil {
		SetEngine(NewInMemoryEngine())
	}
}

type ZeroRequest[T any] struct {
	Metadata metadata `json:"metadata"`
	Data     T        `json:"data"`
}

type metadata struct {
	AppSecret string `json:"appSecret,omitempty"`
	UserID    string `json:"userId,omitempty"`
	ProfileID string `json:"profileId,omitempty"`
}

// Auth handles the authentication
func Auth[T any](ctx context.Context, headers map[string][]string, body io.Reader) (T, Metadata, error) {
	zr := ZeroRequest[T]{}
	meta := Metadata{}
	if appSecret == "" {
		return zr.Data, meta, zerror(ctx, fmt.Errorf("app secret is not set"))
	}
	if setter == nil || getter == nil {
		return zr.Data, meta, zerror(ctx, fmt.Errorf("engine is not set and/or the library is not initialised"))
	}
	uuid := getUUIDHeader(headers)
	authenticating := getAuthHeader(headers) == "true"

	if !authenticating {
		bytes, err := io.ReadAll(body)
		if err != nil {
			return zr.Data, meta, zerror(ctx, fmt.Errorf("error getting data from body: %w", err))
		}
		if bytes == nil || len(bytes) == 0 {
			return zr.Data, meta, zerror(ctx, fmt.Errorf("error getting data from body"))
		}

		if err := json.Unmarshal(bytes, &zr); err != nil {
			return zr.Data, meta, zerror(ctx, fmt.Errorf("cannot unmarshal data %w", err))
		}
		if zr.Metadata.AppSecret != appSecret {
			return zr.Data, meta, zerror(ctx, fmt.Errorf("incorrect app secret"))
		}
		zr.Metadata.AppSecret = ""
		cleanData, err := json.Marshal(zr)
		if err != nil {
			return zr.Data, meta, zerror(ctx, fmt.Errorf("incorrect app secret, err: %w", err))
		}
		if err = save(ctx, uuid, cleanData); err != nil {
			return zr.Data, meta, zerror(ctx, err)
		}
		meta.userID = zr.Metadata.UserID
		meta.profileID = zr.Metadata.ProfileID
		return zr.Data, meta, nil
	}

	newZR, err := authorize[T](ctx, uuid)
	if err != nil {
		return newZR.Data, meta, zerror(ctx, fmt.Errorf("cannot authorise: %v", err))
	}
	meta.userID = newZR.Metadata.UserID
	meta.profileID = newZR.Metadata.ProfileID
	return newZR.Data, meta, nil
}

func save(ctx context.Context, uuid string, data []byte) error {
	if uuid == "" {
		return fmt.Errorf("uuid is not provided")
	}

	if err := setter(ctx, uuid, data); err != nil {
		return fmt.Errorf("engine error: cannot set. err: %v, key: %s, value: %s", err, uuid, string(data))
	}
	return nil
}

func authorize[T any](ctx context.Context, uuid string) (ZeroRequest[T], error) {
	zr := ZeroRequest[T]{}
	if uuid == "" {
		return zr, fmt.Errorf("uuid is not provided")
	}

	v, err := getter(ctx, uuid)
	if err != nil || v == nil {
		return zr, fmt.Errorf("engine error: key is not found, err: %v, key: %s, value:%s", err, uuid, string(v))
	}

	if err = json.Unmarshal(v, &zr); err != nil {
		return zr, fmt.Errorf("engine error: cannot unmarshal data, err: %v, key: %s, value:%s", err, uuid, string(v))
	}
	return zr, err
}
