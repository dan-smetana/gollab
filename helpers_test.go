package gollab_test

import (
	"github.com/danielslee/gollab/runetoken"
	"math/rand"
	"testing"

	"github.com/danielslee/gollab"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ가나다라마바사아자차카타파하")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randomPrimitive() gollab.PrimitiveOp {
	switch rand.Intn(3) {
	case 0:
		return gollab.Retain{Count: 1}
	case 1:
		return gollab.Insert{Tokens: runetoken.Array(randString(1))}
	case 2:
		return gollab.Delete{Count: 1}
	}
	return nil
}

func randomPrimitiveOps(inLength, outLength int) []gollab.PrimitiveOp {
	ops := gollab.CompositeOp([]gollab.PrimitiveOp{})
	for {
		if ops.InputLength() == inLength && ops.OutputLength() == outLength {
			return ops
		}

		if ops.InputLength() == inLength {
			ops = append(ops, gollab.Insert{Tokens: runetoken.Array(randString(1))})
		} else if ops.OutputLength() == outLength {
			ops = append(ops, gollab.Delete{Count: 1})
		} else {
			ops = append(ops, randomPrimitive())
		}
	}
}

func randomCompositeOp(inLength, outLength int) gollab.CompositeOp {
	return gollab.NewCompositeOp(randomPrimitiveOps(inLength, outLength)...)
}

func testEquality(t *testing.T, op1, op2 gollab.Op) {
	if op1.InputLength() != op2.InputLength() {
		t.Error("length mismatch")
		return
	}
	input := randString(op1.InputLength())
	op1Output, err := runetoken.ApplyToString(op1, input)
	if err != nil {
		t.Error(err)
		return
	}
	op2Output, err := runetoken.ApplyToString(op2, input)
	if err != nil {
		t.Error(err)
		return
	}

	if op1Output != op2Output {
		t.Errorf("op1(%s) != op2(%s)", op1Output, op2Output)
	}
}
