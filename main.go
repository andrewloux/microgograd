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
	b := micrograd.NewValue(3.0, micrograd.WithName("b"))
	c := a.Mul(b).(*micrograd.Value[float64])
	c.Name = "c"
	d := c.Add(a).(*micrograd.Value[float64])
	d.Name = "result"

	// Generate DOT file with default node labels
	err := plot.WriteGraph(d, "graph.dot")
	if err != nil {
		log.Fatalf("Error generating graph: %v", err)
	}

	// Generate interactive HTML version
	err = plot.WriteInteractiveHTML(d, "graph.html")
	if err != nil {
		log.Fatalf("Error generating interactive graph: %v", err)
	}

	fmt.Println("Generated graph.dot")
	fmt.Println("To visualize static graph, run:")
	fmt.Println("dot -Tsvg graph.dot -o graph.svg")
	fmt.Println("\nGenerated interactive graph.html")
	fmt.Println("Open graph.html in a browser to interact with the graph")
}
