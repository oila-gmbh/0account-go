package zeroaccount

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const expireTimeDuration = 1 * time.Minute

// authorizingUser represents an in memory object for authorising users
type authorizingUser struct {
	expiresAt time.Time
	data      []byte
}

// InMemoryEngine is a cache in-memory engine
type InMemoryEngine struct {
	authorizingUsers map[string]authorizingUser
	mu               sync.RWMutex
}

// NewInMemoryEngine return an instance of InMemoryEngine
func NewInMemoryEngine() *InMemoryEngine {
	engine := &InMemoryEngine{}
	engine.authorizingUsers = make(map[string]authorizingUser)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				engine.mu.Lock()
				now := time.Now()
				for k, v := range engine.authorizingUsers {
					if v.expiresAt.Before(now) {
						delete(engine.authorizingUsers, k)
					}
				}
				engine.mu.Unlock()
			}
		}
	}()
	return engine
}

// Set stores an authorising user in memory
func (engine *InMemoryEngine) Set(ctx context.Context, k string, v []byte) error {
	// we don't need a sophisticated way to handle context here, so we just check
	// if it is already cancelled and return, otherwise ignore the context and proceed
	select {
	case <-ctx.Done():
		return nil
	default:
		engine.mu.Lock()
		engine.authorizingUsers[k] = authorizingUser{
			expiresAt: time.Now().Add(expireTimeDuration),
			data:      v,
		}
		engine.mu.Unlock()
	}

	return nil
}

// Get retrieves an authorising user from memory
func (engine *InMemoryEngine) Get(ctx context.Context, k string) ([]byte, error) {
	// we don't need a sophisticated way to handle context here, so we just check
	// if it is already cancelled and return, otherwise ignore the context and proceed
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		engine.mu.RLock()
		v, ok := engine.authorizingUsers[k]
		engine.mu.RUnlock()
		if !ok || v.expiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("no item found or item expired for key: %s", k)
		}
		defer func() {
			engine.mu.Lock()
			delete(engine.authorizingUsers, k)
			engine.mu.Unlock()
		}()
		return v.data, nil
	}
}
