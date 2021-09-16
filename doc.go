/*
Package gollab implements four basic operations: Insert, Delete, Retain, NoOp. These operations form a Composite
Operation which can then be applied to a string.

For example, suppose we have a composite operation `a`:

	Delete(1), Insert("H"), Retain(4), Insert(", World")

which when applied to the string `hello` produces the following output `Hello, World`. This operation can only be
applied to a string of length 5. (the Delete operation consumes one character, and the Retain operation consumes 4, with
the remaining operations consuming zero characters, giving us a total of 5 characters)

Given an operation `b`, defined as:

	Retain(5), Insert("!")

`b` can be applied to the same initial string `hello` and we end up with `hello!`

# Transformation

Now suppose these two operations happened concurrently. We need to ensure we end up with the same result, regardless of
the order we apply them in. This is where OT comes in.

Gollab implements two operations that can be applied to composite operations such as `a` and `b` above: Transform and
Compose. Operation Transformation, as suggested by its name, refers to the former.

OT is defined as `OT(a, b) -> a', b'`, where applying `a` first and `b'` second yields the same result as applying `b`
first and `a'` second.

Let's apply OT to `a` and `b` above. We end up with two new operations: `a'` and `b'`

	a' := Insert("H"), Delete(1), Retain(4), Insert(", World"), Retain(1)
	b' := Retain(12), Insert("!")


We can verify that the above transformation is correct by applying `a` to `hello`, resulting in `Hello, World`, which
gives us `Hello, World!` after applying `b'` to it. We arrive to the same result by applying `b` to hello
getting `hello!` followed by applying `a'` making it into `Hello, World!`.

# Composition

When building a collaborative editor we may want to combine two or more composite operations which were
applied in order into one. For example, we may want to merge `Insert("H")`, `Retain(1), Insert("e")`,
`Retain(2), Insert("l")`, `Retain(3), Insert("l")`, `Retain(4), Insert("o")` into `Insert("Hello")`.

# Tokens

Instead of being hard-coded to only operate on plain text, Gollab operates on a string of abstract `tokens`.
These tokens can be anything. Gollab includes an implementation of this for a string of runes (unicode codepoints) in
the runetoken package.

Check out the documentation for TokenReader, TokenWriter and TokenArray interfaces for more details on how to provide
your own implementation for use cases requiring more than plain unicode text.

# Client & Server

The above two operations are only a half of what's needed to build a collaborative editor. For that we need a client
and a server.

Gollab implements a client and a server in the `client` and `server` packages.

# Examples

See below for examples of how to implement the pseudocode above using gollab.
*/
package gollab
