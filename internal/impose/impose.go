package impose

import (
	"fmt"
	"math"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func CalculatePageOrder(totalPages int) []int {
	paddedPages := int(math.Ceil(float64(totalPages)/4) * 4)
	sheets := paddedPages / 4
	pageOrder := make([]int, 0, paddedPages)

	for sheet := 0; sheet < sheets; sheet++ {
		// Front of sheet
		pageOrder = append(pageOrder, paddedPages-1-2*sheet) // Last page of current set
		pageOrder = append(pageOrder, 2*sheet)               // First page of current set
		// Back of sheet (when paper is flipped)
		pageOrder = append(pageOrder, 2*sheet+1)             // Second page of current set
		pageOrder = append(pageOrder, paddedPages-2-2*sheet) // Second-to-last page of current set
	}

	return pageOrder
}

func PDF(inputPath, outputPath string) error {
	// Get page count
	ctx, err := api.ReadContextFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading PDF: %v", err)
	}

	pages := ctx.PageCount
	pageOrder := CalculatePageOrder(pages)

	// Convert page numbers to selection string
	var pageList []string
	for _, p := range pageOrder {
		if p < pages {
			pageList = append(pageList, fmt.Sprintf("%d", p+1))
		}
	}

	// Open input file
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("opening input file: %v", err)
	}
	defer inFile.Close()

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outFile.Close()

	// Collect pages in specified order
	if err := api.Collect(inFile, outFile, pageList, nil); err != nil {
		return fmt.Errorf("reordering pages: %v", err)
	}

	return nil
}
