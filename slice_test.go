package gollab_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/danielslee/gollab"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func testSlice(t *testing.T, op1, op2 gollab.CompositeOp,
	aLengthFunc, bLengthFunc func(op gollab.Op) int) {
	op1Sliced, op2Sliced := gollab.Slice(op1, op2, aLengthFunc, bLengthFunc)

	t.Run("op1-equality", func(t *testing.T) {
		testEquality(t, op1, gollab.CompositeOp(op1Sliced))
	})

	t.Run("op2-equality", func(t *testing.T) {
		testEquality(t, op2, gollab.CompositeOp(op2Sliced))
	})

	if len(op1Sliced) != len(op2Sliced) {
		t.Errorf("len(op1Sliced)(%d) != len(op2Sliced)(%d)",
			len(op1Sliced), len(op2Sliced))
	}

	for i := range op1Sliced {
		a, b := op1Sliced[i], op2Sliced[i]

		if aLengthFunc(a) != bLengthFunc(b) {
			t.Errorf("length mismatch: %v(%d), %v(%d)", a, b,
				aLengthFunc(a), bLengthFunc(b))
		}
	}
}

func TestSliceInputInput(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("rand-%d", i), func(t *testing.T) {
			t.Parallel()
			l := rand.Intn(20) + 6
			op1 := randomCompositeOp(l, l+rand.Intn(10)-5)
			op2 := randomCompositeOp(l, l+rand.Intn(10)-5)
			testSlice(t, op1, op2, gollab.InputLengthFunc, gollab.InputLengthFunc)
		})
	}
}

func TestSliceInputOutput(t *testing.T) {
	for i := 0; i < 10000; i++ {
		t.Run(fmt.Sprintf("rand-%d", i), func(t *testing.T) {
			t.Parallel()
			inputLength := rand.Intn(20) + 5
			outLength := inputLength + rand.Intn(10) - 5
			op1 := randomCompositeOp(inputLength, outLength)
			op2 := randomCompositeOp(outLength, outLength+rand.Intn(5))

			testSlice(t, op1, op2, gollab.OutputLengthFunc, gollab.InputLengthFunc)
		})
	}
}
