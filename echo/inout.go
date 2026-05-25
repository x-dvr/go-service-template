package echo

import (
	"github.com/0xde86/go-service-template/core"
)

// In is an input DTO of echo module
type In struct {
	From      string
	Data      string
	UseCached bool
	WithNoise bool
}

func (in In) Validate() error {
	if in.From == "" {
		return core.NewError(core.ErrValidation, nil).WithContext("From field is not set")
	}
	if in.Data == "" {
		return core.NewError(core.ErrValidation, nil).WithContext("Data field is not set")
	}
	return nil
}

// Out is an output DTO of echo module
type Out struct {
	Message string
}
