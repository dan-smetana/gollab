package gollab_test

import (
	"fmt"
	"github.com/danielslee/gollab"
	"github.com/danielslee/gollab/runetoken"
)

func Example() {
	a := gollab.NewCompositeOp(
		gollab.Delete{Count: 1},
		gollab.Insert{Tokens: runetoken.Array("H")},
		gollab.Retain{Count: 4},
		gollab.Insert{Tokens: runetoken.Array(", World")})

	b := gollab.NewCompositeOp(gollab.Retain{Count: 5}, gollab.Insert{Tokens: runetoken.Array("!")})

	afterA, _ := runetoken.ApplyToString(a, "hello")
	afterB, _ := runetoken.ApplyToString(b, "hello")

	fmt.Println("after a:", afterA)
	fmt.Println("after b:", afterB)

	// Output:
	// after a: Hello, World
	// after b: hello!
}

func Example_transform() {
	a := gollab.NewCompositeOp(
		gollab.Delete{Count: 1},
		gollab.Insert{Tokens: runetoken.Array("H")},
		gollab.Retain{Count: 4},
		gollab.Insert{Tokens: runetoken.Array(", World")})

	b := gollab.NewCompositeOp(gollab.Retain{Count: 5}, gollab.Insert{Tokens: runetoken.Array("!")})

	aPrime, bPrime := a.Transform(b)
	fmt.Println("a':", aPrime)
	fmt.Println("b':", bPrime)

	// Output:
	//a': [Insert(H) Delete(1) Retain(4) Insert(, World) Retain(1)]
	//b': [Retain(12) Insert(!)]
}

func Example_compose() {
	composedOp := gollab.Compose(
		gollab.NewCompositeOp(gollab.Insert{Tokens: runetoken.Array("H")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 1}, gollab.Insert{Tokens: runetoken.Array("e")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 2}, gollab.Insert{Tokens: runetoken.Array("l")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 3}, gollab.Insert{Tokens: runetoken.Array("l")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 4}, gollab.Insert{Tokens: runetoken.Array("o")}))

	fmt.Println("composedOp:", composedOp)
	// Output: composedOp: [Insert(Hello)]
}

func ExampleCompose() {
	composedOp := gollab.Compose(
		gollab.NewCompositeOp(gollab.Insert{Tokens: runetoken.Array("H")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 1}, gollab.Insert{Tokens: runetoken.Array("e")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 2}, gollab.Insert{Tokens: runetoken.Array("l")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 3}, gollab.Insert{Tokens: runetoken.Array("l")}),
		gollab.NewCompositeOp(gollab.Retain{Count: 4}, gollab.Insert{Tokens: runetoken.Array("o")}))

	fmt.Println("composedOp:", composedOp)
	// Output: composedOp: [Insert(Hello)]
}

func ExampleCompositeOp_Transform() {
	a := gollab.NewCompositeOp(
		gollab.Delete{Count: 1},
		gollab.Insert{Tokens: runetoken.Array("H")},
		gollab.Retain{Count: 4},
		gollab.Insert{Tokens: runetoken.Array(", World")})

	b := gollab.NewCompositeOp(gollab.Retain{Count: 5}, gollab.Insert{Tokens: runetoken.Array("!")})

	aPrime, bPrime := a.Transform(b)

	fmt.Println("a':", aPrime)
	fmt.Println("b':", bPrime)

	const initialString = "hello"

	// try applying a first followed by b'
	afterA, _ := runetoken.ApplyToString(a, initialString)
	fmt.Println("after a:", afterA)
	afterAAndBPrime, _ := runetoken.ApplyToString(bPrime, afterA)
	fmt.Println("after a and b':", afterAAndBPrime)

	// try applying b first followed by a'
	afterB, _ := runetoken.ApplyToString(b, initialString)
	fmt.Println("after b:", afterB)
	afterBAndAPrime, _ := runetoken.ApplyToString(aPrime, afterB)
	fmt.Println("after b and a':", afterBAndAPrime)
	// Output:
	// a': [Insert(H) Delete(1) Retain(4) Insert(, World) Retain(1)]
	// b': [Retain(12) Insert(!)]
	// after a: Hello, World
	// after a and b': Hello, World!
	// after b: hello!
	// after b and a': Hello, World!
}
