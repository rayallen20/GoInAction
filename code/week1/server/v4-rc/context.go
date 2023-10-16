package v4_rc

import "net/http"

// Context 路由处理函数的上下文
type Context struct {
	// Req HTTP请求
	Req *http.Request
	// Resp HTTP响应
	Resp http.ResponseWriter
}
