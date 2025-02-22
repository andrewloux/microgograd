package manual_backprop

import (
	"fmt"

	"microgograd/micrograd"
	"microgograd/plot"
)

func Run() error {
	// Create a simple computation graph
	a := micrograd.NewValue(2.0, micrograd.WithName("a"))
	b := micrograd.NewValue(-3.0, micrograd.WithName("b"))
	c := micrograd.NewValue(10.0, micrograd.WithName("c"))

	e := a.Mul(b).SetName("e")
	d := e.Add(c).SetName("d")

	f := micrograd.NewValue(-2.0).SetName("f")
	L := d.Mul(f).SetName("L")

	L.SetGradient(1.0)
	L.Backtrack()

	// Generate interactive HTML version
	err := plot.WriteInteractiveHTML(L, "graph.html")
	if err != nil {
		return fmt.Errorf("error generating interactive graph: %v", err)
	}

	fmt.Println("Generated graph.html")
	fmt.Println("Open graph.html in a browser to interact with the graph")
	return nil
}
