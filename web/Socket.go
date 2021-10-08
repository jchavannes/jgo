package web

import (
	"github.com/gorilla/websocket"
	"github.com/jchavannes/jgo/jerr"
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

func (s *Socket) Ping() error {
	err := s.ws.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		return jerr.Get("error writing ping message", err)
	}
	return nil
}

func (s *Socket) Close() error {
	if err := s.ws.Close(); err != nil {
		return jerr.Get("error closing socket", err)
	}
	return nil
}

// OnClose event does not get triggered unless messages are read.
func (s *Socket) OnClose(closeHandler func()) {
	s.ws.SetCloseHandler(func(code int, text string) error {
		closeHandler()
		return nil
	})
}

// ReadAllUntilClose throws away all incoming messages. If you want to read messages, don't use this function.
func (s *Socket) ReadAllUntilClose(closeHandler chan error) {
	go func() {
		for {
			if _, err := s.ReadMessage(); err != nil {
				closeHandler <- jerr.Get("verification socket closed", err)
				return
			}
		}
	}()
}
