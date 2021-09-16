package gollab

import (
	"errors"
)

// Op represents any operation which can be applied.
type Op interface {
	InputLength() int
	OutputLength() int
	Apply(reader TokenReader, writer TokenWriter) error
}

// PrimitiveOp represents a primitive operation, which can be composed and transformed with other PrimitiveOps, and in
// addition sliced.
//
// There are four PrimitiveOps: NoOp, Retain, Delete and Insert.
type PrimitiveOp interface {
	Op
	Slice(start, end int) PrimitiveOp
	Compose(b PrimitiveOp) PrimitiveOp
	Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp)
}

// ErrLengthMismatch length mismatch error
var ErrLengthMismatch = errors.New("length mismatch")

// ErrUnexpectedOp unexpected operation error
var ErrUnexpectedOp = errors.New("unexpected operation")

// ErrInvalidSlice invalid slice error
var ErrInvalidSlice = errors.New("invalid slice")

func checkComposeLength(a, b Op) {
	if a.OutputLength() != b.InputLength() {
		panic(ErrLengthMismatch)
	}
}

func checkTransformLength(a, b Op) {
	if a.InputLength() != b.InputLength() {
		panic(ErrLengthMismatch)
	}
}

func checkSliceValidity(start, end int) {
	if end < start {
		panic(ErrInvalidSlice)
	}
}
