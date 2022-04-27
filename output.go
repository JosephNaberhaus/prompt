package prompt

import (
	"fmt"
	"github.com/rivo/uniseg"
	escapes "github.com/snugfox/ansi-escapes"
	"os"
	"strings"
	"unicode"
)

type output struct {
	outputWidth             int
	cursorColumn, cursorRow int

	numExtraLinesInBuffer int
	numExtraLinesWritten  int
	buffer                strings.Builder
}

func newOutput() (*output, error) {
	dimensions, err := escapes.GetConsoleSize(os.Stdout.Fd())
	if err != nil {
		return nil, fmt.Errorf("couldn't get console size: %v", err)
	}

	return &output{outputWidth: dimensions.Cols}, nil
}

func (o *output) write(content string) {
	gc := uniseg.NewGraphemes(content)
	for gc.Next() {
		current := gc.Str()
		isNonAscii := len(current) != 1 && current[0] > unicode.MaxASCII

		if isNonAscii {
			o.buffer.WriteString(escapes.CursorSavePosition)
		}

		if current == "\n" {
			o.nextLine()
			continue
		}

		o.buffer.WriteString(current)

		if isNonAscii {
			o.buffer.WriteString(escapes.CursorRestorePosition)
			o.buffer.WriteString(escapes.CursorForward)
		}

		o.cursorColumn++
		o.wrapCursor()
	}
}

func (o *output) writeColor(content string, c color) {
	o.buffer.WriteString(c.toTextEscapes())
	o.write(content)
	o.buffer.WriteString(colorWhite.toTextEscapes())
}

func (o *output) writeLn(content string) {
	o.write(content)
	o.nextLine()
}

func (o *output) writeColorLn(content string, color color) {
	o.writeColor(content, color)
	o.nextLine()
}

func (o *output) setCursor(row, col int) {
	if row > o.numExtraLinesWritten {
		o.nextLine()
	}

	o.moveCursor(row-o.cursorRow, col-o.cursorColumn)
}

func (o *output) moveCursor(numRows, numCols int) {
	if numRows == 0 && numCols == 0 {
		return
	}

	o.cursorColumn += numCols
	o.cursorRow += numRows

	o.buffer.WriteString(escapes.CursorMove(numCols, numRows))
}

func (o *output) nextLine() {
	if o.cursorColumn == o.outputWidth {
		return
	}

	if o.cursorRow+1 < o.numExtraLinesWritten {
		o.setCursor(o.cursorRow+1, 0)
		return
	}

	o.buffer.WriteString("\n ")
	o.cursorRow++
	o.cursorColumn = 1

	o.moveCursor(0, -1)

	o.numExtraLinesWritten = max(o.cursorRow, o.numExtraLinesWritten)
	o.numExtraLinesInBuffer = max(o.cursorRow, o.numExtraLinesInBuffer)
}

func (o *output) wrapCursor() {
	deltaRow := (o.cursorColumn - 1) / o.outputWidth
	o.cursorRow += deltaRow
	o.cursorColumn -= deltaRow * o.outputWidth

	o.numExtraLinesWritten = max(o.cursorRow, o.numExtraLinesWritten)
	o.numExtraLinesInBuffer = max(o.cursorRow, o.numExtraLinesInBuffer)
}

func (o *output) hideCursor() {
	o.buffer.WriteString(escapes.CursorHide)
}

func (o *output) showCursor() {
	o.buffer.WriteString(escapes.CursorShow)
}

func (o *output) clear() {
	for row := 0; row <= o.numExtraLinesInBuffer; row++ {
		o.setCursor(row, 0)
		o.buffer.WriteString(escapes.EraseLine)
	}

	o.setCursor(0, 0)
	o.numExtraLinesInBuffer = 0
}

func (o *output) flush() {
	fmt.Print(o.buffer.String())
	o.buffer.Reset()
}

func (o *output) commit() {
	o.setCursor(o.numExtraLinesInBuffer, 0)
	o.nextLine()
	o.flush()
}

func (o *output) uncommit() {
	o.clear()
	o.numExtraLinesInBuffer = 0
	o.numExtraLinesWritten = 0
	o.flush()
}
