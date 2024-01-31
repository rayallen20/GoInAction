package accessLog

// accessLog 本结构体用于定义日志内容
type accessLog struct {
	Host       string `json:"host,omitempty"`        // Host 请求的主机地址
	Route      string `json:"route,omitempty"`       // Route 命中的路由
	HTTPMethod string `json:"http_method,omitempty"` // HTTPMethod 请求的HTTP动词
	Path       string `json:"path,omitempty"`        // Path 请求的路径 即:uri
}
