package server

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

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func getMetadata() (header metadata.MD, trailer metadata.MD) {
	header = getMetadataByMap(map[string]string{"server_time": time.Now().Format("2006-01-02T15-04-05Z07:00"), "server_header_data": "true"})
	trailer = getMetadataByKV("server_trailer_data", "true")
	return
}

func (EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	// 响应请求，发送元数据
	header, trailer := getMetadata()
	defer grpc.SetTrailer(ctx, trailer)
	grpc.SendHeader(ctx, header)

	//获取请求中的元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("server recv UnaryEcho 获取元数据失败")
	} else {
		fmt.Println("server recv UnaryEcho 获取元数据 time :", md.Get("time"))
	}

	fmt.Printf("server recv : %v\n", in.Message)
	return &echo.EchoResponse{
		Message: "server send message",
	}, nil
}
func (EchoServer) ServerStreamingEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	// 响应请求，发送元数据
	header, trailer := getMetadata()
	// trailer,服务端调用结束后填充的数据
	defer stream.SetTrailer(trailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	//获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("server recv ServerStreamingEcho 获取元数据失败")
	} else {
		fmt.Println("server recv ServerStreamingEcho 获取元数据 :", md)
	}

	fmt.Printf("server recv : %v\n", in.Message)
	filePath := "echo-server/files/server.jpg"
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
			return err
		}
		stream.Send(&echo.EchoResponse{
			Message: "server sending file",
			Bytes:   buf[:n],
			Time:    timestamppb.New(time.Now()),
			Length:  int32(n),
		})
	}
	// 服务端return nil 或者 error 即表示流结束
	return nil
}
func (EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	// 响应请求，发送元数据
	header, trailer := getMetadata()
	// trailer,服务端调用结束后填充的数据
	defer stream.SetTrailer(trailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	//获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("server recv ClientStreamingEcho 获取元数据失败")
	} else {
		fmt.Println("server recv ClientStreamingEcho 获取元数据 :", md)
	}

	filePath := "echo-server/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		file.Write(req.Bytes[:req.Length])
		//fmt.Printf("server recv : %v\n", req.Message)
	}
	err = stream.SendAndClose(&echo.EchoResponse{
		Message: "server send message",
	})
	return err
}
func (EchoServer) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	// 响应请求，发送元数据
	header, trailer := getMetadata()
	// trailer,服务端调用结束后填充的数据
	defer stream.SetTrailer(trailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	//获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("server recv BidirectionalStreamingEcho 获取元数据失败")
	} else {
		fmt.Println("server recv BidirectionalStreamingEcho 获取元数据 :", md)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//接收客户端流，保存文件
		filePath := "echo-server/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".jpg"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return
			}
			file.Write(req.Bytes[:req.Length])
			//fmt.Printf("server recv : %v\n", req.Message)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath := "echo-server/files/server.jpg"
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
				return
			}
			stream.Send(&echo.EchoResponse{
				Message: "server sending file",
				Bytes:   buf[:n],
				Time:    timestamppb.New(time.Now()),
				Length:  int32(n),
			})
		}
	}()
	wg.Wait()
	// 服务端return nil 或者 error 即表示流结束
	return nil
}
