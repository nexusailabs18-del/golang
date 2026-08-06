package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// 1. Slice holding three names
	names := []string{"Gladys", "Samantha", "Darrin"}

	// 2. Pass the slice to Hellos()
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Print the entire map output

	nameList := greetings.GetNames(messages)
	fmt.Println("Extracted Names:", nameList)
}
