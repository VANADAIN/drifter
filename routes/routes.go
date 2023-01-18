package routes

import (
	"github.com/VANADAIN/drifter/types"
)

func Route(msg *types.Message) {
	mtype := msg.Body.Type

	if mtype == "text" {
		textHandler(msg)
	}
	if mtype == "barter" {
		barterHandler(msg)
	}
	if mtype == "tunneljob" {
		tjHandler(msg)
	}
}
