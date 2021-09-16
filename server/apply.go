package server

import (
	"errors"
	"github.com/danielslee/gollab"
)

// ErrUnknownRevision is an error indicating that an unknown revision was encountered.
var ErrUnknownRevision = errors.New("unknown revision")

// ErrInvalidOperation is an error indication that the operation cannot be applied.
var ErrInvalidOperation = errors.New("invalid operation")

// ApplyClientOpInput contains data needed to apply a client operation. It is used as the input to ApplyClientOp.
type ApplyClientOpInput struct {
	CurrentDocument gollab.TokenArray
	CurrentRevision int

	Op           gollab.CompositeOp
	TransformOps []gollab.CompositeOp
}

// ApplyClientOpOutput serves as the output of ApplyClientOp.
type ApplyClientOpOutput struct {
	Document gollab.TokenArray
	Op       gollab.CompositeOp
	Revision int
}

// ApplyClientOp Applies a client operation. This function is intended to be used by a StateStore implementation.
func ApplyClientOp(i ApplyClientOpInput) (o ApplyClientOpOutput, err error) {
	o.Op = i.Op
	for _, transformOp := range i.TransformOps {
		if o.Op.InputLength() != transformOp.InputLength() {
			err = ErrInvalidOperation
			return
		}
		o.Op, _ = o.Op.Transform(transformOp)
	}

	writer := i.CurrentDocument.Type().NewBuilder()
	err = o.Op.Apply(gollab.NewTokenArrayReader(i.CurrentDocument), writer)
	if err != nil {
		return
	}

	o.Document = writer.TokenArray()
	o.Revision = i.CurrentRevision + 1

	return
}
