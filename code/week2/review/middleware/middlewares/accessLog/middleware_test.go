package accessLog

import (
	"fmt"
	"testing"
	"web"
)

// Test_Middleware 本函数用于测试 accessLog 是否工作正常
func Test_Middleware(t *testing.T) {
	// step1. 创建中间件
	logFunc := func(content string) {
		fmt.Printf("%#v\n", content)
	}
	middlewareBuilder := AccessMiddlewareBuilder{
		logFunc: logFunc,
	}
	accessLogMiddleware := middlewareBuilder.Build()

	// step2. 创建HTTPServer
	middlewareOption := web.ServerWithMiddlewares(accessLogMiddleware)
	httpServer := web.NewHTTPServer(middlewareOption)

	// step3. 启动HTTPServer
	handleFunc := func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello"))
	}
	httpServer.GET("/user/show", handleFunc)

	httpServer.Start(":8081")
}
