package micrograd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_Basic(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		gradient float64
	}{
		{
			name:     "simple value",
			value:    2.0,
			gradient: 0.0,
		},
		{
			name:     "with gradient",
			value:    3.0,
			gradient: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValue(tt.value, WithGradient(tt.gradient))
			assert.Equal(t, tt.value, v.GetValue())
			assert.Equal(t, tt.gradient, v.GetGradient())
		})
	}
}

func TestValue_Operations(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		a := NewValue(2.0)
		b := NewValue(3.0)
		c := a.Add(b)
		assert.Equal(t, 5.0, c.GetValue())
	})

	t.Run("multiplication", func(t *testing.T) {
		a := NewValue(2.0)
		b := NewValue(3.0)
		c := a.Mul(b)
		assert.Equal(t, 6.0, c.GetValue())
	})

	t.Run("subtraction", func(t *testing.T) {
		a := NewValue(5.0)
		b := NewValue(3.0)
		c := a.Sub(b)
		assert.Equal(t, 2.0, c.GetValue())
	})
}

func TestValue_Backpropagation(t *testing.T) {
	// Test case: f(a,b) = a * b + b
	// df/da = b
	// df/db = a + 1
	a := NewValue(2.0, WithName("a"))
	b := NewValue(3.0, WithName("b"))

	c := a.Mul(b).SetName("c") // c = a * b
	d := c.Add(b).SetName("d") // d = c + b = (a * b) + b

	d.SetGradient(1.0)
	d.Backtrack()

	assert.Equal(t, 3.0, a.GetGradient()) // df/da = b = 3
	assert.Equal(t, 3.0, b.GetGradient()) // df/db = a + 1 = 3
}

func TestValue_Options(t *testing.T) {
	t.Run("with name", func(t *testing.T) {
		v := NewValue(2.0, WithName("test"))
		assert.Equal(t, "test", v.GetName())
	})

	t.Run("with gradient", func(t *testing.T) {
		v := NewValue(2.0, WithGradient(1.0))
		assert.Equal(t, 1.0, v.GetGradient())
	})

	t.Run("with value", func(t *testing.T) {
		v := NewValue(2.0, WithValue(3.0))
		assert.Equal(t, 3.0, v.GetValue())
	})
}
