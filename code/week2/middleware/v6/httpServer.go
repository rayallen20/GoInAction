package web

import (
	"fmt"
	"net"
	"net/http"
)

// 为确保HTTPServer结构体为Server接口的实现而定义的变量
var _ Server = &HTTPServer{}

// HTTPServer HTTP服务器
type HTTPServer struct {
	router                                   // router 路由树
	middlewares []Middleware                 // middlewares 中间件切片.表示HTTPServer需要按顺序执行的的中间件链
	logFunc     func(msg string, arg ...any) // logFunc 日志函数
}

// NewHTTPServer 创建HTTP服务器
// 这里选项的含义其实是指不同的 Option 函数
// 每一个 Option 函数都会对 HTTPServer 实例的不同属性进行设置
func NewHTTPServer(opts ...Option) *HTTPServer {
	server := &HTTPServer{
		router: newRouter(),
		logFunc: func(msg string, arg ...any) {
			fmt.Printf(msg, arg...)
		},
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

// ServerWithMiddleware 本函数用于为 HTTPServer 实例添加中间件
// 即: 本函数用于设置 HTTPServer 实例的middlewares属性
func ServerWithMiddleware(middlewares ...Middleware) Option {
	return func(server *HTTPServer) {
		server.middlewares = middlewares
	}
}

// ServeHTTP WEB框架入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构建上下文
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 执行中间件链
	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		root = s.middlewares[i](root)
	}

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	// 最后注册将响应数据和响应码写入到响应体中的中间件
	// 确保这个中间件是执行完所有对响应码和响应数据的读写操作后才执行的
	// 换言之,确保这个中间件是返回响应之前最后一个执行的
	root = m(root)

	// 查找路由树并执行命中的业务逻辑
	root(ctx)
}

// flashResp 将响应数据和响应码写入到响应体中
func (s *HTTPServer) flashResp(ctx *Context) {
	// 若使用者设置了响应码 则刷到响应上
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}

	// 刷响应数据到响应上
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		s.logFunc("响应数据写入失败: %v", err)
	}

	if n != len(ctx.RespData) {
		s.logFunc("响应数据写入不完全, 期望写入: %d 字节, 实际写入: %d 字节", len(ctx.RespData), n)
	}
}

// serve 查找路由树并执行命中的业务逻辑
func (s *HTTPServer) serve(ctx *Context) {
	method := ctx.Req.Method
	path := ctx.Req.URL.Path
	targetNode, ok := s.findRoute(method, path)
	// 没有在路由树中找到对应的路由节点 或 找到了路由节点的处理函数为空(即NPE:none pointer exception 的问题)
	// 则返回404
	if !ok || targetNode.node.HandleFunc == nil {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("Not Found")
		return
	}

	// 命中节点则将路径参数名值对设置到上下文中
	ctx.PathParams = targetNode.pathParams

	// 命中节点则将节点的路由设置到上下文中
	ctx.MatchRoute = targetNode.node.route
	// 执行路由节点的处理函数
	targetNode.node.HandleFunc(ctx)
}

// Start 启动WEB服务器
func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 在监听端口之后,启动服务之前做一些操作
	// 例如在微服务框架中,启动服务之前需要注册服务

	return http.Serve(l, s)
}

// GET 注册GET请求路由
func (s *HTTPServer) GET(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodGet, path, handleFunc)
}

// POST 注册POST请求路由
func (s *HTTPServer) POST(path string, handleFunc HandleFunc) {
	s.addRoute(http.MethodPost, path, handleFunc)
}

// Use 执行路由匹配 仅当匹配到路由时 才执行中间件
func (s *HTTPServer) Use(path string, middlewares ...Middleware) {
	panic("implement me")
}
