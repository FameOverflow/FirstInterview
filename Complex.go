package main

import (
	"fmt"
	"math"
)

type Complex struct {
	Real float64
	Imag float64
}

type Complexer interface {
	Add(c Complex) Complex
	Sub(c Complex) Complex
	Mul(c Complex) Complex
	Div(c Complex) Complex
	Mod(c Complex) Complex
	ToString(c Complex) string
}

func (c1 Complex) Add(c2 Complex) Complex {
	return Complex{c1.Real + c2.Real, c1.Imag + c2.Imag}
}

func (c1 Complex) Sub(c2 Complex) Complex {
	return Complex{c1.Real - c2.Real, c1.Imag - c2.Imag}
}

func (c1 Complex) Mul(c2 Complex) Complex {
	return Complex{c1.Real*c2.Real - c1.Imag*c2.Imag, c1.Real*c2.Imag + c1.Imag*c2.Real}
}

func (c1 Complex) Div(c2 Complex) Complex {
	return Complex{(c1.Real*c2.Real + c1.Imag*c2.Imag) / (c2.Real*c2.Real + c2.Imag*c2.Imag), (c1.Imag*c2.Real - c1.Real*c2.Imag) / (c2.Real*c2.Real + c2.Imag*c2.Imag)}
}

func (c Complex) Mod() float64 {
	return math.Sqrt(c.Real*c.Real + c.Imag*c.Imag)
}

func (c Complex) ToString() string {
	return fmt.Sprintf("%f + %fi", c.Real, c.Imag)
}
