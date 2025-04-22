package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"grpc/echo"
	"grpc/echo-server/server"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	port = flag.Int("port", 50051, "")
)

func getOptions() (opts []grpc.ServerOption) {
	opts = make([]grpc.ServerOption, 0)
	//opts = append(opts, server.GetTlsOpt())
	opts = append(opts, server.GetMTlsOpt())
	opts = append(opts, grpc.UnaryInterceptor(server.UnaryInterceptor))
	opts = append(opts, grpc.StreamInterceptor(server.StreamInterceptor))
	opts = append(opts, server.GetKeepaliveOpt()...)
	return opts
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(getOptions()...)
	echo.RegisterEchoServer(s, &server.EchoServer{})

	//以多路复用的方式注册健康检查服务
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	log.Printf("server listening at : %v\n", lis.Addr())
	//启动echo server
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	//向 nameserver 注册服务信息，并保活
	nameServer := server.NewNameServer("localhost:60051")
	serviceName := "myecho"
	addr := fmt.Sprintf("localhost:%d", *port)
	go func() {
		nameServer.RegisterName(serviceName, addr)
		nameServer.Keepalive(serviceName, addr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()
	//停止服务，删除注册信息
	nameServer.Delete(serviceName, addr)
	nameServer.Close()
}
