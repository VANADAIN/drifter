package routes

import (
	"encoding/json"

	"github.com/VANADAIN/drifter/server"
	"github.com/VANADAIN/drifter/types"
)

// barter represents info swaps between nodes

func barterHandler(s *server.Server, msg *types.Message) {
	var barter_conns []string
	json.Unmarshal(msg.Body.Payload, &barter_conns)

	// compare connections known by server and bartered connections
	for _, barter_conn := range barter_conns {
		if !inlist(barter_conn, s.CH.KnownConns) {
			// concurrent safe ???
			s.CH.KnownConns = append(s.CH.KnownConns, barter_conn)
		}
	}
}

func inlist(conn string, conns []string) bool {
	for _, value := range conns {
		if value == conn {
			return true
		}
	}

	return false
}
