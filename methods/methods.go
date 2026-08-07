package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	x := Vertex{3, 4}
	fmt.Println(x.Abs())
	y := Vertex{5, 12}
	fmt.Println(Abs(y))
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Calc())
	x.Scale(10)
	fmt.Println(x.Abs())
	v := Vertex{3, 4}
	Scalex(&v, 10)
	fmt.Println(Abs(v))
	just_see()

}
func just_see() {
	v := Vertex{3, 4}
	fmt.Printf("Before scaling: %+v, Abs: %v\n", v, v.Abs())
	v.Scale(5)
	fmt.Printf("After scaling: %+v, Abs: %v\n", v, v.Abs())
}
func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type MyFloat float64

func (f MyFloat) Calc() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func Scalex(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}
