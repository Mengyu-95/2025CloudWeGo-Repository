package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc/echo"
	"grpc/echo-client/client"
	client_pool "grpc/echo-client/client-pool"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50051", "")
)

func getOptions() (opts []grpc.DialOption) {
	opts = make([]grpc.DialOption, 0)
	//opts = append(opts, client.GetTlsOpt())
	opts = append(opts, client.GetMTlsOpt())
	opts = append(opts, grpc.WithUnaryInterceptor(client.UnaryInterceptor))
	opts = append(opts, grpc.WithStreamInterceptor(client.StreamInterceptor))
	opts = append(opts, client.GetAuth(client.FetchToken()))
	opts = append(opts, client.GetKeepaliveOpt())
	opts = append(opts, client.GetNameResolver(client.NewNameServer("localhost:60051")))
	return opts
}

func main() {
	flag.Parse()
	//根据地址访问
	//conn, err := grpc.Dial(*addr, getOptions()...)
	//根据 协议 + 服务名 通过名称服务解析，访问服务器
	//conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", client.MyScheme, client.MyServiceName), getOptions()...)

	pool, err := client_pool.GetPool(fmt.Sprintf("%s:///%s", client.MyScheme, client.MyServiceName), getOptions()...)
	if err != nil {
		log.Fatal(err)
	}
	conn := pool.Get()
	defer pool.Put(conn)

	//defer conn.Close()
	c := echo.NewEchoClient(conn)
	client.CallUnary(c)
	time.Sleep(5 * time.Second)
	client.CallServerStream(c)
	time.Sleep(5 * time.Second)
	client.CallClientStream(c)
	time.Sleep(5 * time.Second)
	client.CallBidirectional(c)
}
