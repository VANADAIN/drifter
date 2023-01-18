package routes

import (
	"github.com/VANADAIN/drifter/types"
)

func Route(msg *types.Message) {
	mtype := msg.Body.Type

	switch mtype {

	case "text":
		textHandler(msg)

	case "barter":
		barterHandler(msg)

	case "tunneljob":
		tjHandler(msg)

	}
}
