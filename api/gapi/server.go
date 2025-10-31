package gapi

import (
	"github.com/jasonwebb3152/webb-ai/api/pb"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
)

type Server struct {
	store db.Store
	pb.UnimplementedWebbAiServer
}

func NewServer(store db.Store) *Server {
	return &Server{
		store: store,
	}
}
