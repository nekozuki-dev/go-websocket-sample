package wsservice

import (
	"fmt"
	"go-websocket-sample/app"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn interface {
	Run(readCh chan []byte, closeCh chan bool)
	Write([]byte)
	Close()
}

func NewConn(ws *websocket.Conn) app.Conn {
	return &conn{
		ws: ws,
	}
}

type conn struct {
	ws      *websocket.Conn
	wg      *sync.WaitGroup
	writeCh chan []byte
}

func (c *conn) Run(readCh chan []byte, closeCh chan bool) {
	c.wg = &sync.WaitGroup{}
	c.writeCh = make(chan []byte)

	// Wait Read Goroutine
	errCh := make(chan error)
	c.wg.Add(1)
	go c.waitRead(readCh, errCh)

	// Wait Write Groutine
	c.wg.Add(1)
	go c.waitWrite()

	// Process
	for {
		select {
		case <-errCh:
			close(c.writeCh)
			c.wg.Wait()

			close(closeCh)
			return
		}
	}
}

func (c *conn) Write(bytes []byte) {
	c.writeCh <- bytes
}

func (c *conn) Close() {
	c.ws.Close()
}

func (c *conn) waitWrite() {
	defer c.wg.Done()

	fmt.Println("Begin WaitWrite Goroutine.")
	for bytes := range c.writeCh {
		if err := c.ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
			fmt.Println("Error", err)
			break
		}
	}
	c.Close()
	fmt.Println("End WaitWrite Goroutine.")
}

func (c *conn) waitRead(readCh chan []byte, errCh chan error) {
	defer c.wg.Done()

	fmt.Println("Begin WaitRead Goroutine.")
	for {
		_, readBytes, err := c.ws.ReadMessage()
		if err != nil {
			errCh <- err
			break
		}
		readCh <- readBytes
	}
	c.Close()
	fmt.Println("End WaitRead Goroutine.")
}
