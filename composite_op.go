package gollab

// CompositeOp is an operation composed of multiple PrimitiveOps.
//
// CompositeOp implements MarshalJSON and UnmarshalJSON allowing it to be encoded to/decoded from JSON as an array of
// objects each carrying their type and depending on the type, count or text properties. Here is an example containing
// all four possible types:
// 	[
//	 {"type": "noop"},
//	 {"type": "retain", "count": 1},
//	 {"type": "delete", "count": 1},
//	 {"type": "insert", "tokens": ...}
//	]
// The tokens in an insert operation will be serialized based on its MarshalJSON and UnmarshalJSON methods.
type CompositeOp []PrimitiveOp

// NewCompositeOp creates a new CompositeOp, normalizing the given PrimitiveOps as needed. (i.e. multiple consecutive
// operations of the same type will be merged into one)
func NewCompositeOp(ops ...PrimitiveOp) CompositeOp {
	return normalize(ops)
}

// InputLength returns the required input length calculated as the sum of all contained operation's InputLength.
func (c CompositeOp) InputLength() (length int) {
	for _, p := range c {
		length += p.InputLength()
	}
	return
}

// OutputLength returns the length of an output this operation produces calculated as the sum of all contained
// operation's OutputLength.
func (c CompositeOp) OutputLength() (length int) {
	for _, p := range c {
		length += p.OutputLength()
	}
	return
}

// Apply applies the operation reading from an TokenReader and writing to an TokenWriter.
func (c CompositeOp) Apply(reader TokenReader, writer TokenWriter) error {
	for _, p := range c {
		if err := p.Apply(reader, writer); err != nil {
			return err
		}
	}
	return nil
}

// Transform implements OT - Operation Transformation.
//
// Suppose we have two operations which were applied simultaneously.
// We need to ensure we end up with the same result, regardless of the order we apply them in.
// This is where OT comes in.
//
// Gollab implements two operations that can be applied to composite operations such as `a` and `b` above: Transform and
// Compose. Operation Transformation, as suggested by its name, refers to the former.
//
// OT is defined as `OT(a, b) -> a', b'`, where applying `a` first and `b'` second yields the same result as applying
// `b` first and `a'` second.
//
// Let's apply OT to `a` and `b` above. We end up with two new operations: `a'` and `b'`
//
//  a' := Insert("H"), Delete(1), Retain(4), Insert(", World"), Retain(1)
//  b' := Retain(12), Insert("!")
// We can verify that the above transformation is correct by applying `a` to `hello`, resulting in `Hello, World`, which
// gives us `Hello, World!` after applying `b'` to it. We arrive to the same result by applying `b` to hello getting
// `hello!` followed by applying `a'` making it into `Hello, World!`.
func (c CompositeOp) Transform(b CompositeOp) (aPrime, bPrime CompositeOp) {
	slicedA, slicedB := slice(c, b,
		inputLengthFunc, inputLengthFunc)

	for i := range slicedA {
		aOp, bOp := slicedA[i], slicedB[i]
		aOpPrime, bOpPrime := aOp.Transform(bOp)
		aPrime = append(aPrime, aOpPrime)
		bPrime = append(bPrime, bOpPrime)
	}
	return NewCompositeOp(aPrime...), NewCompositeOp(bPrime...)
}

// Compose composes two operations which happened in order into one.
//
// See also: func Compose which takes any number of operations as opposed to two and the attached example.
func (c CompositeOp) Compose(b CompositeOp) CompositeOp {
	slicedA, slicedB := slice(c, b,
		outputLengthFunc, inputLengthFunc)

	var res []PrimitiveOp
	for i := range slicedA {
		aOp, bOp := slicedA[i], slicedB[i]
		c := aOp.Compose(bOp)
		res = append(res, c)
	}
	return NewCompositeOp(res...)
}

// Compose composes multiple operations which happened in order into one.
//
// When building a collaborative editor we may want to combine two or more composite operations which were
// applied in order into one.
func Compose(ops ...CompositeOp) CompositeOp {
	if l := len(ops); l == 0 {
		return CompositeOp{}
	} else if l < 2 {
		return ops[0]
	}

	op := ops[0]
	for i := 1; i < len(ops); i++ {
		op = op.Compose(ops[i])
	}
	return op
}
