package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	// 1. Allocate outer slice of length dy (rows)
	picture := make([][]uint8, dy)

	// 2. Loop through each row index y
	for y := 0; y < dy; y++ {
		// Allocate inner slice of length dx (columns)
		picture[y] = make([]uint8, dx)

		// 3. Loop through each column index x
		for x := 0; x < dx; x++ {
			// Compute pixel value using a formula and convert to uint8
			picture[y][x] = uint8((x + y) / 2)
		}
	}

	// 4. Return the completed 2D slice
	return picture
}

func main() {
	pic.Show(Pic)
}
