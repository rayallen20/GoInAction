package main

import (
	"fmt"
	"option/pattern"
)

func main() {
	// 初始化SomeOption时,若SomeOption的字段有变化,仅需调整options切片即可
	options := []pattern.OptionFunc{
		pattern.WithB(10),
		pattern.WithC(true),
	}

	someOption := pattern.NewSomeOption("a", options...)
	fmt.Printf("%#v\n", someOption)
}
