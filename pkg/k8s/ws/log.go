package ws

import (
	"bufio"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"time"
)

type PodLogSession struct {
	wsConn       *websocket.Conn
	doneChan     chan struct{}
	ioReadCloser io.ReadCloser
	bufReader    *bufio.Reader
}

func NewPodLogSession(w http.ResponseWriter, r *http.Request) (*PodLogSession, error) {

	upgrader := &websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	session := &PodLogSession{
		wsConn:   conn,
		doneChan: make(chan struct{}),
	}

	return session, nil
}

func (p *PodLogSession) SetReaderCloser(readCloser io.ReadCloser) {
	p.ioReadCloser = readCloser
	p.bufReader = bufio.NewReader(readCloser)
}

func (p *PodLogSession) Read() {
	_, _, err := p.wsConn.ReadMessage()
	if err != nil {
		close(p.doneChan)
	}
}

func (p *PodLogSession) Write() error {
	for {
		bytes, err := p.bufReader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			return err
		}
		err = p.wsConn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			return err
		}
	}
}

func (p *PodLogSession) Close() error {
	return p.wsConn.Close()
}

func (p *PodLogSession) Done() chan struct{} {
	return p.doneChan
}
