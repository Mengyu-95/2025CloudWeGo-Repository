package server

import (
	"context"
	"grpc/name"
	"io"
	"log"
)

type NameServer struct {
	name.UnimplementedNameServer
}

func (NameServer) Register(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	for _, address := range in.Address {
		Register(in.ServiceName, address)
	}
	log.Println(GetByServiceName(in.ServiceName))
	return &name.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}
func (NameServer) Delete(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	for _, address := range in.Address {
		Delete(in.ServiceName, address)
	}
	log.Println(GetByServiceName(in.ServiceName))
	return &name.NameResponse{
		ServiceName: in.ServiceName,
	}, nil
}
func (NameServer) Keepalive(stream name.Name_KeepaliveServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println(req)
		for _, address := range req.Address {
			Keepalive(req.ServiceName, address)
		}
	}
	return stream.SendAndClose(&name.NameResponse{})
}
func (NameServer) GetAddress(ctx context.Context, in *name.NameRequest) (*name.NameResponse, error) {
	address := GetByServiceName(in.ServiceName)
	log.Println(address)
	return &name.NameResponse{
		ServiceName: in.ServiceName,
		Address:     address,
	}, nil
}
