package gollab

import (
	"encoding/json"
	"errors"
)

type jsonOp struct {
	Type   string     `json:"type"`
	Tokens TokenArray `json:"tokens,omitempty"`
	Count  int        `json:"count,omitempty"`
}

// ErrUnknownOp represents an unknown operation error which can occur when parsing operations from JSON.
var ErrUnknownOp = errors.New("unknown operation")

func newJSONOp(op PrimitiveOp) (jsonOp, error) {
	switch op := op.(type) {
	case NoOp:
		return jsonOp{
			Type: "noop",
		}, nil
	case Retain:
		return jsonOp{
			Type:  "retain",
			Count: op.Count,
		}, nil
	case Delete:
		return jsonOp{
			Type:  "delete",
			Count: op.Count,
		}, nil
	case Insert:
		return jsonOp{
			Type:   "insert",
			Tokens: op.Tokens,
		}, nil
	}
	return jsonOp{}, ErrUnknownOp
}

func (j jsonOp) toOp() (PrimitiveOp, error) {
	switch j.Type {
	case "noop":
		return NoOp{}, nil
	case "retain":
		return Retain{Count: j.Count}, nil
	case "delete":
		return Delete{Count: j.Count}, nil
	case "insert":
		return Insert{Tokens: j.Tokens}, nil
	}
	return nil, ErrUnknownOp
}

// UnmarshalJSON decodes a CompositeOp from JSON.
//
// Check the type documentation for information on the format.
func (c *CompositeOp) UnmarshalJSON(data []byte) error {
	var ops []jsonOp
	if err := json.Unmarshal(data, &ops); err != nil {
		return err
	}

	newC := make(CompositeOp, len(ops))
	for i, op := range ops {
		unmarshalled, err := op.toOp()
		if err != nil {
			return err
		}
		newC[i] = unmarshalled
	}

	*c = newC
	return nil
}

// MarshalJSON encode a CompositeOp into JSON.
//
// Check the type documentation for information on the format.
func (c CompositeOp) MarshalJSON() ([]byte, error) {
	ops := make([]jsonOp, len(c))
	for i, op := range c {
		j, err := newJSONOp(op)
		if err != nil {
			return nil, err
		}
		ops[i] = j
	}

	return json.Marshal(ops)
}
