package web

// Middleware 中间件
type Middleware func(HandleFunc) HandleFunc
