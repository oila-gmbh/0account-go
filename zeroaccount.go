package zeroaccount

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Setter func(ctx context.Context, k string, v []byte) error
type Getter func(ctx context.Context, k string) ([]byte, error)
type ErrorListener func(ctx context.Context, err error)

type Engine interface {
	Set(ctx context.Context, k string, v []byte) error
	Get(ctx context.Context, k string) ([]byte, error)
}

// Client struct to handle middleware logic
type Client struct {
	Engine             Engine
	GetterSetterEngine *GetterSetterEngine
	Client             *http.Client
	ErrorListener      ErrorListener
}

func httpClient() *http.Client {
	var netTransport = &http.Transport{
		DialContext: (net.Dialer{
			Timeout: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return netClient
}

// New returns an instances of Client middleware
func New(options ...option) *Client {
	zero := Client{}
	for _, option := range options {
		option(&zero)
	}
	if zero.Engine == nil && zero.GetterSetterEngine == nil {
		zero.Engine = NewInMemoryEngine()
	} else if zero.Engine == nil {
		zero.Engine = zero.GetterSetterEngine
	}

	if zero.Client == nil {
		zero.Client = httpClient()
	}

	return &zero
}

// Auth handles the authentication
func (zero *Client) Auth(ctx context.Context, headers map[string]string, body []byte) ([]byte, error) {
	if zero == nil || zero.Engine == nil {
		return nil, zero.error(ctx, fmt.Errorf("engine is not provided and/or the library is not initialised"))
	}
	uuid := headers["0account-uuid"]
	token, err := BearerFromHeader(headers)

	if err != nil || token == "" {
		if body == nil {
			return nil, zero.error(ctx, fmt.Errorf("body cannot be nil"))
		}
		err := zero.save(ctx, uuid, body)
		if err != nil {
			return nil, zero.error(ctx, err)
		}
		return nil, nil
	}

	data, err := zero.authorize(ctx, token, uuid)
	if err != nil {
		return nil, zero.error(ctx, fmt.Errorf("cannot authorise: %v", err))
	}
	return data, nil
}

func (zero *Client) save(ctx context.Context, uuid string, data []byte) error {
	err := zero.Engine.Set(ctx, uuid, data)
	if err != nil {
		return fmt.Errorf("engine error: cannot set. err: %v, key: %s, value: %s", err, uuid, string(data))
	}
	return nil
}

func (zero *Client) authorize(ctx context.Context, token, uuid string) ([]byte, error) {
	if token == "" {
		return nil, fmt.Errorf("empty or wrong bearer token")
	}
	if uuid == "" {
		return nil, fmt.Errorf("uuid is not provided")
	}

	v, err := zero.Engine.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("engine error: key is not found, err: %v, key: %s, value:%s", err, uuid, string(v))
	}
	return v, nil
}

func (zero *Client) error(ctx context.Context, err error) error {
	if err != nil && zero.ErrorListener != nil {
		zero.ErrorListener(ctx, err)
	}
	return err
}
