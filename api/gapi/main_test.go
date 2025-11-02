package gapi

import (
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util"
)

func newTestServer(store db.Store) *Server {
	server, err := NewServer(store, util.RandomString(32))
	if err != nil {
		panic("failed to create the test server")
	}
	return server
}
