package web

import "strconv"

// ReqValue 用于承载来自请求中各部分输入的值 并提供统一的类型转换API
type ReqValue struct {
	value string // value 来自请求中不同部分的值 以string类型表示
	err   error  // err 承载接收请求中的参数时出现的错误
}

// AsInt64 将ReqValue中的值转换为int64类型
func (r ReqValue) AsInt64() (value int64, err error) {
	if r.err != nil {
		return 0, r.err
	}

	return strconv.ParseInt(r.value, 10, 64)
}

// AsUint64 将ReqValue中的值转换为uint64类型
func (r ReqValue) AsUint64() (value uint64, err error) {
	if r.err != nil {
		return 0, r.err
	}

	return strconv.ParseUint(r.value, 10, 64)
}

// AsFloat64 将ReqValue中的值转换为float64类型
func (r ReqValue) AsFloat64() (value float64, err error) {
	if r.err != nil {
		return 0, r.err
	}

	return strconv.ParseFloat(r.value, 64)
}
