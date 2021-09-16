package gollab_test

import (
	"fmt"
	"github.com/danielslee/gollab/runetoken"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/danielslee/gollab/client"
	"github.com/danielslee/gollab/server"
)

func randomWait() {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func slowedDownChannel(c chan<- server.ClientMessage) chan<- server.ClientMessage {
	newChan := make(chan server.ClientMessage)

	go func() {
		defer close(c)
		for msg := range newChan {
			randomWait()
			c <- msg
		}
	}()

	return newChan
}

type otClient struct {
	mux   sync.Mutex
	state client.State
	id    string
	numID int

	document string

	sendChan    chan<- server.ClientMessage
	receiveChan <-chan interface{}
}

func (c *otClient) typeRandomStuff() {
	for i := 0; i < 10; i++ {
		randomWait()

		c.mux.Lock()
		docLength := len([]rune(c.document))
		op := randomCompositeOp(docLength, docLength+3)
		newState, sendAwaiting := c.state.ApplyClientOp(op)

		if sendAwaiting {
			c.sendChan <- server.ClientMessage{
				ClientID: c.numID,
				Message: server.OpMessage{
					AuthorID: c.id,
					Op:       newState.Awaiting,
					Revision: newState.Revision,
				},
			}
		}

		c.print("applying client op:", op)
		c.printStateChange(c.state, newState)
		c.state = newState
		newDoc, err := runetoken.ApplyToString(op, c.document)
		if err != nil {
			panic(err)
		}

		c.document = newDoc

		c.mux.Unlock()
	}
}

func (c *otClient) run() {
	for srvMsg := range c.receiveChan {
		randomWait()

		msg, ok := srvMsg.(server.OpMessage)
		if !ok {
			continue
		}

		if msg.AuthorID == c.id {
			c.mux.Lock()

			newState, sendAwaiting := c.state.ApplyServerAck()
			if sendAwaiting {
				c.sendChan <- server.ClientMessage{
					ClientID: c.numID,
					Message: server.OpMessage{
						AuthorID: c.id,
						Op:       newState.Awaiting,
						Revision: newState.Revision,
					},
				}
			}

			c.print("applying ack")
			c.printStateChange(c.state, newState)
			c.state = newState

			c.mux.Unlock()
		} else {
			c.mux.Lock()
			newState, docOp := c.state.ApplyServerOp(msg.Op)
			c.print("applying server op:", msg.Op)
			c.printStateChange(c.state, newState)
			c.state = newState

			newDoc, err := runetoken.ApplyToString(docOp, c.document)
			if err != nil {
				panic(err)
			}

			c.document = newDoc

			c.mux.Unlock()
		}
	}
}

func TestClientServer(t *testing.T) {
	memoryStore := server.NewMemoryStateStore(runetoken.Array{})
	d := server.NewDocumentServer(memoryStore)
	serverChan := slowedDownChannel(d.ReceiveChan())

	var clients []*otClient
	for i := 0; i < 15; i++ {
		id, clientChan := d.NewClient()
		clients = append(clients, &otClient{
			id:          strconv.Itoa(id),
			numID:       id,
			sendChan:    serverChan,
			receiveChan: clientChan,
		})
	}

	var wg sync.WaitGroup
	var typeStuffWG sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.Run()
	}()

	for _, c := range clients {
		wg.Add(1)
		go func(client *otClient) {
			defer wg.Done()
			client.run()
		}(c)

		typeStuffWG.Add(1)
		go func(client *otClient) {
			defer typeStuffWG.Done()
			client.typeRandomStuff()
		}(c)
	}

	typeStuffWG.Wait()

	time.Sleep(10 * time.Second)
	log.Println("closing server chan")
	close(serverChan)
	wg.Wait()

	doc, _, _ := memoryStore.Current()
	for _, c := range clients {
		c.print(c.document)

		runeDoc := doc.(runetoken.Array)
		if runeDoc.String() != c.document {
			t.Errorf("client #%s: document is not the same", c.id)
		}
	}
}

func (c *otClient) print(args ...interface{}) {
	log.Println(append([]interface{}{fmt.Sprintf("\033[32mCLIENT(%s):\033[0m", c.id)}, args...)...)
}

func (c *otClient) printStateChange(oldState, newState client.State) {
	c.print("before:", oldState)
	c.print("after:", newState)
}
