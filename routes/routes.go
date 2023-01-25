package routes

import (
	"github.com/VANADAIN/drifter/server"
	"github.com/VANADAIN/drifter/types"
)

// type Router struct {

// }

func Route(s *server.Server, msg *types.Message) {
	mtype := msg.Body.Type

	switch mtype {

	case "text":
		textHandler(msg)

	case "barter":
		barterHandler(s, msg)

	case "tunneljob":
		tjHandler(msg)
	}
}
