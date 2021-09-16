/*
Package runetoken provides an implementation of the gollab.TokenReader, gollab.TokenWriter and gollab.TokenArray
for plain unicode strings.
 */
package runetoken

import (
	"encoding/json"
	"errors"
	"github.com/danielslee/gollab"
	"strings"
)

// Array implements the TokenArray interface as a rune slice.
type Array []rune

// Type returns the ArrayType.
func (Array) Type() gollab.TokenArrayType {
	return ArrayType{}
}

// ArrayType contains methods related to Array.
type ArrayType struct{}

// NewBuilder creates a new ArrayBuilder.
func (ArrayType) NewBuilder() gollab.TokenArrayBuilder {
	return &ArrayBuilder{
		Array: []rune{},
	}
}

// Concat concatenates two TokenArrays and returns the result.
func (ArrayType) Concat(a, b gollab.TokenArray) gollab.TokenArray {
	rta, ok := a.(Array)
	if !ok {
		panic("ArrayType.Concat: expected param a to be of type Array")
	}
	rtb, ok := b.(Array)
	if !ok {
		panic("ArrayType.Concat: expected param b to be of type Array")
	}

	newTokens := Array(make([]rune, rta.Len()+rtb.Len()))
	copy(newTokens, rta)
	copy(newTokens[rta.Len():], rtb)
	return newTokens
}

// ArrayBuilder is a TokenArrayBuilder implementation using Array.
type ArrayBuilder struct {
	Array Array
}

// WriteToken appends a given token to the builder's Array.
func (b *ArrayBuilder) WriteToken(token interface{}) error {
	r, ok := token.(rune)
	if !ok {
		return errors.New("ArrayBuilder.WriteToken: expected a rune")
	}
	b.Array = append(b.Array, r)
	return nil
}

// TokenArray returns the Array.
func (b ArrayBuilder) TokenArray() gollab.TokenArray {
	return b.Array
}

// At returns an element given its index.
func (t Array) At(idx int) interface{} {
	return t[idx]
}

// Slice slices the array.
func (t Array) Slice(start, end int) gollab.TokenArray {
	return t[start:end]
}

// Len returns the length of the array.
func (t Array) Len() int {
	return len(t)
}

// String returns the Array as a string.
func (t Array) String() string {
	return string(t)
}

// MarshalJSON encodes the Array as a JSON string.
func (t Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

// UnmarshalJSON parses the Array from a JSON string.
func (t *Array) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*t = []rune(str)
	return nil
}

// StringReader implements a TokenReader using a strings.Reader, which can be created using a plain Go string.
type StringReader struct {
	Reader *strings.Reader
}

// ReadToken reads and returns a rune using the internal strings.Reader.
func (r StringReader) ReadToken() (interface{}, error) {
	res, _, err := r.Reader.ReadRune()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// StringWriter implements the TokenWriter interface using a strings.Builder, making it easy to get a Go string out.
type StringWriter struct {
	strings.Builder
}

// WriteToken writes a rune token
func (w *StringWriter) WriteToken(token interface{}) error {
	if r, ok := token.(rune); ok {
		_, err := w.WriteRune(r)
		return err
	}
	return errors.New("StringWriter.WriteToken: expected a rune")
}

// ApplyToString applies the operation to a plain Go string, returning the result as a plain Go string.
// This assumes the operation consumes and outputs runes only.
func ApplyToString(op gollab.Op, text string) (string, error) {
	if op.InputLength() != len([]rune(text)) {
		return "", gollab.ErrLengthMismatch
	}
	reader := StringReader{Reader: strings.NewReader(text)}
	var out StringWriter
	if err := op.Apply(reader, &out); err != nil {
		return "", err
	}
	return out.String(), nil
}
