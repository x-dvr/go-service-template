package core

import (
	"github.com/x-dvr/go-service-template/core"
)

type EchoIn struct {
	From      string
	Data      string
	UseCached bool
	WithNoise bool
}

func (in EchoIn) Validate() error {
	if in.From == "" {
		return core.NewError(core.ErrValidation, nil).WithContext("From field is not set")
	}
	if in.Data == "" {
		return core.NewError(core.ErrValidation, nil).WithContext("Data field is not set")
	}
	return nil
}

type EchoOut struct {
	Message string
}
