package main

import (
	"fmt"
	"log"

	"microgograd/micrograd"
	"microgograd/plot"
)

func main() {
	// Create a simple computation graph
	a := micrograd.NewValue(2.0, micrograd.WithName("a"))
	b := micrograd.NewValue(-3.0, micrograd.WithName("b"))
	c := micrograd.NewValue(10.0, micrograd.WithName("c"))

	e := a.Mul(b).SetName("c")
	d := e.Add(c).SetName("d")

	f := micrograd.NewValue(-2.0).SetName("f")
	L := d.Mul(f).SetName("L")

	L.SetGradient(1.0)

	// Generate interactive HTML version
	err := plot.WriteInteractiveHTML(L, "graph.html")
	if err != nil {
		log.Fatalf("Error generating interactive graph: %v", err)
	}

	fmt.Println("Generated graph.dot")
	fmt.Println("To visualize static graph, run:")
	fmt.Println("dot -Tsvg graph.dot -o graph.svg")
	fmt.Println("\nGenerated interactive graph.html")
	fmt.Println("Open graph.html in a browser to interact with the graph")
}
