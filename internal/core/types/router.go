package types

import (
	"github.com/gorilla/mux"
)

type Router interface {
	RegisterRoutes(router *mux.Router)
}
