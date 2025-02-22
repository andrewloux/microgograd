package micrograd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair_Basic(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		var zero float64
		p := NewPair(zero, zero)
		assert.Equal(t, zero, p.X())
		assert.Equal(t, zero, p.Y())
	})

	t.Run("non-zero values", func(t *testing.T) {
		p := NewPair(1.0, 2.0)
		assert.Equal(t, 1.0, p.X())
		assert.Equal(t, 2.0, p.Y())
	})
}

func TestPair_WithValues(t *testing.T) {
	t.Run("numeric values", func(t *testing.T) {
		a := NewValue(2.0)
		b := NewValue(3.0)
		p := NewPair[*Value[float64]](a, b)

		assert.Equal(t, 2.0, p.X().GetValue())
		assert.Equal(t, 3.0, p.Y().GetValue())
	})

	t.Run("with operations", func(t *testing.T) {
		a := NewValue(2.0)
		b := NewValue(3.0)
		c := a.Add(b).(*Value[float64])
		p := NewPair[*Value[float64]](a, c)

		assert.Equal(t, 2.0, p.X().GetValue())
		assert.Equal(t, 5.0, p.Y().GetValue())
	})
}

func TestPair_At(t *testing.T) {
	p := NewPair(1.0, 2.0)

	t.Run("valid indices", func(t *testing.T) {
		assert.Equal(t, 1.0, p.At(0))
		assert.Equal(t, 2.0, p.At(1))
	})

	t.Run("panic on invalid index", func(t *testing.T) {
		assert.Panics(t, func() {
			p.At(2)
		})
	})
}

func TestPair_Accessors(t *testing.T) {
	t.Run("first and second", func(t *testing.T) {
		p := NewPair(1.0, 2.0)
		assert.Equal(t, 1.0, p.first)
		assert.Equal(t, 2.0, p.second)
	})

	t.Run("X and Y", func(t *testing.T) {
		p := NewPair(1.0, 2.0)
		assert.Equal(t, 1.0, p.X())
		assert.Equal(t, 2.0, p.Y())
	})
}
