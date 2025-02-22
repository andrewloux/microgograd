package main

import (
	"fmt"

	tensor "gorgonia.org/tensor"
)

func main() {

	fmt.Println("Hello World")
	a := tensor.New(tensor.WithShape(1, 1), tensor.WithBacking([]int{1}))

}
