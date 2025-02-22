package main

import (
	"fmt"
	"log"
	"os"

	"microgograd/examples/manual_backprop"
)

func main() {
	fmt.Println("MicroGrad Examples")
	fmt.Println("=================")
	fmt.Println("\nAvailable examples:")
	fmt.Println("1. Manual Backpropagation")
	fmt.Println("\nSelect an example (1) or press Ctrl+C to exit:")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		fmt.Println("\nRunning Manual Backpropagation Example...")
		fmt.Println("(To run directly: go run cmd/manual_backprop/main.go)")
		fmt.Println()

		if err := manual_backprop.Run(); err != nil {
			log.Fatalf("Error running example: %v", err)
		}
	default:
		fmt.Println("Invalid choice")
		os.Exit(1)
	}
}
