package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	words := strings.Fields(s)

	for _, w := range words {
		counts[w]++
	}
	return counts
}

func main() {
	text1 := "go is fun and go is fast"
	fmt.Println("Input:", text1)
	fmt.Println("Output:", WordCount(text1))

	fmt.Println("---")

	// Test 2: Sentence with multiple spaces and repeated words
	text2 := "apple banana   apple orange  banana apple"
	fmt.Println("Input:", text2)
	fmt.Println("Output:", WordCount(text2))
}
