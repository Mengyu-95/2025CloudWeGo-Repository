package server

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("server UnaryInterceptor")
	fmt.Println(info)
	if info.FullMethod != "/grpc.health.v1.Health/Check" {
		err = oauth2Valid(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}

func StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Println("server StreamInterceptor")
	fmt.Println(info)
	err := oauth2Valid(ss.Context())
	if err != nil {
		return err
	}
	return handler(srv, ss)
}

func oauth2Valid(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("元数据获取失败，身份认证失败")
	}
	authorization := md["authorization"]
	if !valid(authorization) {
		return errors.New("身份令牌校验失败，身份认证失败")
	}
	return nil
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == fetchToken()
}
func fetchToken() string {
	return "some-secret-token"
}
