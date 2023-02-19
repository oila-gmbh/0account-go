package oneaccount

import (
	"context"
	"fmt"
	"net/http"
)

type option func(zero *ZeroAccount)

// SetOnErrorListener is used to track errors
func SetOnErrorListener(errorListener ErrorListener) func(zero *ZeroAccount) {
	return func(zero *ZeroAccount) {
		zero.ErrorListener = errorListener
	}
}

// SetEngine is used to change the cache engine.
// If not set default in memory cache engine is used
func SetEngine(e Engine) func(zero *ZeroAccount) {
	return func(zero *ZeroAccount) {
		zero.Engine = e
	}
}

// SetClient can be used to change the underlying http client
func SetClient(client *http.Client) func(zero *ZeroAccount) {
	return func(zero *ZeroAccount) {
		zero.Client = client
	}
}

// GetterSetterEngine is a convenience object to be used instead of
// creating the Engine
type GetterSetterEngine struct {
	Setter Setter
	Getter Getter
}

func (g GetterSetterEngine) Set(ctx context.Context, k string, v []byte) error {
	if g.Setter == nil {
		return fmt.Errorf("engine setter is not set")
	}
	return g.Setter(ctx, k, v)
}

func (g GetterSetterEngine) Get(ctx context.Context, k string) ([]byte, error) {
	if g.Getter == nil {
		return nil, fmt.Errorf("engine getter is not set")
	}
	return g.Getter(ctx, k)
}

// SetEngineSetter can be used to set the setter function
// for the engine
func SetEngineSetter(setter Setter) func(oa *ZeroAccount) {
	return func(oa *ZeroAccount) {
		oa.Engine = nil
		if oa.GetterSetterEngine == nil {
			oa.GetterSetterEngine = &GetterSetterEngine{}
		}
		oa.GetterSetterEngine.Setter = setter
	}
}

// SetEngineGetter can be used to set the getter function
// for the engine
func SetEngineGetter(getter Getter) func(oa *ZeroAccount) {
	return func(oa *ZeroAccount) {
		oa.Engine = nil
		if oa.GetterSetterEngine == nil {
			oa.GetterSetterEngine = &GetterSetterEngine{}
		}
		oa.GetterSetterEngine.Getter = getter
	}
}
