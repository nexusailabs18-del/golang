package main

import (
	"fmt"
	"strings"
)

func main() {
	i, j := 42, 2701
	p := &i
	fmt.Println(*p)

	*p = 21
	fmt.Println(i)

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j

	mainy()
	fmt.Println(v1, pv, v2, v3)

	array()
	names_slices()
	slice_literals()
	slice_length()
	slice_with_make()
	tic_tac_toe_board()
	append_a_slice()
	sqaure_root_of_2()
}

type Vertex struct {
	X int
	Y int
}

var (
	v1 = Vertex{1, 2}
	v2 = Vertex{X: 1}
	v3 = Vertex{}
	pv = &Vertex{1, 2}
)

func mainy() {
	fmt.Println(Vertex{1, 2})

	v := Vertex{1, 2}
	p := &v
	p.X = 1e9
	fmt.Println(v)
}

func array() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	var s []int = primes[1:4]
	fmt.Println(s)
}

func names_slices() {
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)
	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)

	b[0] = "XXX"
	fmt.Println(a, b)
	fmt.Println(names)
}

func slice_literals() {
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)
	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	sap := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(sap)

}
func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
func slice_length() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	s = s[:0]
	printSlice(s)
	s = s[:4]
	printSlice(s)

	s = s[2:]
	printSlice(s)
}

func slice_with_make() {
	a := make([]int, 5)
	printSlicey("a", a)
	b := make([]int, 0, 5)
	printSlicey("b", b)

	c := b[:2]
	printSlicey("c", c)

	d := c[2:5]
	printSlicey("d", d)
}
func printSlicey(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

func tic_tac_toe_board() {
	board := [][]string{
		{"_", "_", "_"},
		{"_", "_", "_"},
		{"_", "_", "_"},
	}

	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}
func append_a_slice() {
	var s []int
	printSlice(s)

	s = append(s, 0)
	printSlice(s)
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 2, 3, 4)
	printSlice(s)
}

var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func sqaure_root_of_2() {
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}
}
