# imposer

A command-line tool for imposing PDF files into booklet format. It reorders pages so that when printed double-sided and folded, they create a properly ordered booklet.

## Installation

### Via go install
```bash
go install github.com/cosmix/imposer/cmd/impose@latest
```

### Build from source
```bash
git clone https://github.com/cosmix/imposer.git
cd imposer
go build ./cmd/impose
```

## Usage

```bash
# Using flags
impose -i input.pdf -o output.pdf

# Using positional arguments
impose input.pdf output.pdf
```

## How It Works

The tool performs the following operations:
1. Reads the input PDF file
2. Calculates the correct page order for booklet printing
3. Reorders pages so that when printed double-sided and folded:
   - Pages appear in the correct reading order
   - The total page count is padded to a multiple of 4 if needed
   - Pages are arranged in printer spreads (last-first, first-last pattern)

## Example

For an 8-page document, pages will be arranged as:
- Sheet 1 front: Page 8 | Page 1
- Sheet 1 back: Page 2 | Page 7
- Sheet 2 front: Page 6 | Page 3
- Sheet 2 back: Page 4 | Page 5

When printed double-sided, folded in half, and stacked, the pages will read in order from 1 to 8.

## License

MIT License - see [LICENSE](LICENSE) for details
