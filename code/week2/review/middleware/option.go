package web

// Option 本类型为 HttpServer 的选项函数
// 本类型的每个不同实例均用于修改 HttpServer 的不同字段值
type Option func(server *HTTPServer)
