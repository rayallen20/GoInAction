package pattern

// WithB 本函数用于修改 someOption 实例的 b 字段值
func WithB(b int) IOption {
	withFunc := func(someOption *someOption) {
		someOption.b = b
	}

	return &funcOption{withFunc: withFunc}
}

// WithC 本函数用于修改 someOption 实例的 c 字段值
func WithC(c bool) IOption {
	withFunc := func(someOption *someOption) {
		someOption.c = c
	}

	return &funcOption{withFunc: withFunc}
}
