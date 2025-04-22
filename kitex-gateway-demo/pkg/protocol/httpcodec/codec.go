package httpcodec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/yourusername/kitex-gateway-demo/kitex_gen/api"
	// "github.com/cloudwego/hertz/pkg/binding"
)

type httpCodec struct{}

// Name implements remote.Codec.
func (c *httpCodec) Name() string {
	panic("unimplemented")
}

func NewHTTPCodec() remote.Codec {
	return &httpCodec{}
}

type Message struct {
	// 其他字段
	serviceName string
	method      string
	data        interface{}
}

func (m *Message) SetServiceName(name string) {
	m.serviceName = name
}
func (m *Message) SetMethod(method string) {
	m.method = method
}
func (m *Message) SetData(data interface{}) {
	m.data = data
}

func (c *httpCodec) Encode(ctx context.Context, msg remote.Message, out remote.ByteBuffer) error {
	// 实现HTTP响应编码
	resp := msg.Data().(*Response)

	// 设置HTTP状态码
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	// 编码JSON响应体
	jsonData, err := json.Marshal(resp.Body)
	if err != nil {
		return err
	}

	// 写入HTTP响应
	_, err = out.WriteString("HTTP/1.1 ")
	_, err = out.WriteString(http.StatusText(resp.StatusCode))
	_, err = out.WriteString("\r\n")
	_, err = out.WriteString("Content-Type: application/json\r\n")
	_, err = out.WriteString("\r\n")
	_, err = out.Write(jsonData)

	return err
}

func (c *httpCodec) Decode(ctx context.Context, msg remote.Message, in remote.ByteBuffer) error {
	// 解析HTTP请求
	req := &Request{
		Header: make(http.Header),
	}

	// 解析HTTP方法、路径和版本
	line, err := in.ReadString('\n')
	if err != nil {
		return err
	}
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) < 3 {
		return errors.New("invalid http request line")
	}
	req.Method = parts[0]
	req.URL.Path = parts[1]

	// 解析HTTP头
	for {
		line, err := in.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) == 2 {
			req.Header.Add(parts[0], strings.TrimSpace(parts[1]))
		}
	}

	// 解析HTTP体
	if req.Method == "POST" || req.Method == "PUT" {
		contentLength := req.Header.Get("Content-Length")
		if contentLength != "" {
			length, _ := strconv.Atoi(contentLength)
			body := make([]byte, length)
			_, err := in.Read(body)
			if err != nil {
				return err
			}
			req.Body = body
		}
	}

	// 将HTTP请求映射到Thrift方法
	pathParts := strings.Split(strings.TrimPrefix(req.URL.Path, "/api/"), "/")
	if len(pathParts) < 2 {
		return errors.New("invalid http path")
	}

	svcName := pathParts[0]
	methodName := pathParts[1]

	msg.SetServiceName(svcName)
	msg.SetMethod(methodName)

	// 解析请求体到Thrift结构
	var thriftReq interface{}
	switch svcName {
	case "UserService":
		switch methodName {
		case "GetUser":
			thriftReq = &api.GetUserRequest{}
		}
	}

	if thriftReq != nil {
		if err := binding.BindAndValidate(req, thriftReq); err != nil {
			return err
		}
		msg.SetData(thriftReq)
	}

	return nil
}

type Request struct {
	Method string
	URL    struct {
		Path string
	}
	Header http.Header
	Body   []byte
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       interface{}
}
