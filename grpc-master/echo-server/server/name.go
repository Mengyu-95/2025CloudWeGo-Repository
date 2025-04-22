package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/name"
	"log"
	"time"
)

type NameServer struct {
	conn *grpc.ClientConn
}

func NewNameServer(addr string) *NameServer {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	return &NameServer{
		conn: conn,
	}
}

func (ns *NameServer) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	ns.conn.Close()
}

func (ns *NameServer) RegisterName(serviceName, address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
		Address:     []string{address},
	}
	_, err := client.Register(context.Background(), in)
	if err != nil {
		log.Println(err)
	}
}

func (ns *NameServer) Delete(serviceName, address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
		Address:     []string{address},
	}
	_, err := client.Delete(context.Background(), in)
	if err != nil {
		log.Println(err)
	}
}

func (ns *NameServer) Keepalive(serviceName, address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	client := name.NewNameClient(ns.conn)
	in := &name.NameRequest{
		ServiceName: serviceName,
		Address:     []string{address},
	}
	stream, err := client.Keepalive(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	for {
		stream.Send(in)
		time.Sleep(time.Second)
	}
}
