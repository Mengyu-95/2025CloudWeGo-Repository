package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/yourusername/kitex-gateway-demo/kitex_gen/api/userservice"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	svr := userservice.NewServer(
		new(UserServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "UserService",
		}),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
