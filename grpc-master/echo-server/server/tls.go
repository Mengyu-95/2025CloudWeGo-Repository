package server

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

func GetTlsOpt() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile("x509/server_cert.pem", "x509/server_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	return grpc.Creds(creds)
}

func GetMTlsOpt() grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair("x509/server_cert.pem", "x509/server_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	ca := x509.NewCertPool()
	caFilePath := "x509/client_ca_cert.pem"
	bytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if ok := ca.AppendCertsFromPEM(bytes); !ok {
		log.Fatal("ca append failed")
	}
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}
