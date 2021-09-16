package server

import (
	"log"
	"sync"

	"github.com/danielslee/gollab"
)

// InitMessage is the initial message sent by the server to a new client.
type InitMessage struct {
	Document gollab.TokenArray `json:"document"`
	Revision int               `json:"revision"`
}

// OpMessage is a message containing an operation along with its author's id and revision.
type OpMessage struct {
	AuthorID string             `json:"authorID"`
	Op       gollab.CompositeOp `json:"op"`
	Revision int                `json:"revision"`
}

// ClientMessage contains an OpMessage along with its sender's client id.
type ClientMessage struct {
	ClientID int
	Message  OpMessage
}

// ErrorMessage is a message signifying an error has occurred.
type ErrorMessage struct {
	Error string `json:"error"`
}

// DocumentServer implements a server serving a single document.
type DocumentServer struct {
	state StateStore

	receiveChan chan ClientMessage

	sendChannelsMux sync.RWMutex
	sendChannels    map[int]chan<- interface{}
	channelCounter  int
}

// NewDocumentServer creates a new document server given a StateStore.
func NewDocumentServer(stateStore StateStore) *DocumentServer {
	return &DocumentServer{
		state:        stateStore,
		receiveChan:  make(chan ClientMessage, 128),
		sendChannels: make(map[int]chan<- interface{}),
	}
}

// Run start serving clients.
func (d *DocumentServer) Run() {
	defer func() {
		d.sendChannelsMux.Lock()
		defer d.sendChannelsMux.Unlock()
		for _, c := range d.sendChannels {
			close(c)
		}
		d.sendChannels = make(map[int]chan<- interface{})
	}()

	for {
		select {
		case clientMsg, more := <-d.receiveChan:
			if !more {
				return
			}

			msg := clientMsg.Message

			err := d.state.ApplyClient(msg)
			if err != nil {
				log.Println("err applying operation:", err)
				d.sendError(clientMsg.ClientID, "invalid operation")
			}
		case op := <-d.state.OperationStream():
			d.send(op)
		}
	}
}

func (d *DocumentServer) send(msg OpMessage) {
	d.sendChannelsMux.RLock()
	defer d.sendChannelsMux.RUnlock()

	for _, c := range d.sendChannels {
		c <- msg
	}
}

func (d *DocumentServer) sendError(clientID int, err string) {
	d.sendChannelsMux.RLock()
	defer d.sendChannelsMux.RUnlock()
	if clientChan, ok := d.sendChannels[clientID]; ok {
		clientChan <- ErrorMessage{err}
		delete(d.sendChannels, clientID)
		close(clientChan)
	}
}

// NewClient creates and attaches a new client. It returns the client's id number and a channel on which the client
// can receive messages from the server.
func (d *DocumentServer) NewClient() (clientID int, sendToClientChan <-chan interface{}) {
	d.sendChannelsMux.Lock()
	defer d.sendChannelsMux.Unlock()

	clientID = d.channelCounter
	d.channelCounter++

	c := make(chan interface{}, 128)

	doc, rev, err := d.state.Current()
	if err != nil {
		panic(err)
	}

	c <- InitMessage{
		Document: doc,
		Revision: rev,
	}

	d.sendChannels[clientID] = c
	return clientID, c
}

// RemoveClient detaches a client.
func (d *DocumentServer) RemoveClient(id int) {
	d.sendChannelsMux.Lock()
	defer d.sendChannelsMux.Unlock()

	delete(d.sendChannels, id)
}

// ReceiveChan returns a channel on which the DocumentServer receiver messages from clients.
func (d *DocumentServer) ReceiveChan() chan<- ClientMessage {
	return d.receiveChan
}
