package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"time"
)

type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

func NewTerminalSession(w http.ResponseWriter, r *http.Request) (*TerminalSession, error) {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		// 检测请求来源
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		return copy(p, "\u0004"), err
	}
	var msg TerminalMessage
	if err = json.Unmarshal(message, &msg); err != nil {
		return copy(p, "\u0004"), err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		return copy(p, "\u0004"), fmt.Errorf("unknown message type")
	}
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		return 0, err
	}
	if err = t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t *TerminalSession) Done() {
	close(t.doneChan)
}

func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}
