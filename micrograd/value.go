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
	constraints.Integer | constraints.Float | constraints.Complex
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
}

type Value[K BaseNumeric] struct {
	Name string

	datum     K
	gradient  K
	children  Pair[Numeric[K]]
	operation OperationEnum
}

var _ Numeric[int] = NewValue(0)

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
