package client

import (
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func GetAuth(token string) grpc.DialOption {
	perRPC := GetPerRPCCredentials(token)
	return grpc.WithPerRPCCredentials(perRPC)
}

func GetPerRPCCredentials(token string) credentials.PerRPCCredentials {
	return oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
}

func FetchToken() string {
	return "some-secret-token"
}
