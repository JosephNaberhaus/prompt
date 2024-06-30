package prompt

import (
	"fmt"
	"github.com/rivo/uniseg"
	"strings"
)

const defaultNumLinesShown = 7

type SelectionOption struct {
	ID          string
	Name        string
	Description string
}

type line struct {
	optionIndex int
	text        string
	isFirst     bool
}

type Select struct {
	base

	// The question to display to the user
	Question string

	// Array of options for the user to select
	Options []SelectionOption

	// The number of lines that will be shown at a time
	// Default is 7
	NumLinesShown int

	// Called when a key is pressed but before it is processed. Return `false` to cancel the event.
	OnKeyFunc func(Prompt, Key) bool

	offset int
	cursor int
	filter string

	lines []line
}

func (s *Select) Show() error {
	err := s.show()
	if err != nil {
		return err
	}

	s.computeLines()
	s.offset = s.NumLinesToShow() / 2

	s.output.hideCursor()
	s.render(false)

	for s.State() == Showing {
		nextKey, err := s.nextKey()
		if err != nil {
			s.output.showCursor()
			s.finish()
			return err
		}

		s.handleInput(nextKey)
	}

	return nil
}

func (s *Select) handleInput(input Key) {
	if s.OnKeyFunc != nil && !s.OnKeyFunc(s, input) {
		return
	}

	if input == ControlUp {
		if len(s.filteredOptions()) > 1 {
			for {
				s.cursor--
				if s.cursor < 0 {
					s.cursor = len(s.lines) - 1
				}

				curLine := s.lines[s.cursor]

				// If we haven't reached the first line then keep going.
				if !curLine.isFirst {
					continue
				}

				// If we're on a filtered out line then keep going
				if !s.matchesFilter(s.Options[curLine.optionIndex]) {
					continue
				}

				// Otherwise, we've completed our search.
				break
			}
		}
	} else if input == ControlDown {
		if len(s.filteredOptions()) > 1 {
			for {
				s.cursor++
				if s.cursor == len(s.lines) {
					s.cursor = 0
				}

				curLine := s.lines[s.cursor]

				// If we're on a filtered out line then keep going
				if !s.matchesFilter(s.Options[curLine.optionIndex]) {
					continue
				}

				// Otherwise, we've completed our search.
				break
			}
		}
	} else if input == ControlEnter {
		if len(s.filteredOptions()) != 0 {
			s.output.showCursor()
			s.render(true)
			s.finish()
			return
		}
	} else if input.IsText() {
		s.filter += string(input.Rune())

		if len(s.filteredOptions()) > 0 && !s.matchesFilter(s.curOption()) {
			closestValidLine := -1
			for i, line := range s.lines {
				// We're only interested in the first lines.
				if !line.isFirst {
					continue
				}

				// This line doesn't match the filter
				if !s.matchesFilter(s.Options[line.optionIndex]) {
					continue
				}

				if closestValidLine == -1 {
					closestValidLine = i
				} else if abs(s.cursor-i) < abs(s.cursor-closestValidLine) {
					closestValidLine = i
				}
			}

			s.cursor = closestValidLine
		}
	} else if input == ControlBackspace {
		if s.filter != "" {
			s.filter = s.filter[:len(s.filter)-1]
		}
	}

	if s.State() != Waiting {
		s.render(false)
	}
}

func (s *Select) curOption() SelectionOption {
	return s.Options[s.lines[s.cursor].optionIndex]
}

func (s *Select) matchesFilter(option SelectionOption) bool {
	if s.filter == "" {
		return true
	}

	if strings.Contains(strings.ToLower(option.Name), strings.ToLower(s.filter)) {
		return true
	}

	return false
}

func (s *Select) filteredOptions() []SelectionOption {
	var result []SelectionOption
	for _, option := range s.Options {
		if s.matchesFilter(option) {
			result = append(result, option)
		}
	}

	return result
}

func (s *Select) NumLinesToShow() int {
	if s.NumLinesShown <= 0 {
		return min(defaultNumLinesShown, len(s.lines))
	}

	return min(s.NumLinesShown, len(s.lines))
}

