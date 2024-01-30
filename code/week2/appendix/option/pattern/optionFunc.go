package pattern

// OptionFunc 函数选项 该类型用于设置 SomeOption 的各属性值(即各选项)
type OptionFunc func(*SomeOption)

// WithB 本函数返回一个 OptionFunc ,该 OptionFunc 将
// SomeOption 实例的 B 字段值设置为给定值
func WithB(b int) OptionFunc {
	return func(someOption *SomeOption) {
		someOption.B = b
	}
}

// WithC 本函数返回一个 OptionFunc ,该 OptionFunc 将
// SomeOption 实例的 C 字段值设置为给定值
func WithC(c bool) OptionFunc {
	return func(someOption *SomeOption) {
		someOption.C = c
	}
}
