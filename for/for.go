package main

import (
	"fmt"
	"math"
)

func only_condition() {
	sum := 1
	for sum < 10 {
		sum += sum
	}
	fmt.Println(sum)
}
func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	return lim
}

func Sqrt(x float64) float64 {
	z := 1.0
	for {
		prev := z
		// Newton's method formula
		z -= (z*z - x) / (2 * z)

		// Stop when the change becomes smaller than a tiny threshold
		if math.Abs(z-prev) < 1e-10 {
			break
		}
	}
	return z
}

func main() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
	only_condition()
	fmt.Println(sqrt(2), sqrt(-4))
	fmt.Println(pow(3, 2, 10), pow(3, 3, 20))

	x := 2.0
	myResult := Sqrt(x)
	stdResult := math.Sqrt(x)

	fmt.Printf("Custom Sqrt(%g): %v\n", x, myResult)
	fmt.Printf("math.Sqrt(%g):   %v\n", x, stdResult)
	fmt.Printf("Difference:       %g\n", math.Abs(myResult-stdResult))
	f := 27.0
	fmt.Printf("CubeRoot(%g) = %v\n", f, CubeRoot(f))
	fmt.Printf("math.Cbrt(%g) = %v\n", f, math.Cbrt(f))
	fmt.Printf("Difference:       %g\n", math.Abs(CubeRoot(f)-math.Cbrt(f)))

	start := 27
	steps, maxVal := Collatz(start)

	fmt.Printf("Starting Number: %d\n", start)
	fmt.Printf("Steps to reach 1: %d\n", steps)
	fmt.Printf("Highest Peak: %d\n", maxVal)
}

func CubeRoot(x float64) float64 {
	z := 1.0
	for {
		prev := z
		z -= (z*z*z - x) / (3 * z * z)
		if math.Abs(z-prev) < 1e-10 {
			break
		}
	}
	return z
}

func Collatz(n int) (steps int, maxVal int) {
	steps = 0
	maxVal = n

	for n > 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = n*3 + 1
		}
		if n > maxVal {
			maxVal = n
		}
		steps++
	}
	return steps, maxVal

}
