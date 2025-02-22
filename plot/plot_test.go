package plot

import (
	"os"
	"strings"
	"testing"

	"microgograd/micrograd"

	"github.com/stretchr/testify/assert"
)

func TestWriteInteractiveHTML(t *testing.T) {
	t.Run("basic graph", func(t *testing.T) {
		// Create a simple graph
		a := micrograd.NewValue(2.0, micrograd.WithName("a"))
		b := micrograd.NewValue(3.0, micrograd.WithName("b"))
		c := a.Mul(b).(*micrograd.Value[float64])
		c.SetName("c")

		// Generate HTML
		err := WriteInteractiveHTML(c, "test_graph.html")
		assert.NoError(t, err)

		// Read generated file
		data, err := os.ReadFile("test_graph.html")
		assert.NoError(t, err)
		content := string(data)

		// Check for essential elements
		assert.Contains(t, content, `<!DOCTYPE html>`)
		assert.Contains(t, content, `cytoscape`)
		assert.Contains(t, content, `Name: a`)
		assert.Contains(t, content, `Name: b`)
		assert.Contains(t, content, `Name: c`)

		// Clean up
		os.Remove("test_graph.html")
	})
}

func TestNodeLabel(t *testing.T) {
	t.Run("default label format", func(t *testing.T) {
		v := micrograd.NewValue(2.5, micrograd.WithName("test"))
		v.SetGradient(1.0)

		label := defaultNodeLabel(v)
		assert.Contains(t, label, "Name: test")
		assert.Contains(t, label, "Value: 2.5")
		assert.Contains(t, label, "Grad: 1.0")
	})

	t.Run("custom label function", func(t *testing.T) {
		v := micrograd.NewValue(2.5, micrograd.WithName("test"))
		customLabel := func(node micrograd.Numeric[float64]) string {
			return "Custom: " + node.GetName()
		}

		err := WriteInteractiveHTML(v, "test_graph.html", WithNodeLabelFunc(customLabel))
		assert.NoError(t, err)

		data, err := os.ReadFile("test_graph.html")
		assert.NoError(t, err)
		content := string(data)

		assert.Contains(t, content, "Custom: test")
		assert.NotContains(t, content, "Value: 2.5")

		// Clean up
		os.Remove("test_graph.html")
	})
}

func TestDOTGeneration(t *testing.T) {
	t.Run("operation nodes", func(t *testing.T) {
		a := micrograd.NewValue(2.0, micrograd.WithName("a"))
		b := micrograd.NewValue(3.0, micrograd.WithName("b"))
		c := a.Mul(b).(*micrograd.Value[float64])
		c.SetName("c")

		cfg := &plotConfig[float64]{
			labelFunc: defaultNodeLabel[float64],
		}
		dot := dotFromValue(c, cfg)

		// Check for multiplication operation node
		assert.Contains(t, dot, `shape="ellipse"`)
		assert.Contains(t, dot, string(rune(micrograd.MUL)))

		// Check for edges
		assert.True(t, strings.Contains(dot, "->"), "DOT should contain edges")
	})
}
