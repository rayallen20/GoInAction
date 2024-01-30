package pattern

// IOption 本接口用于定义函数选项行为
type IOption interface {
	// apply 本方法用于修改 someOption 结构体的字段值
	apply(someOption *someOption)
}
