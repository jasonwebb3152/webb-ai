package gapi

import (
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
)

func newTestServer(store db.Store) *Server {
	server := NewServer(store)
	return server
}
