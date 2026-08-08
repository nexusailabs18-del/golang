package main

import (
	"fmt"
	"math"
)

// Interfaces
type Abser interface {
	Abs() float64
}

type I interface {
	M()
}

// Types and Methods for Abser
type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Types and Methods for I
type T struct {
	S string
}

func (t *T) M() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func main() {
	// Part 1: Abser interface demo
	var a Abser

	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f
	fmt.Println(a.Abs()) // Output: 1.4142135623730951

	a = &v
	fmt.Println(a.Abs()) // Output: 5

	fmt.Println("---")

	// Part 2: I interface demo
	var i I = &T{S: "Hello"} // Added "Hello" to make output visible
	i.M()                    // Output: Hello
	describe(i)              // Output: (&{Hello}, *main.T)

	i = F(math.Pi)
	describe(i) // Output: (3.141592653589793, main.F)
	i.M()       // Output: 3.141592653589793

	do(21)
	do("hello")
	do(true)

	g := Person{"Arthur Dent", 42}
	z := Person{"Zaphod Beeblebrox", 9001}
	fmt.Println(g, z)
}

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}
