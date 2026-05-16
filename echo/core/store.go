package core

import "context"

type Store interface {
	Save(context.Context, Echo) error
	Load(ctx context.Context, from string) (*Echo, error)
}
