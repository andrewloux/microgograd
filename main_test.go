package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"microgograd/micrograd"
	"microgograd/plot"
)

func TestDOTOutput(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *micrograd.Value[float64]
		wantDot []string
		wantNot []string
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
			wantDot: []string{
				`label="{a | 2}"`,
				`label="{b | -3}"`,
				`label="{c | -6 | op: *}"`,
				`->`,
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
			wantDot: []string{
				`label="{y | 3}"`,
				`label="{z | 4}"`,
				`label="{out |`, // partial match since value depends on evaluation
				`op: +`,         // addition operation
				`op: *`,         // multiplication operation
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.setup()
			err := plot.WriteGraph[float64](v, "test_graph.dot")
			assert.NoError(t, err)

			// Read the generated file
			data, err := os.ReadFile("test_graph.dot")
			assert.NoError(t, err)
			got := string(data)

			for _, want := range tt.wantDot {
				assert.Contains(t, got, want)
			}
			for _, notWant := range tt.wantNot {
				assert.NotContains(t, got, notWant)
			}
		})
	}
}

func TestDOTFormatting(t *testing.T) {
	t.Run("node_style", func(t *testing.T) {
		v := micrograd.NewValue(1.0, micrograd.WithName("test"))
		err := plot.WriteGraph[float64](v, "test_graph.dot")
		assert.NoError(t, err)

		data, err := os.ReadFile("test_graph.dot")
		assert.NoError(t, err)
		dot := string(data)

		assert.Contains(t, dot, `shape="record"`)
		assert.Contains(t, dot, `style="filled"`)
		assert.Contains(t, dot, `fillcolor="white"`)
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

		err := plot.WriteGraph[float64](d, "test_graph.dot")
		assert.NoError(t, err)

		data, err := os.ReadFile("test_graph.dot")
		assert.NoError(t, err)
		dot := string(data)

		// Look for label definitions to check order
		aIdx := strings.Index(dot, `label="{a |`)
		cIdx := strings.Index(dot, `label="{c |`)
		dIdx := strings.Index(dot, `label="{d |`)

		assert.True(t, aIdx < cIdx, "a should come before c")
		assert.True(t, cIdx < dIdx, "c should come before d")
	})
}
