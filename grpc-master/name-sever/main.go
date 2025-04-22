package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc/name"
	"grpc/name-sever/server"
	"log"
	"net"
	"time"
)

var (
	port = flag.Int("port", 60051, "")
)

func main() {
	//testData()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	name.RegisterNameServer(s, &server.NameServer{})
	log.Printf("server listening at : %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}

}

func testData() {
	server.Register("echo", "localhost:50051")
	server.Register("echo", "localhost:50052")
	time.Sleep(time.Second * 2)
	server.Register("echo", "localhost:50053")
	server.Register("echo", "localhost:50054")
	time.Sleep(time.Second * 2)
	server.Register("echo", "localhost:50055")
	server.Register("echo", "localhost:50056")
	server.Register("echo", "localhost:50051")
	server.Register("echo", "localhost:50052")
	time.Sleep(time.Second * 2)
	server.Register("echo", "localhost:50053")
	server.Register("echo", "localhost:50054")
	time.Sleep(time.Second * 2)
	server.Register("echo", "localhost:50055")
	server.Register("echo", "localhost:50056")
	allData := server.GetAllData()
	fmt.Println(allData)
	server.Delete("echo", "localhost:50056")
	fmt.Println(server.GetByServiceName("echo"))
	allData = server.GetAllData()
	fmt.Println(allData)

}
