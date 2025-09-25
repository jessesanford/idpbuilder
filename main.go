package main

import (
	"fmt"
	"os"
)

const version = "split-001-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("idpbuilder %s (Split 001 - API Types Only)\n", version)
		return
	}

	fmt.Println("idpbuilder Split 001: API Types and Core")
	fmt.Println("This split contains only type definitions.")
	fmt.Println("Use --version to see version info.")
}
