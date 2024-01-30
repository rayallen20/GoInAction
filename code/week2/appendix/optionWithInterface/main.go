package main

import "optionWithInterface/pattern"

func main() {
	options := []pattern.IOption{
		pattern.WithB(10),
		pattern.WithC(true),
	}

	pattern.NewSomeOption("a", options...)
}
