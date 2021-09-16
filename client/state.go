/*
Package client implements an immutable State struct representing the client's state.
*/
package client

import (
	"fmt"

	"github.com/danielslee/gollab"
)

// State represents immutable client state. Calling methods on it results in a new State.
// It consists of a revision and two buffers:
//  1. Awaiting - an operation that has been sent to the server but hasn't been returned yet (server ack hasn't been
//                received yet). When an ack is received Awaiting is flushed and any operation in Buffer is moved over
//                to awaiting and sent to the server.
//
//  2. Buffer  -  if there is anything in Awaiting, all new operations made by the client are placed in/appended to
//                (using gollab.CompositeOp.Compose) the Buffer.
type State struct {
	Revision int
	Buffer   gollab.CompositeOp
	Awaiting gollab.CompositeOp
}

// ApplyServerOp applies an operation received from server, returning new state
// and a transformed operation to be applied to the client's document.
func (s State) ApplyServerOp(op gollab.CompositeOp) (newState State, documentOp gollab.CompositeOp) {
	newState.Revision = s.Revision + 1

	if s.Awaiting == nil && s.Buffer == nil {
		documentOp = op
	} else if s.Awaiting != nil && s.Buffer == nil {
		awaitingPrime, opPrime := s.Awaiting.Transform(op)
		newState.Awaiting = awaitingPrime
		documentOp = opPrime
	} else {
		awaitingPrime, opPrime := s.Awaiting.Transform(op)
		newState.Awaiting = awaitingPrime

		bufferPrime, opDoublePrime := s.Buffer.Transform(opPrime)
		newState.Buffer = bufferPrime
		documentOp = opDoublePrime
	}

	return
}

// ApplyServerAck registers server acknowledgment, returning new state and
// a boolean signalling whether to send whatever is in the awaiting buffer.
func (s State) ApplyServerAck() (newState State, sendAwaiting bool) {
	newState.Revision = s.Revision + 1

	if s.Awaiting == nil && s.Buffer == nil {
		panic("received ack while not awaiting anything")
	} else if s.Awaiting != nil && s.Buffer == nil {
		newState.Awaiting = nil
	} else {
		newState.Awaiting = s.Buffer
		newState.Buffer = nil
		sendAwaiting = true
	}
	return
}

// ApplyClientOp registers a new client side operation, returning new state and
// a boolean signalling whether to send whatever is in the awaiting buffer.
func (s State) ApplyClientOp(op gollab.CompositeOp) (newState State,
	sendAwaiting bool) {
	newState.Revision = s.Revision

	if s.Awaiting == nil && s.Buffer == nil {
		newState.Awaiting = op
		sendAwaiting = true
	} else if s.Awaiting != nil && s.Buffer == nil {
		newState.Awaiting = s.Awaiting
		newState.Buffer = op
	} else {
		newState.Awaiting = s.Awaiting
		newState.Buffer = s.Buffer.Compose(op)
	}
	return
}

// String returns a string representation useful for debugging.
func (s State) String() string {
	return fmt.Sprintf(
		"client.State(revision: %d, Awaiting %v, Buffer: %v))",
		s.Revision, s.Awaiting, s.Buffer)
}
