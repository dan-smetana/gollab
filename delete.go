package gollab

import (
	"fmt"
)

// Delete represents a delete operation which consumes a given number of tokens from the input, and writes nothing
// to the output.
type Delete struct {
	Count int
}

// InputLength returns the required input length.
//
// For a Delete operation, this equals its Count property.
func (d Delete) InputLength() int {
	return d.Count
}

// OutputLength returns the output length.
//
// This will always be zero for a Delete operation.
func (d Delete) OutputLength() int {
	return 0
}

// Slice slices the operation.
func (d Delete) Slice(start, end int) PrimitiveOp {
	checkSliceValidity(start, end)
	return Delete{Count: end - start}
}

// Apply applies the operation.
func (d Delete) Apply(reader TokenReader, _ TokenWriter) error {
	for i := 0; i < d.Count; i++ {
		if _, err := reader.ReadToken(); err != nil {
			return err
		}
	}
	return nil
}

// Join joins the operation with another one if possible. This is used for normalization.
func (d Delete) Join(next PrimitiveOp) PrimitiveOp {
	if nextDelete, ok := next.(Delete); ok {
		return Delete{Count: d.Count + nextDelete.Count}
	}
	return nil
}

// Swap returns true if it should be swapped with the next operation. This is used for normalization.
func (d Delete) Swap(next PrimitiveOp) bool {
	if _, ok := next.(Insert); ok {
		return true
	}
	return false
}

// Compose attempts to compose the delete operation with another PrimitiveOp.
//
// The output length of the first operation has to equal the input length of the other one.
func (d Delete) Compose(other PrimitiveOp) PrimitiveOp {
	checkComposeLength(d, other)

	switch other.(type) {
	case NoOp:
		return d
	default:
		panic(ErrUnexpectedOp)
	}
}

// Transform attempts to perform Operation Transformation.
//
// The input length of both operations has to be equal.
func (d Delete) Transform(other PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(d, other)

	switch b := other.(type) {
	case Delete:
		return NoOp{}, NoOp{}
	case Retain:
		return Delete{Count: b.Count}, NoOp{}
	default:
		panic(ErrUnexpectedOp)
	}
}

// String provides a string representation of the Delete operation.
func (d Delete) String() string {
	return fmt.Sprintf("Delete(%d)", d.Count)
}
