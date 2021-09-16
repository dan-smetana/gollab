package gollab

import "io"

// TokenReader is an interface wrapping around a ReadToken method.
type TokenReader interface {
	ReadToken() (interface{}, error)
}

// TokenWriter is an interface wrapping around a WriteToken method.
type TokenWriter interface {
	WriteToken(interface{}) error
}

// TokenArray is an interface defining a read-only array of tokens.
type TokenArray interface {
	Type() TokenArrayType
	At(idx int) interface{}
	Slice(start, end int) TokenArray
	Len() int
}

// TokenArrayType is returned by a TokenArray and implements "static" methods provided by the TokenArray's type.
type TokenArrayType interface {
	NewBuilder() TokenArrayBuilder
	Concat(a, b TokenArray) TokenArray
}

// TokenArrayBuilder is a TokenWriter with a TokenArray() method returning the resulting array added on to it.
type TokenArrayBuilder interface {
	TokenWriter
	TokenArray() TokenArray
}

// TokenArrayReader is a struct implementing a TokenReader given a TokenArray.
type TokenArrayReader struct {
	tokenArray TokenArray
	idx        int
}

// NewTokenArrayReader create a new TokenArrayReader given a TokenArray.
func NewTokenArrayReader(t TokenArray) *TokenArrayReader {
	return &TokenArrayReader{
		tokenArray: t,
	}
}

// ReadToken Reads a single token from the token array. Returns io.EOF when all tokens have been read.
func (r *TokenArrayReader) ReadToken() (interface{}, error) {
	if r.idx < r.tokenArray.Len() {
		token := r.tokenArray.At(r.idx)
		r.idx++
		return token, nil
	}
	return nil, io.EOF
}
