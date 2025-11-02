package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	grpcGatewayIp        = "x-forwarded-for"
	grpcCallUserAgent    = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if agents := md.Get(grpcGatewayUserAgent); len(agents) > 0 {
			mtdt.UserAgent = agents[0]
		}

		if agents := md.Get(grpcCallUserAgent); len(agents) > 0 {
			mtdt.UserAgent = agents[0]
		}

		if ips := md.Get(grpcGatewayIp); len(ips) > 0 {
			mtdt.ClientIp = ips[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIp = p.Addr.String()
	}

	return mtdt
}
