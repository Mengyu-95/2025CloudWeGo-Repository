package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

func main() {
	// 1. 初始化 Thrift 泛化调用客户端（指向 Kitex 服务端口）
	provider, _ := generic.NewThriftFileProvider("../idl/user.thrift") // 根据实际路径调整
	g, _ := generic.HTTPThriftGeneric(provider)
	cli, _ := genericclient.NewClient(
		"UserService",
		g,
		client.WithHostPorts("101.126.12.5:8888"),
		// client.WithHostPorts("127.0.0.1:8888"), // Kitex 服务地址
	)

	// 2. 启动 Hertz HTTP 服务
	h := server.Default(server.WithHostPorts(":8080")) // 网关监听 8080 端口
	h.POST("/user/stream", func(c context.Context, ctx *app.RequestContext) {
		// 将 HTTP 请求转发给 Kitex 服务
		resp, err := cli.GenericCall(c, "StreamUsers", ctx.Request.Body())
		if err != nil {
			ctx.JSON(500, map[string]string{"error": err.Error()})
			return
		}
		ctx.Data(200, "application/json", resp.([]byte))
	})
	h.Spin()
}
