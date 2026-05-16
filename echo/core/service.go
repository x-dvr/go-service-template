package core

import (
	"context"
	"fmt"
	"math/rand/v2"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Echo(ctx context.Context, in EchoIn) (EchoOut, error) {
	if in.UseCached {
		e, err := s.store.Load(ctx, in.From)
		if err != nil {
			return EchoOut{}, fmt.Errorf("load from cache for %q: %w", in.From, err)
		}
		return EchoOut{
			Message: echo(e.Data, in.WithNoise),
		}, nil
	}
	err := s.store.Save(ctx, Echo{
		From: in.From,
		Data: in.Data,
	})
	if err != nil {
		return EchoOut{}, fmt.Errorf("save to cache for %q: %w", in.From, err)
	}
	return EchoOut{
		Message: echo(in.Data, in.WithNoise),
	}, nil
}

func echo(msg string, withNoise bool) string {
	if !withNoise {
		return msg
	}
	return fmt.Sprintf("%s (((%s))) %s", generateNoise(5), msg, generateNoise(5))
}

const noiseSet = "}=+#!)*-%:]&<>?'`\\/;,._{~|\"@(^$["

func generateNoise(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = noiseSet[rand.IntN(len(noiseSet))]
	}
	return string(b)
}
