package gapi

import (
	"github.com/jasonwebb3152/webb-ai/api/pb"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util/token"
)

type Server struct {
	store db.Store
	maker token.Maker
	pb.UnimplementedWebbAiServer
}

func NewServer(store db.Store, symmetricKey string) (*Server, error) {
	pasetoMaker, err := token.NewPasetoMaker(symmetricKey)
	if err != nil {
		return nil, err
	}

	return &Server{
		store: store,
		maker: pasetoMaker,
	}, nil
}
