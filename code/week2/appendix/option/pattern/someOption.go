package pattern

type SomeOption struct {
	A string
	B int
	C bool
}

const (
	// defaultValueB SomeOption 实例的 B 字段的默认值
	defaultValueB = 100
)

// NewSomeOption SomeOption 的构造函数
func NewSomeOption(a string, options ...OptionFunc) *SomeOption {
	someOption := &SomeOption{
		A: a,
		B: defaultValueB,
	}

	// 若客户端提供了修改B字段的选项函数 则在此处会覆盖掉默认值
	for _, option := range options {
		option(someOption)
	}

	return someOption
}
