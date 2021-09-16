package gollab_test

import (
	"fmt"
	"github.com/danielslee/gollab/runetoken"
	"math/rand"
	"testing"

	"github.com/danielslee/gollab"
)

func TestCompositeOpApplyString(t *testing.T) {
	op := gollab.NewCompositeOp(
		gollab.Delete{Count: 1},
		gollab.Insert{Tokens: runetoken.Array("H")},
		gollab.Retain{Count: 4},
		gollab.Insert{Tokens: runetoken.Array(",")},
		gollab.Retain{Count: 1},
		gollab.Delete{Count: 1},
		gollab.Insert{Tokens: runetoken.Array("W")},
		gollab.Retain{Count: 4},
		gollab.Insert{Tokens: runetoken.Array("!")},
	)

	applied, err := runetoken.ApplyToString(op, "hello world")

	if err != nil {
		t.Error(err)
		return
	}

	if applied != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got '%s'", applied)
	}
}

func TestNewCompositeOp(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("new-op-%d", i), func(t *testing.T) {
			t.Parallel()
			inLength := rand.Intn(10) + 5
			outLength := inLength + rand.Intn(10) - 5
			primitives := randomPrimitiveOps(inLength, outLength)

			primitivesCopy := make([]gollab.PrimitiveOp, len(primitives))
			copy(primitivesCopy, primitives)

			op := gollab.NewCompositeOp(primitivesCopy...)

			for i := 0; i < len(op)-1; i++ {
				a, b := op[i], op[i+1]
				if joinable, ok := a.(gollab.Joinable); ok {
					if joinable.Join(b) != nil {
						t.Errorf("joinable CompositeOp found: %v, ops: %v", joinable, op)
						return
					}
				}

				if swappable, ok := a.(gollab.Swappable); ok {
					if swappable.Swap(b) {
						t.Errorf("swappable CompositeOp found: %v, ops: %v", swappable, op)
						return
					}
				}
			}

			t.Run("original-simplified-equality", func(t *testing.T) {
				testEquality(t, gollab.CompositeOp(primitives), op)
			})
		})
	}
}
