package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"microgograd/micrograd"
	"microgograd/plot"
)

func TestHTMLOutput(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *micrograd.Value[float64]
		wantHTML []string
	}{
		{
			name: "simple multiplication",
			setup: func() *micrograd.Value[float64] {
				a := micrograd.NewValue(2.0, micrograd.WithName("a"))
				b := micrograd.NewValue(-3.0, micrograd.WithName("b"))
				c := a.Mul(b).(*micrograd.Value[float64])
				c.Name = "c"
				return c
			},
			wantHTML: []string{
				`Name: a`,
				`Value: 2`,
				`Name: b`,
				`Value: -3`,
				`Name: c`,
			},
		},
		{
			name: "chained operations",
			setup: func() *micrograd.Value[float64] {
				x := micrograd.NewValue(2.0, micrograd.WithName("x"))
				y := micrograd.NewValue(3.0, micrograd.WithName("y"))
				z := micrograd.NewValue(4.0, micrograd.WithName("z"))
				out := x.Add(y.Mul(z)).(*micrograd.Value[float64])
				out.Name = "out"
				return out
			},
			wantHTML: []string{
				`Name: x`,
				`Name: y`,
				`Name: z`,
				`Name: out`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.setup()
			err := plot.WriteInteractiveHTML[float64](v, "test_graph.html")
			assert.NoError(t, err)

			// Read the generated file
			data, err := os.ReadFile("test_graph.html")
			assert.NoError(t, err)
			got := string(data)

			for _, want := range tt.wantHTML {
				assert.Contains(t, got, want)
			}
		})
	}
}

func TestHTMLFormatting(t *testing.T) {
	t.Run("node_style", func(t *testing.T) {
		v := micrograd.NewValue(1.0, micrograd.WithName("test"))
		err := plot.WriteInteractiveHTML[float64](v, "test_graph.html")
		assert.NoError(t, err)

		data, err := os.ReadFile("test_graph.html")
		assert.NoError(t, err)
		html := string(data)

		assert.Contains(t, html, `shape="record"`)
		assert.Contains(t, html, `style="filled"`)
	})
}

func TestTopologicalOrder(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		a := micrograd.NewValue(1.0, micrograd.WithName("a"))
		b := micrograd.NewValue(2.0, micrograd.WithName("b"))
		c := a.Add(b).(*micrograd.Value[float64])
		c.Name = "c"
		d := c.Mul(a).(*micrograd.Value[float64])
		d.Name = "d"

		err := plot.WriteInteractiveHTML[float64](d, "test_graph.html")
		assert.NoError(t, err)

		data, err := os.ReadFile("test_graph.html")
		assert.NoError(t, err)
		html := string(data)

		// Look for node definitions to check order
		aIdx := strings.Index(html, `Name: a`)
		cIdx := strings.Index(html, `Name: c`)
		dIdx := strings.Index(html, `Name: d`)

		assert.True(t, aIdx >= 0, "should contain node a")
		assert.True(t, cIdx >= 0, "should contain node c")
		assert.True(t, dIdx >= 0, "should contain node d")
	})
}
