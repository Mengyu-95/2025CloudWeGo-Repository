package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func GetKeepaliveOpt() (opts []grpc.ServerOption) {
	//服务端轻质保活策略，客户端违反该策略将被关闭
	var kaep = keepalive.EnforcementPolicy{
		//客户端ping服务器，最小时间间隔，小于该时间间隔强制关闭连接
		MinTime: 5 * time.Second,
		//当前没有任何活动流的情况下，是否允许被ping
		PermitWithoutStream: true,
	}

	var kasp = keepalive.ServerParameters{
		//客户端空闲15秒发送goaway 指令(尝试断开连接)
		MaxConnectionIdle: 15 * time.Second,
		//最大连接时长30s，超时发送goaway
		MaxConnectionAge: 30 * time.Second,
		//强制关闭前等待时长
		MaxConnectionAgeGrace: 5 * time.Second,
		//客户端空闲5秒，发送ping保活
		Time: 5 * time.Second,
		// ping 超时时间
		Timeout: 1 * time.Second,
	}
	return []grpc.ServerOption{grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp)}
}