func (s *Select) render(isFinished bool) {
	s.output.clear()

	s.output.writeColor("? ", colorGreen)
	s.output.write(s.Question)
	s.output.write(": ")
	if isFinished {
		s.output.writeColor(fmt.Sprintf("%s: %s", s.Response().Name, s.Response().Description), colorCyan)
		return
	} else {
		s.output.writeColor("(Use arrow keys) (Type to filter)", colorGreen)
	}
	s.output.nextLine()

	cursorLine := s.lines[s.cursor]

	startOffset := (-s.NumLinesToShow() / 2) + s.offset
	endOffset := (s.NumLinesToShow() / 2) + s.offset

	fillRemainingWithBlank := false
	if len(s.filteredOptions()) == 0 {
		s.output.writeColorLn(s.filter, colorRed)
		fillRemainingWithBlank = true
		startOffset++
	}

	for offset := startOffset; offset <= endOffset; offset++ {
		lineIndex := s.actualLineNumber(s.cursor + offset)

		// We've looped back to the start
		if offset != startOffset && lineIndex == s.actualLineNumber(s.cursor+startOffset) {
			fillRemainingWithBlank = true
		}
		if fillRemainingWithBlank {
			s.output.nextLine()
			continue
		}

		line := s.lines[lineIndex]
		option := s.Options[line.optionIndex]

		if !s.matchesFilter(option) {
			endOffset++
			continue
		}

		if lineIndex == s.cursor {
			s.output.writeColor("> ", colorCyan)
		} else {
			s.output.write("  ")
		}

		if line.isFirst {
			option := s.Options[line.optionIndex]

			redRemaining := 0
			for i, c := range line.text {
				if i >= len(option.Name) {
					if line.optionIndex == cursorLine.optionIndex {
						s.output.writeColor(string(c), colorCyan)
					} else {
						s.output.write(string(c))
					}
				} else {
					if s.filter != "" && strings.HasPrefix(strings.ToLower(option.Name[i:]), strings.ToLower(s.filter)) {
						s.output.writeColor(string(c), colorRed)
						redRemaining = len(s.filter) - 1
					} else if redRemaining > 0 {
						s.output.writeColor(string(c), colorRed)
						redRemaining--
					} else if line.optionIndex == cursorLine.optionIndex {
						s.output.writeColor(string(c), colorCyan)
					} else {
						s.output.write(string(c))
					}
				}
			}

			s.output.nextLine()
		} else {
			if line.optionIndex == cursorLine.optionIndex {
				s.output.writeColorLn(line.text, colorCyan)
			} else {
				s.output.writeLn(line.text)
			}
		}
	}

	if len(s.lines) > s.NumLinesToShow() {
		s.output.writeColor("(Move up and down to reveal more choices)", colorGreen)
	}

	s.output.flush()
}

func (s *Select) computeLines() {
	s.lines = make([]line, 0, len(s.Options))

	longestName := s.longestName()

	for optionIndex, option := range s.Options {
		wrappedDescription := wrapString(option.Description, s.output.outputWidth-longestName-4)

		for i, wrapped := range wrappedDescription {
			var currentLineText string
			if i == 0 {
				padding := strings.Repeat(" ", longestName-uniseg.GraphemeClusterCount(option.Name))
				currentLineText = fmt.Sprintf("%s: %s%s", option.Name, padding, wrapped)
			} else {
				currentLineText = fmt.Sprintf("%s  %s", strings.Repeat(" ", longestName), wrapped)
			}

			s.lines = append(s.lines, line{
				optionIndex: optionIndex,
				text:        currentLineText,
				isFirst:     i == 0,
			})
		}
	}
}

func (s *Select) actualLineNumber(line int) int {
	if line < 0 {
		return line + len(s.lines)
	} else if line >= len(s.lines) {
		return line - len(s.lines)
	}

	return line
}

func (s *Select) Response() SelectionOption {
	return s.curOption()
}

func (s *Select) longestName() int {
	longestName := 0
	for _, option := range s.Options {
		longestName = max(uniseg.GraphemeClusterCount(option.Name), longestName)
	}

	return longestName
}
