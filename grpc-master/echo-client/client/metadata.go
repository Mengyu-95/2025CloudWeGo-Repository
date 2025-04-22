package client

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func getMetadataByMap(mp map[string]string) metadata.MD {
	//通过map 初始化metadata
	md := metadata.New(mp)
	return md
}

func getMetadataByKV(kv ...string) metadata.MD {
	//通过键值对的方式初始化metadata
	md := metadata.Pairs(kv...)
	return md
}

func getOutgoingContext(ctx context.Context, md metadata.MD) context.Context {
	// OutgoingContext 用于请求发送方，包装数据传递出去
	// IncomingContext 用于请求接收方，用于获取发送方传递的数据
	// Context 通过序列化成 http2 header 的方式传输
	// new 方法会覆盖ctx 原有元数据
	return metadata.NewOutgoingContext(ctx, md)
}

// 将数据附加到OutgoingContext
func appendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, kv...)
}
