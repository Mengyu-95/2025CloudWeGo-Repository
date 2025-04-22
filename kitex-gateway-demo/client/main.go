package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudwego/kitex/client"
	"github.com/yourusername/kitex-gateway-demo/kitex_gen/api"
	"github.com/yourusername/kitex-gateway-demo/kitex_gen/api/userservice"
	// "github.com/yourusername/kitex-gateway-demo/kitex_gen/api/userservice"
)

func main() {
	// 测试HTTP请求
	testHTTPRequest()

	// 测试Thrift请求
	testThriftRequest()
}

func testHTTPRequest() {
	// 发送HTTP GET请求
	resp, err := http.Get("http://localhost:8888/api/UserService/GetUser?user_id=123")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("HTTP Response:", string(body))

	// 发送HTTP POST请求
	reqBody := map[string]interface{}{"user_id": 123}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err = http.Post("http://localhost:8888/api/UserService/GetUser", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("HTTP POST Response:", string(body))
}

func testThriftRequest() {
	// 创建Thrift客户端
	c, err := userservice.NewClient("UserService", client.WithHostPorts("localhost:8888"))
	if err != nil {
		panic(err)
	}

	// 发送Thrift请求
	req := &api.GetUserRequest{UserId: 123}
	resp, err := c.GetUser(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Thrift Response: %+v\n", resp)
}
