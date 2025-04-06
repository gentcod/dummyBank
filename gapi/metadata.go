package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgent    = "user-agent"
	gatewayUserAgent = "grpcgateway-user-agent"
	gatewayClientIP  = "x-forwarded-for"
)

type Metadata struct {
	ClientIP  string
	UserAgent string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcUserAgent); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(gatewayUserAgent); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(gatewayClientIP); len(clientIps) > 0 {
			mtdt.ClientIP = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
