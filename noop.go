package gollab

// NoOp represents a NoOp operation which preforms nothing.
type NoOp struct{}

// InputLength returns the required input length. This is always zero.
func (n NoOp) InputLength() int {
	return 0
}

// OutputLength returns the resulting output length. This is always zero.
func (n NoOp) OutputLength() int {
	return 0
}

// Slice slices the operation, returning a new empty NoOp.
func (n NoOp) Slice(_, _ int) PrimitiveOp {
	return NoOp{}
}

// Apply applies the operation, doing nothing.
func (n NoOp) Apply(TokenReader, TokenWriter) error {
	return nil
}

// Compose attempts to compose the delete operation with another PrimitiveOp.
//
// The output length of the first operation has to equal the input length of the other one.
func (n NoOp) Compose(b PrimitiveOp) PrimitiveOp {
	switch b := b.(type) {
	case Insert:
		return b
	default:
		panic(ErrUnexpectedOp)
	}
}

// Transform attempts to perform Operation Transformation.
//
// The input length of both operations has to be equal.
func (n NoOp) Transform(b PrimitiveOp) (aPrime, bPrime PrimitiveOp) {
	switch b := b.(type) {
	case Insert:
		return Retain{Count: b.OutputLength()}, b
	default:
		panic(ErrUnexpectedOp)
	}
}

// String provides a string representation of the NoOp operation.
func (n NoOp) String() string {
	return "NoOp"
}
