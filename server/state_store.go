package server

import (
	"sync"

	"github.com/danielslee/gollab"
)

// StateStore stores the current document-state server side. See MemoryStateStore for a basic implementation. More
// complex implementations can be implemented to use a database like Redis.
type StateStore interface {
	Current() (document gollab.TokenArray, revision int, err error)
	ApplyClient(opMsg OpMessage) error
	OperationStream() <-chan OpMessage
}

// MemoryStateStore implements a basic StateStore.
type MemoryStateStore struct {
	mux sync.RWMutex

	document gollab.TokenArray
	ops      []gollab.CompositeOp
	opStream chan OpMessage
}

// NewMemoryStateStore Creates a new NewMemoryStateStore.
func NewMemoryStateStore(document gollab.TokenArray) *MemoryStateStore {
	return &MemoryStateStore{
		opStream: make(chan OpMessage, 128),
		document: document,
	}
}

// Current returns the current state consisting of the document and its revision number.
func (m *MemoryStateStore) Current() (document gollab.TokenArray, revision int, err error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	document = m.document
	revision = len(m.ops)
	return
}

// ApplyClient applies a client-side operation.
func (m *MemoryStateStore) ApplyClient(opMsg OpMessage) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if opMsg.Revision < 0 || (opMsg.Revision > len(m.ops)) {
		return ErrUnknownRevision
	}

	res, err := ApplyClientOp(ApplyClientOpInput{
		CurrentDocument: m.document,
		CurrentRevision: len(m.ops),
		Op:              opMsg.Op,
		TransformOps:    m.ops[opMsg.Revision:],
	})

	if err != nil {
		return err
	}

	m.document = res.Document
	m.ops = append(m.ops, res.Op)

	m.opStream <- OpMessage{
		AuthorID: opMsg.AuthorID,
		Op:       res.Op,
		Revision: res.Revision,
	}

	return nil
}

// OperationStream is a channel returning operations to be broadcast to all clients.
func (m *MemoryStateStore) OperationStream() <-chan OpMessage {
	return m.opStream
}
