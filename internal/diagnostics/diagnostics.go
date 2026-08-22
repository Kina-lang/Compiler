package diagnostics

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// File is one source file + its line index
// Spans are byte offsets, line/column is derived on demand,
// not stored per token.
type File struct {
	Name  string
	src   []byte
	lines []int // byte offset of each line start
}

func NewFile(name string, src []byte) *File {
	lines := []int{0}

	// Loop over every byte in the source
	for i, b := range src {
		// If the byte is newline, append offset of the next byte into
		// the line starts array
		// TODO: Add support for \r (\r\n) should be supported by this
		if b == '\n' {
			lines = append(lines, i+1)
		}
	}

	// Create new File struct and return pointer
	return &File{name, src, lines}
}

// Converts a byte offset to 1-based line and column
func (file *File) LineCol(offset int) (line, col int) {
	// Finds the line that contains the byte offset
	// Finds lowest line number that has offset larger than the wanted
	// offset and then subtracts 1 (=> get the line that includes the offset)
	i := sort.Search(len(file.lines), func(i int) bool {
		return file.lines[i] > offset
	}) - 1

	// Return the line (1-based offset) and column (offset - line offset, 1-based)
	return i + 1, offset - file.lines[i] + 1
}

// Gets text on the specified line
func (file *File) LineText(line int) string {
	start := file.lines[line-1] // start byte offset
	end := len(file.src)        // last available byte offset

	// If not last line
	if line < len(file.lines) {
		// Byte offset of next line - 1 => end byte offset of wanted line
		end = file.lines[line] - 1
	}

	// Return string slice from the source code
	return string(file.src[start:end])
}

type Diagnostic struct {
	File       *File
	Code       string // "Exxxx"
	Start, End int
	Message    string
}

// Bag collects diagnostics for every module in the program
type Bag struct {
	list  []Diagnostic
	Color bool
}

// Reporter is a Bag bound to a specific file
type Reporter struct {
	bag  *Bag
	file *File
}

// Creates a new diagnostics Bag
func NewBag(color bool) *Bag {
	return &Bag{
		Color: color,
	}
}

// Get reporter for the specified file
func (bag *Bag) For(file *File) *Reporter {
	return &Reporter{bag, file}
}

// Appends error to the Bag, has var args support
func (reporter *Reporter) Errorf(start, end int, code, format string, args ...any) {
	reporter.bag.list = append(reporter.bag.list, Diagnostic{reporter.file, code, start, end, fmt.Sprintf(format, args...)})
}

// Checks whenever the bag has any recorded errors
func (bag *Bag) HasErrors() bool {
	return len(bag.list) > 0
}

// Prints error if any errors got reported
func (bag *Bag) Err(w io.Writer) error {
	if len(bag.list) == 0 {
		return nil
	}

	bag.Render(w)

	return fmt.Errorf("%d error(s)", len(bag.list))
}

// Render errors to the given writer
func (bag *Bag) Render(w io.Writer) {
	// Sorts the errors by byte offsets and file names (alphabetical)
	// "stable" keeps the equal elements in the same position (order of being reported)
	slices.SortStableFunc(bag.list, func(a, b Diagnostic) int {
		if a.File != b.File {
			return strings.Compare(a.File.Name, b.File.Name)
		}

		return a.Start - b.Start
	})

	for _, diagnostic := range bag.list {
		var line, col int
		var positionDefined bool = diagnostic.Start >= 0 && diagnostic.End >= 0

		if positionDefined {
			line, col = diagnostic.File.LineCol(diagnostic.Start) // Get line and column of the diagnostic (the one where the diagnostic starts)
		}

		var posString string
		if positionDefined {
			posString = fmt.Sprintf("%s:%d:%d", diagnostic.File.Name, line, col)
		} else {
			posString = diagnostic.File.Name
		}

		// Print file:line:col (clickable in IDEs), code and error message
		if bag.Color {
			fmt.Fprintf(w, "%s %s%s %s\n", color.BlackString(posString), color.HiRedString(diagnostic.Code), color.BlackString(":"), diagnostic.Message)
		} else {
			fmt.Fprintf(w, "%s %s: %s\n", posString, diagnostic.Code, diagnostic.Message)
		}

		if !positionDefined {
			continue
		}

		// Print source code line and pointer to the columns affected
		lineText := diagnostic.File.LineText(line)

		fmt.Fprintf(w, "  %s\n", lineText)
		fmt.Fprintf(w, "  %s%s\n", strings.Repeat(" ", col-1), strings.Repeat("^", max(1, min(diagnostic.End-diagnostic.Start, len(lineText)))))
	}
}
