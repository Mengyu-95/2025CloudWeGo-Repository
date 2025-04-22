package detection

import (
	"context"
	"errors"
	"net"
	"regexp"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2"
)

var httpReg = regexp.MustCompile(`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)`)

type svrTransHandler struct {
	remote.ServerTransHandler
}

func (t *svrTransHandler) ProtocolMatch(ctx context.Context, conn net.Conn) error {
	buf := make([]byte, 4)
	n, _ := conn.Read(buf)
	if n < 4 {
		return errors.New("insufficient data")
	}

	if httpReg.Match(buf) {
		return nil
	}

	return errors.New("protocol not supported")
}

func NewSvrTransHandlerFactory() remote.ServerTransHandlerFactory {
	return &svrTransHandlerFactory{}
}

type svrTransHandlerFactory struct{}

func (f *svrTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	// 使用新版Kitex的API创建gRPC传输处理器
	th, err := nphttp2.NewSvrTransHandlerFactory().NewTransHandler(opt)
	if err != nil {
		return nil, err
	}

	return &svrTransHandler{
		ServerTransHandler: th,
	}, nil
}
