package server

import (
	"github.com/cloudwego/kitex/server"
	"github.com/yourusername/kitex-gateway-demo/pkg/protocol/detection"
	"github.com/yourusername/kitex-gateway-demo/pkg/protocol/httpcodec"
)

func WithMultiProtocol() server.Option {
	return server.WithTransHandlerFactory(detection.NewSvrTransHandlerFactory())
}

func WithHTTPCodec() server.Option {
	return server.WithCodec(httpcodec.NewHTTPCodec())
}
