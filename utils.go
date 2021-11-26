package prompt

import (
	"github.com/rivo/uniseg"
	"strings"
)

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func wrapString(toWrap string, width int) []string {
	if len(toWrap) == 0 {
		return []string{""}
	}

	wrapped := make([]string, 0, (len(toWrap)/width)+1)
	gc := uniseg.NewGraphemes(toWrap)

	current := strings.Builder{}
	currentLength := 0
	for gc.Next() {
		current.WriteString(gc.Str())
		currentLength++

		if currentLength == width {
			wrapped = append(wrapped, current.String())
			current.Reset()
			currentLength = 0
		}
	}

	if currentLength > 0 {
		wrapped = append(wrapped, current.String())
	}

	return wrapped
}
