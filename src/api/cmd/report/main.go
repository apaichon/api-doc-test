package main

import (
	"fmt"
	"os"

	"api/test"
)

func main() {
	if err := test.GenerateReport(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}
} 