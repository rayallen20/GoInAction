package v4_rc

import (
	"net/http"
	"testing"
)

// TestServer_Start 测试服务器启动
func TestServer_Start(t *testing.T) {
	s := NewHTTPServer()
	handleFunc := func(ctx *Context) {
		ctx.Resp.Write([]byte("hello order detail"))
	}
	s.addRoute(http.MethodGet, "/order/detail", handleFunc)
	err := s.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}
