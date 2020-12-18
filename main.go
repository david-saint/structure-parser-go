package main

import (
	"calhoun/appscript/parser"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("in.json")
	defer file.Close()

	out, err := parser.Parse(file)

	if err != nil {
		fmt.Print("Error", err)
		os.Exit(1)
	}

	fmt.Print(out)
}
