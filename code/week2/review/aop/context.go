package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Context HandleFunc的上下文
type Context struct {
	Req        *http.Request       // Req 请求
	Resp       http.ResponseWriter // Resp 响应
	PathParams map[string]string   // PathParams 路径参数名值对
	queryValue url.Values          // queryValue 查询参数名值对
}

// BindJSON 绑定请求体中的JSON数据到给定的目标对象上 这个目标对象可能是某个结构体的实例 也有可能是个map
func (c *Context) BindJSON(target any) error {
	if target == nil {
		return errors.New("web绑定错误: 给定的实例为空")
	}

	if c.Req.Body == nil {
		return errors.New("web绑定错误: 请求体为空")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(target)
}

// FormValue 获取表单中给定的key对应的值
func (c *Context) FormValue(key string) (value ReqValue) {
	err := c.Req.ParseForm()
	if err != nil {
		return ReqValue{err: err}
	}

	_, ok := c.Req.Form[key]
	if !ok {
		return ReqValue{err: errors.New("web绑定错误: 表单中没有给定的key: " + key)}
	}

	return ReqValue{value: c.Req.FormValue(key)}
}

// QueryValue 获取URL中给定的key对应的值
func (c *Context) QueryValue(key string) (value ReqValue) {
	if c.queryValue == nil {
		c.queryValue = c.Req.URL.Query()
	}

	values, ok := c.queryValue[key]
	if !ok {
		return ReqValue{err: errors.New("web绑定错误: URL中没有给定的key: " + key)}
	}

	if len(values) == 0 {
		return ReqValue{err: errors.New("web绑定错误: URL中给定的key没有对应的值: " + key)}
	}

	return ReqValue{value: values[0]}
}

// PathValue 获取路径参数中给定的key对应的值
func (c *Context) PathValue(key string) (value ReqValue) {
	if c.PathParams == nil {
		return ReqValue{err: errors.New("web绑定错误: 路径参数为空")}
	}

	val, ok := c.PathParams[key]
	if !ok {
		return ReqValue{err: errors.New("web绑定错误: 路径中没有给定的key: " + key)}
	}

	return ReqValue{value: val}
}
