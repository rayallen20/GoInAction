package main

import "github.com/google/uuid"

func main() {
	for i := 0; i < 3; i++ {
		uuid := uuid.New()
		println(uuid.String())
	}
}
