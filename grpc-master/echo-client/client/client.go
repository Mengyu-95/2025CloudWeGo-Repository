package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc/echo"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func getContext(ctx context.Context) context.Context {
	md := getMetadataByMap(map[string]string{"time": time.Now().Format("2006-01-02T15-04-05Z07:00"), "header_data": "true"})
	//将数据写入ctx
	ctx = getOutgoingContext(ctx, md)
	//将数据附加到ctx
	ctx = appendToOutgoingContext(ctx, "k1", "value1", "k2", "value2")
	//md1, _ := metadata.FromOutgoingContext(ctx)
	//fmt.Println(md1)
	return ctx
}

func CallUnary(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = getContext(ctx)
	in := &echo.EchoRequest{
		Message: "client send message",
		Time:    timestamppb.New(time.Now()),
	}
	var header, trailer metadata.MD
	res, err := client.UnaryEcho(ctx, in, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client recv: %v\n", res.Message)
	fmt.Println("client recv UnaryEcho metadata header: ", header)
	fmt.Println("client recv UnaryEcho metadata trailer: ", trailer)
}
func CallServerStream(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = getContext(ctx)
	in := &echo.EchoRequest{
		Message: "client send message",
		Time:    timestamppb.New(time.Now()),
	}
	stream, err := client.ServerStreamingEcho(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	header, _ := stream.Header()
	fmt.Println("client recv ServerStreamingEcho metadata header: ", header)

	filename := "echo-client/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		file.Write(res.Bytes[:res.Length])
		//fmt.Printf("client recv : %v\n", res.Message)
	}
	stream.CloseSend()

	trailer := stream.Trailer()
	fmt.Println("client recv ServerStreamingEcho metadata trailer: ", trailer)
}
func CallClientStream(client echo.EchoClient) {
	filePath := "echo-client/files/client.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	ctx = getContext(ctx)
	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	header, _ := stream.Header()
	fmt.Println("client recv ClientStreamingEcho metadata header: ", header)

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stream.Send(&echo.EchoRequest{
			Message: "client sending file",
			Bytes:   buf[:n],
			Time:    timestamppb.New(time.Now()),
			Length:  int32(n),
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client recv : %v\n", res.Message)
	trailer := stream.Trailer()
	fmt.Println("client recv ClientStreamingEcho metadata trailer: ", trailer)
}
func CallBidirectional(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = getContext(ctx)
	stream, err := client.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	header, _ := stream.Header()
	fmt.Println("client recv BidirectionalStreamingEcho metadata header: ", header)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath := "echo-client/files/client.jpg"
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			stream.Send(&echo.EchoRequest{
				Message: "client sending file",
				Bytes:   buf[:n],
				Time:    timestamppb.New(time.Now()),
				Length:  int32(n),
			})
		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filename := "echo-client/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			file.Write(res.Bytes[:res.Length])
			//fmt.Printf("client recv : %v\n", res.Message)
		}
	}()
	wg.Wait()
	trailer := stream.Trailer()
	fmt.Println("client recv BidirectionalStreamingEcho metadata trailer: ", trailer)
}
