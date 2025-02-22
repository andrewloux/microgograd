package micrograd

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type OperationEnum int

const (
	UNSET = iota
	ADD   = '+'
	SUB   = '-'
	MUL   = '*'
)

type BaseNumeric interface {
	constraints.Float
}

type Numeric[K BaseNumeric] interface {
	GetName() string
	SetName(string) *Value[K]
	Add(Numeric[K]) Numeric[K]
	Sub(Numeric[K]) Numeric[K]
	Mul(Numeric[K]) Numeric[K]
	GetValue() K
	SetValue(K) *Value[K]
	GetGradient() K
	SetGradient(K) *Value[K]
	GetChildren() Pair[Numeric[K]]
	GetOperation() OperationEnum
	Backtrack()
}

type Value[K BaseNumeric] struct {
	Name string

	datum     K
	gradient  K
	children  Pair[Numeric[K]]
	operation OperationEnum
}

var _ Numeric[float64] = NewValue(0.0)

func (v *Value[K]) Add(input Numeric[K]) Numeric[K] {
	return &Value[K]{
		datum:     v.datum + input.GetValue(),
		children:  NewPair[Numeric[K]](v, input),
		operation: ADD,
	}
}

func (v *Value[K]) Sub(k Numeric[K]) Numeric[K] {
	return &Value[K]{
		datum:     v.datum - k.GetValue(),
		operation: SUB,
		children:  NewPair[Numeric[K]](v, k),
	}
}

func (v *Value[K]) Mul(k Numeric[K]) Numeric[K] {
	return &Value[K]{
		datum:     v.datum * k.GetValue(),
		operation: MUL,
		children:  NewPair[Numeric[K]](v, k),
	}
}

func (v *Value[K]) GetName() string {
	return v.Name
}

func (v *Value[K]) SetName(input string) *Value[K] {
	v.Name = input
	return v
}

func (v *Value[K]) GetValue() K {
	return v.datum
}

func (v *Value[K]) SetValue(input K) *Value[K] {
	v.datum = input
	return v
}

func (v *Value[K]) GetChildren() Pair[Numeric[K]] {
	return v.children
}

func (v *Value[K]) GetOperation() OperationEnum {
	return v.operation
}

func (v *Value[K]) String() string {
	return fmt.Sprintf("Value(data=%v)", v.GetValue())
}

func (v *Value[K]) GetGradient() K {
	return v.gradient
}

func (v *Value[K]) SetGradient(input K) *Value[K] {
	v.gradient = input
	return v
}

func (v *Value[K]) Backtrack() {
	// have: dO[utput]/dv
	// want: dO/da, dO/db -- i.e. we want to know how each leaf node (input)
	// 		 affects the overall output of the system
	a, b := v.GetChildren().first, v.GetChildren().second
	if a == nil && b == nil {
		return
	}

	switch v.GetOperation() {
	case ADD:
		// a + b
		// dv/da = 1
		// dv/da * dO/dv = dv/da
		// 1 *
		a.SetGradient(a.GetGradient() + 1*v.GetGradient())
		// dv/db = 1
		b.SetGradient(b.GetGradient() + 1*v.GetGradient())
	case MUL:
		// a * b
		// dv/da = b
		// dv/da * dO/dv = d0/da
		a.SetGradient(a.GetGradient() + b.GetValue()*v.GetGradient())
		// dv/db = a
		// dv/db * dO/dv = d0/db
		b.SetGradient(b.GetGradient() + a.GetValue()*v.GetGradient())
	}

	// continue with the rest of the backtrack logic
	a.Backtrack()
	b.Backtrack()
}

type ValueOptions[K BaseNumeric] struct {
	Name     string
	Value    K
	Gradient K
}

type ValueOpt[K BaseNumeric] func(*ValueOptions[K])

func WithName(name string) ValueOpt[float64] {
	return func(cur *ValueOptions[float64]) {
		cur.Name = name
	}
}

func WithGradient(input float64) ValueOpt[float64] {
	return func(cur *ValueOptions[float64]) {
		cur.Gradient = input
	}
}

func WithValue(input float64) ValueOpt[float64] {
	return func(cur *ValueOptions[float64]) {
		cur.Value = input
	}
}

func NewValue[K BaseNumeric](input K, options ...ValueOpt[K]) *Value[K] {
	opts := &ValueOptions[K]{}
	for _, o := range options {
		o(opts)
	}
	v := &Value[K]{
		datum: input,
	}
	if opts.Name != "" {
		v.SetName(opts.Name)
	}
	if opts.Gradient != 0 {
		v.SetGradient(opts.Gradient)
	}
	if opts.Value != 0 {
		v.SetValue(opts.Value)
	}
	return v
}
