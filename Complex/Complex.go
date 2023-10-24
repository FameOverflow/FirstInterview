package main

import (
	"fmt"
	"math"
)

type Complex struct {
	Real float64
	Imag float64
}

type typer interface {
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
	if c.Imag == 0 {
		return fmt.Sprintf("%g", c.Real)
	}
	if c.Imag > 0 {
		return fmt.Sprintf("%g+%gi", c.Real, c.Imag)
	}
	return fmt.Sprintf("%g%gi", c.Real, c.Imag)
}

func main() {
	c1 := Complex{1, 2}
	fmt.Println("c1: " + c1.ToString())
	c2 := Complex{3, 4}
	fmt.Println("c2: " + c2.ToString())
	c3 := c1.Add(c2)
	fmt.Println("c1+c2: " + c3.ToString())
	c3 = c1.Sub(c2)
	fmt.Println("c1-c2: " + c3.ToString())
	c3 = c1.Mul(c2)
	fmt.Println("c1*c2: " + c3.ToString())
	c3 = c1.Div(c2)
	fmt.Println("c1/c2: " + c3.ToString())
	fmt.Println("c1的模: " + fmt.Sprintf("%g", c1.Mod()))
}
