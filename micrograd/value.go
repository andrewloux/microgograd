package micrograd

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Pair[K any] struct {
	data [2]K
}

func (p Pair[K]) String() string {
	return fmt.Sprintf("Pair(%v, %v)", p.data[0], p.data[1])
}

func (p Pair[K]) Data() [2]K {
	return p.data
}

func NewPair[K any](a, b K) Pair[K] {
	return Pair[K]{
		data: [2]K{a, b},
	}
}

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
	Add(Numeric[K]) Numeric[K]
	Sub(Numeric[K]) Numeric[K]
	Mul(Numeric[K]) Numeric[K]
	GetValue() K
	GetChildren() Pair[Numeric[K]]
	GetOperation() OperationEnum
}

type Value[K BaseNumeric] struct {
	Name string

	datum     K
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

func (v *Value[K]) GetValue() K {
	return v.datum
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

type ValueOptions struct {
	Name string
}

type ValueOpt func(*ValueOptions)

func WithName(name string) ValueOpt {
	return func(cur *ValueOptions) {
		cur.Name = name
	}
}

func NewValue[K BaseNumeric](input K, options ...ValueOpt) *Value[K] {
	opts := &ValueOptions{}
	for _, o := range options {
		o(opts)
	}

	return &Value[K]{
		datum: input,
		Name:  opts.Name,
	}
}
