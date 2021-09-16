package gollab

import (
	"fmt"
)

// Insert represents an insert operation which consumes nothing and writes the content of its Tokens TokenArray to the
// output.
type Insert struct {
	Tokens TokenArray
}

// InputLength returns the required input length.
//
// For an Insert operation, this is always zero.
func (i Insert) InputLength() int {
	return 0
}

// OutputLength returns the length of the result after the operation has been applied.
//
// For an Insert operation, this will equal the length of its Tokens.
func (i Insert) OutputLength() int {
	return i.Tokens.Len()
}

// Slice slices the operation.
func (i Insert) Slice(start, end int) PrimitiveOp {
	checkSliceValidity(start, end)
	return Insert{Tokens: i.Tokens.Slice(start, end)}
}

// Apply applies the operation.
func (i Insert) Apply(_ TokenReader, writer TokenWriter) error {
	for idx := 0; idx < i.Tokens.Len(); idx++ {
		token := i.Tokens.At(idx)
		if err := writer.WriteToken(token); err != nil {
			return err
		}
	}
	return nil
}

// Join joins the operation with another one if possible. This is used for normalization.
func (i Insert) Join(next PrimitiveOp) PrimitiveOp {
	if nextInsert, ok := next.(Insert); ok {
		return Insert{Tokens: i.Tokens.Type().Concat(i.Tokens, nextInsert.Tokens)}
	}
	return nil
}

// Compose attempts to compose the delete operation with another PrimitiveOp.
//
// The output length of the first operation has to equal the input length of the other one.
func (i Insert) Compose(b PrimitiveOp) PrimitiveOp {
	checkComposeLength(i, b)

	switch b.(type) {
	case Delete:
		return NoOp{}
	case Retain:
		return i
	default:
		panic(ErrUnexpectedOp)
	}
}

// Transform attempts to perform Operation Transformation.
//
// The input length of both operations has to be equal.
func (i Insert) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	checkTransformLength(i, b)

	switch b.(type) {
	case NoOp:
		return i, Retain{Count: i.OutputLength()}
	default:
		panic(ErrUnexpectedOp)
	}
}

// String provides a string representation of the Insert operation.
func (i Insert) String() string {
	return fmt.Sprintf("Insert(%s)", i.Tokens)
}
