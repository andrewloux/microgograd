package main

import (
	"fmt"
	"log"

	"microgograd/examples/manual_backprop"
)

func main() {
	fmt.Println("Manual Backpropagation Example")
	fmt.Println("=============================")

	if err := manual_backprop.Run(); err != nil {
		log.Fatalf("Error running example: %v", err)
	}
}
