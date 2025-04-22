package client

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func GetInsecureOpt() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func GetTlsOpt() grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile("x509/ca_cert.pem", "echo.grpc.0voice.com")
	if err != nil {
		log.Fatal(err)
	}
	opt := grpc.WithTransportCredentials(creds)
	return opt

}

func GetMTlsOpt() grpc.DialOption {
	cert, err := tls.LoadX509KeyPair("x509/client_cert.pem", "x509/client_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	ca := x509.NewCertPool()
	caFilePath := "x509/ca_cert.pem"
	bytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if ok := ca.AppendCertsFromPEM(bytes); !ok {
		log.Fatal("ca append failed")
	}
	tlsConfig := &tls.Config{
		ServerName:   "echo.grpc.0voice.com",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
}
