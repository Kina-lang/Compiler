// NOTE: This file does not inline with our error code spec,
// 		 all error codes are made up, we don't need to adhere to
// 		 spec here.

package diagnostics

import (
	"bytes"
	"testing"
)

// offsets: l1 0..17, \n 18 | l2 19..36, \n 37 | l3 38, \n 39
var src = []byte("func main(): int {\n  val a: int = 3.1\n}\n")

func TestLineCol(t *testing.T) {
	file := NewFile("main.kin", src)
	tests := []struct {
		offset int
		line, col int
	}{
		{0, 1, 1}, // first byte
		{17, 1, 18}, // {, last byte of line 1
		{19, 2, 1}, // first byte after a newline
		{34, 2, 16}, // the 3 of 3.1
		{38, 3, 1}, // }
	}

	for _, tt := range tests {
		line, col := file.LineCol(tt.offset)

		if line != tt.line || col != tt.col {
			t.Errorf("LineCol(%d) = %d:%d, want %d:%d", tt.offset, line, col, tt.line, tt.col)
		}
	}
}

func TestLineText(t *testing.T) {
	file := NewFile("main.kin", src)
	tests := []struct {
		line int
		want string
	}{
		{1, "func main(): int {"},
		{2, "  val a: int = 3.1"},
		{3, "}"},
	}

	for _, tt := range tests {
		if got := file.LineText(tt.line); got != tt.want {
			t.Errorf("LineText(%d) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestNoTrailingNewline(t *testing.T) {
	file := NewFile("x.kin", []byte("abc"))
	if line, col := file.LineCol(2); line != 1 || col != 3 {
		t.Errorf("LineCol(2) = %d:%d, want 1:3", line, col)
	}

	if got := file.LineText(1); got != "abc" {
		t.Errorf("LineText(1) = %q, want %q", got, "abc")
	}
}

func TestRenderSortsByPosition(t *testing.T) {
	file := NewFile("main.kin", src)
	bag := &Bag{}
	reporter := bag.For(file)

	// Emitted out of order on purpose
	reporter.Errorf(38, 39, "E0123", "function 'main' must return a value on all paths")
	reporter.Errorf(34, 37, "E0234", "insane programmer wrote this, refusing to compile!")

	var buf bytes.Buffer
	bag.Render(&buf)

	want := "main.kin:2:16: E0234: insane programmer wrote this, refusing to compile!\n" +
			"    val a: int = 3.1\n" +
			"                 ^^^\n" +
			"main.kin:3:1: E0123: function 'main' must return a value on all paths\n" +
            "  }\n" +
            "  ^\n"

    if buf.String() != want {
    	t.Errorf("Render() =\n%s\nwant\n%s", buf.String(), want)
    }
}

// Test synthetic token (inserted semicolon - ASI) with zero-width span - must still produce
// visible caret
func TestZeroWidthSpan(t *testing.T) {
	file := NewFile("main.kin", src)
	bag := &Bag{}
	bag.For(file).Errorf(18, 18, "E0345", "expected ';'")

	var buf bytes.Buffer
	bag.Render(&buf)

	want := "main.kin:1:19: E0345: expected ';'\n" +
			"  func main(): int {\n"+
			"                    ^\n"

	if buf.String() != want {
		t.Errorf("Render() =\n%s\nwant\n%s", buf.String(), want)
	}
}
