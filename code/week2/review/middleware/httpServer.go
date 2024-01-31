package web

import (
	"net"
	"net/http"
)

// 为确保HTTPServer结构体为Server接口的实现而定义的变量
var _ Server = &HTTPServer{}

// HTTPServer HTTP服务器
type HTTPServer struct {
	router                   // router 路由树
	middlewares []Middleware // middlewares Server级别的中间件链 实际上就是责任链 所有的请求都会经过这个链的处理
}

// NewHTTPServer 根据给定的 Option 列表(每个 Option 均表示要修改一个 HttpServer 的成员属性),创建HTTP服务器
func NewHTTPServer(options ...Option) *HTTPServer {
	httpServer := &HTTPServer{
		router: newRouter(),
	}

	for _, option := range options {
		option(httpServer)
	}

	return httpServer
}

// ServerWithMiddlewares 本函数用于根据给定的 Middleware 列表,创建
// 修改 HttpServer 实例的 middlewares 字段值的选项函数
func ServerWithMiddlewares(middlewares ...Middleware) Option {
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

	// 构建责任链
	// step1. 找到请求对应的HandleFunc
	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		// step2. 从后往前构建责任链
		// Tips: 组装的过程是从后向前的 也就是说最后一个被执行的中间件是最先被组装到责任链上的
		root = s.middlewares[i](root)
	}

	// 从责任链的头部开始执行
	root(ctx)
}

// serve 查找路由树并执行命中的业务逻辑
func (s *HTTPServer) serve(ctx *Context) {
	method := ctx.Req.Method
	path := ctx.Req.URL.Path
	targetNode, ok := s.findRoute(method, path)
	if !ok || targetNode.node.HandleFunc == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("Not Found"))
		return
	}

	ctx.PathParams = targetNode.pathParams

	// 命中节点则将节点的全路由设置到上下文中
	ctx.MatchRoute = targetNode.node.fullRoute

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
