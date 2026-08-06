package main

import (
	"fmt"
	"math/cmplx"
)

var (
	c, python int  = 1, 2
	java      bool = true
)

func AddUp(x int, y int) int {
	return x + y
}

func Split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

func spl() {
	fmt.Println(Split(17))
}

func main() {
	fmt.Println(AddUp(42, 49))
	spl()

	var rust, css, java = true, false, "no!"
	k := 3
	fmt.Println(rust, css, c, python, java, k)

	fmt.Printf("Type : %T Value : %v \n ", ToBe, ToBe)
	fmt.Printf("Type : %T Value : %v \n ", MaxInt, MaxInt)
	fmt.Printf("Type : %T Value : %v \n ", z, z)
}
