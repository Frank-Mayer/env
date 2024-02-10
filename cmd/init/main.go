package main

import (
	"fmt"

	"github.com/Frank-Mayer/env/internal"
)

func main() {
	ok, err := internal.Init()
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("Aborted")
		return
	}
	fmt.Println("You're good to go!")
}
