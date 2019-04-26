package wiring

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/handler"
)

// Params represent router params.
type Params struct {
	Logger    *zap.Logger
	Formatter handler.Formatter
	Sender    handler.Sender
}

// NewRouter configure all router.
func NewRouter(params *Params) *mux.Router {
	rtr := mux.NewRouter().StrictSlash(true)
	rtr.Handle("/api/v1/sms/send", handler.Send(params.Logger, params.Sender, params.Formatter)).Methods("POST")
	return rtr
}
