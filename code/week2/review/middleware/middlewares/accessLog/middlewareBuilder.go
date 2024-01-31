package accessLog

import (
	"encoding/json"
	"web"
)

// AccessMiddlewareBuilder 日志中间件构建器
type AccessMiddlewareBuilder struct {
	logFunc func(content string) // logFunc 用于记录日志的函数
}

// SetLogFunc 本方法用于设置记录日志的函数
func (b *AccessMiddlewareBuilder) SetLogFunc(logFunc func(string)) {
	b.logFunc = logFunc
}

// Build 本方法用于构建一个日志中间件
func (b *AccessMiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 构建日志内容
			defer func() {
				log := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}

				// 记录日志
				logJsonBytes, _ := json.Marshal(log)
				b.logFunc(string(logJsonBytes))
			}()
			next(ctx)
		}
	}
}
