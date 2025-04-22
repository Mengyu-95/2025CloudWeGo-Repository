package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func GetKeepaliveOpt() (opt grpc.DialOption) {
	var kacp = keepalive.ClientParameters{
		// 如果没有活动流，则每10s发送一次ping
		Time: 10 * time.Second,
		// ping 超时时长
		Timeout: time.Second,
		//当没任何活动流的情况下，是否允许被ping
		PermitWithoutStream: true,
	}
	return grpc.WithKeepaliveParams(kacp)
}
