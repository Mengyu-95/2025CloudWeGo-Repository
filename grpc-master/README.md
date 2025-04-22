# grpc 
## 无证书的情况下健康检查
``` 
.\grpc_health_probe-windows-amd64.exe -addr localhost:50051
```
## TLS 健康检查
``` 
.\grpc_health_probe-windows-amd64.exe -tls -tls-ca-cert .\x509\ca_cert.pem -tls-server-name echo.grpc.0voice.com -addr localhost:50051
```
## MTLS 健康检查
``` 
.\grpc_health_probe-windows-amd64.exe -tls -tls-ca-cert .\x509\ca_cert.pem -tls-client-cert .\x509\client_cert.pem -tls-client-key .\x509\client_key.pem -tls-server-name echo.grpc.0voice.com -addr localhost:50051
```