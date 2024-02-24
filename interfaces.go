package zeroaccount

import "context"

type Setter func(ctx context.Context, k string, v []byte) error
type Getter func(ctx context.Context, k string) ([]byte, error)
type ErrorListener func(ctx context.Context, err error)

type Engine interface {
	Set(ctx context.Context, k string, v []byte) error
	Get(ctx context.Context, k string) ([]byte, error)
}
