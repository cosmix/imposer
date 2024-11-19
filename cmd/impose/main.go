package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cosmix/imposer/internal/impose"
)

func main() {
	// Define flags
	inputFile := flag.String("i", "", "Input PDF file")
	outputFile := flag.String("o", "", "Output PDF file")
	flag.Parse()

	// Check for positional arguments if flags aren't used
	args := flag.Args()
	if *inputFile == "" && *outputFile == "" && len(args) == 2 {
		*inputFile = args[0]
		*outputFile = args[1]
	}

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Usage:")
		fmt.Println("  impose -i input.pdf -o output.pdf")
		fmt.Println("  impose input.pdf output.pdf")
		os.Exit(1)
	}

	// Process the PDF
	if err := impose.PDF(*inputFile, *outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created imposed PDF: %s\n", *outputFile)
	fmt.Printf("Happy printing!\n")
}
