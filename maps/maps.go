package main

import "fmt"

type Vertex struct {
	Lat, Long float64
}

var m = map[string]Vertex{
	"Bell Labs": {
		40.68433, -74.39967,
	},
	"Google": {
		37.42202, -122.08408,
	},
}

func main() {

	fmt.Println(m)
	mapsx()
}

func mapsx() {
	x := make(map[string]int)
	x["Answer"] = 42
	fmt.Println("The value:", x["Answer"])

	x["Answer"] = 48
	fmt.Println("The value:", x["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
}
