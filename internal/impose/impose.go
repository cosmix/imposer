package impose

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// MinimumPages is the minimum number of pages required for proper booklet imposition
const MinimumPages = 8

// addBlankPages adds blank pages to a PDF file until it reaches the minimum page count
func addBlankPages(inputPath, outputPath string, targetPageCount int) error {
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

	// Get current page count
	pageCount, err := api.PageCount(inFile, nil)
	if err != nil {
		return fmt.Errorf("getting page count: %v", err)
	}

	pagesToAdd := targetPageCount - pageCount
	if pagesToAdd <= 0 {
		// No need to add pages, just copy the file
		if _, err := inFile.Seek(0, 0); err != nil {
			return fmt.Errorf("seeking input file: %v", err)
		}
		if _, err := io.Copy(outFile, inFile); err != nil {
			return fmt.Errorf("copying file: %v", err)
		}
		return nil
	}

	// Instead of trying to add multiple pages in one go, add them one by one
	// Create a temporary copy first
	tempDir := filepath.Dir(outputPath)
	tempFile, err := os.CreateTemp(tempDir, "temp-*.pdf")
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	// Copy the original content to the temp file
	if _, err := inFile.Seek(0, 0); err != nil {
		return fmt.Errorf("seeking input file: %v", err)
	}
	tempOut, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}
	if _, err := io.Copy(tempOut, inFile); err != nil {
		return fmt.Errorf("copying to temp file: %v", err)
	}
	tempOut.Close()

	// For each page we need to add, append a blank page
	for i := 0; i < pagesToAdd; i++ {
		// Create a temporary source file
		srcFile, err := os.Open(tempPath)
		if err != nil {
			return fmt.Errorf("opening temp file: %v", err)
		}

		// Create a temporary destination file
		destFile, err := os.CreateTemp(tempDir, "dest-*.pdf")
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("creating dest file: %v", err)
		}
		destPath := destFile.Name()
		destFile.Close()

		// Open destination for writing
		destOut, err := os.Create(destPath)
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("opening dest file: %v", err)
		}

		// Add a blank page after the last page
		selectedPage := fmt.Sprintf("%d", pageCount+i)
		if err := api.InsertPages(srcFile, destOut, []string{selectedPage}, false, nil, nil); err != nil {
			srcFile.Close()
			destOut.Close()
			os.Remove(destPath)
			return fmt.Errorf("inserting blank page %d: %v", i+1, err)
		}

		// Close files
		srcFile.Close()
		destOut.Close()

		// Replace temp file with dest file
		os.Remove(tempPath)
		if err := os.Rename(destPath, tempPath); err != nil {
			return fmt.Errorf("replacing temp file: %v", err)
		}
	}

	// Copy the final result to the output file
	finalSrc, err := os.Open(tempPath)
	if err != nil {
		return fmt.Errorf("opening final temp file: %v", err)
	}
	defer finalSrc.Close()

	if _, err := io.Copy(outFile, finalSrc); err != nil {
		return fmt.Errorf("copying to output file: %v", err)
	}

	return nil
}

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

	// If we have fewer than MinimumPages pages, we need to add blank pages
	if pages < MinimumPages {
		// Create a temporary file for the padded PDF
		tempDir := filepath.Dir(outputPath)
		tempFile, err := os.CreateTemp(tempDir, "padded-*.pdf")
		if err != nil {
			return fmt.Errorf("creating temp file: %v", err)
		}
		tempPath := tempFile.Name()
		tempFile.Close()
		defer os.Remove(tempPath)

		// Add blank pages to reach MinimumPages
		if err := addBlankPages(inputPath, tempPath, MinimumPages); err != nil {
			return fmt.Errorf("adding blank pages: %v", err)
		}

		// Use the padded file for reordering
		inputPath = tempPath
		pages = MinimumPages
	}

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
