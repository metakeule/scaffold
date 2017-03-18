package main

import (
	"fmt"
	"github.com/metakeule/scaffold"
)

func main() {
	templ, err := scaffold.Scan("models")

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	fmt.Println(string(templ))
}
