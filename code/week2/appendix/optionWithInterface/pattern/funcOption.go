package pattern

// funcOption 函数选项类型 因为要实现接口,所以不能再使用函数类型,只能使用结构体
type funcOption struct {
	// f 具体选项函数
	withFunc func(someOption *someOption)
}

func (o *funcOption) apply(someOption *someOption) {
	o.withFunc(someOption)
}

// newFuncOption funcOption 的构造函数
func newFuncOption(withFunc func(someOption *someOption)) *funcOption {
	return &funcOption{withFunc: withFunc}
}
