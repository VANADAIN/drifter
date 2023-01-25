package server

import "golang.org/x/net/websocket"

type ConnAction struct {
	action string
	conn   *websocket.Conn
}

func NewConnectionAction(action_type string, conn *websocket.Conn) *ConnAction {
	new_conn := &ConnAction{
		action: action_type,
		conn:   conn,
	}

	return new_conn
}
