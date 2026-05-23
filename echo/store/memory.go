package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/x-dvr/go-service-template/core"
	"github.com/x-dvr/go-service-template/echo"
	"github.com/x-dvr/go-service-template/logs"
)

type MemoryStore struct {
	store map[string]echo.Echo
}

func NewInMemory() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]echo.Echo),
	}
}

// Load implements [ports.Store].
func (m *MemoryStore) Load(ctx context.Context, from string) (*echo.Echo, error) {
	e, err := m.loadInternal(ctx, from)
	if errors.Is(err, errNotFound) {
		return nil, core.NewError(core.ErrNotFound, err)
	}
	if err != nil {
		return nil, core.ErrorFrom(err)
	}
	return e, nil
}

// Save implements [ports.Store].
func (m *MemoryStore) Save(ctx context.Context, e echo.Echo) error {
	_, exists := m.store[e.From]
	if exists {
		return core.
			NewError(core.ErrDuplicate, nil).
			WithContext(fmt.Sprintf("already exist: %s", e.From))
	}

	return m.saveInternal(ctx, e)
}

var _ echo.Store = (*MemoryStore)(nil)

func (m *MemoryStore) saveInternal(ctx context.Context, e echo.Echo) error {
	start := time.Now()
	time.Sleep(time.Duration(rand.IntN(100) * int(time.Millisecond)))
	shouldFail := rand.IntN(100) > 70
	if shouldFail {
		return errors.New("failed to save")
	}

	m.store[e.From] = e
	logs.AddAttrs(ctx, slog.Duration("db_save_duration", time.Since(start)))
	return nil
}

var errNotFound = errors.New("not found in memory")

func (m *MemoryStore) loadInternal(ctx context.Context, from string) (*echo.Echo, error) {
	start := time.Now()
	time.Sleep(time.Duration(rand.IntN(100) * int(time.Millisecond)))
	shouldFail := rand.IntN(100) > 90
	if shouldFail {
		return nil, errors.New("failed to load")
	}

	e, found := m.store[from]
	if !found {
		return nil, fmt.Errorf("entity %s: %w", from, errNotFound)
	}

	logs.AddAttrs(ctx, slog.Duration("db_load_duration", time.Since(start)))
	return &e, nil
}
