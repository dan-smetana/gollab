package gollab

import (
	"fmt"
)

// Retain represents a retain operation which copies a given number of tokens (Count) from the input over to the
// output
type Retain struct {
	Count int
}

// InputLength returns the required input length.
//
// For a Retain operation, this equals its Count property.
func (r Retain) InputLength() int {
	return r.Count
}

// OutputLength returns the output length.
//
// For a Retain operation, this equals its Count property.
func (r Retain) OutputLength() int {
	return r.Count
}

// Slice slices the operation.
func (r Retain) Slice(start, end int) PrimitiveOp {
	checkSliceValidity(start, end)
	return Retain{Count: end - start}
}

// Apply applies the operation
func (r Retain) Apply(reader TokenReader, writer TokenWriter) error {
	for i := 0; i < r.Count; i++ {
		token, err := reader.ReadToken()
		if err != nil {
			return err
		}
		if err := writer.WriteToken(token); err != nil {
			return err
		}
	}
	return nil
}

// Join joins the operation with another one if possible. This is used for normalization.
func (r Retain) Join(next PrimitiveOp) PrimitiveOp {
	if nextRetain, ok := next.(Retain); ok {
		return Retain{Count: r.Count + nextRetain.Count}
	}
	return nil
}

// Compose attempts to compose the delete operation with another PrimitiveOp.
//
// The output length of the first operation has to equal the input length of the other one.
func (r Retain) Compose(b PrimitiveOp) PrimitiveOp {
	checkComposeLength(r, b)

	switch b := b.(type) {
	case Retain:
		return r
	case Delete:
		return b
	default:
		panic(ErrUnexpectedOp)
	}
}

// Transform attempts to perform Operation Transformation.
//
// The input length of both operations has to be equal.
func (r Retain) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(r, b)

	switch b := b.(type) {
	case Retain:
		return b, b
	case Delete:
		return NoOp{}, b
	default:
		panic(ErrUnexpectedOp)
	}
}

// String provides a string representation of the Retain operation.
func (r Retain) String() string {
	return fmt.Sprintf("Retain(%d)", r.Count)
}
