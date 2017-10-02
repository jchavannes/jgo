package web

import (
	"github.com/gorilla/websocket"
)

type Socket struct {
	ws *websocket.Conn
}

func (s *Socket) ReadMessage() ([]byte, error) {
	_, msgByte, err := s.ws.ReadMessage()
	return msgByte, err
}

func (s *Socket) WriteJSON(v interface{}) error {
	return s.ws.WriteJSON(v)
}
