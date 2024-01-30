package pattern

import "fmt"

// NewSomeOption 本函数是 someOption 的构造函数
func NewSomeOption(a string, options ...IOption) {
	someOptionObj := &someOption{
		a: a,
	}

	for _, option := range options {
		option.apply(someOptionObj)
	}

	fmt.Printf("%#v\n", someOptionObj)
}
