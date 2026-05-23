package echo

import "context"

// Store is a port for storaga/repository required by echo module
type Store interface {
	Save(context.Context, Echo) error
	Load(ctx context.Context, from string) (*Echo, error)
}
