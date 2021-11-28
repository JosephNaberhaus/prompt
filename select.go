package prompt

import (
	"fmt"
	"github.com/rivo/uniseg"
	"strings"
)

const defaultNumLinesShown = 7

type SelectionOption struct {
	Name        string
	Description string
}

type line struct {
	optionIndex int
	text string
	isFirst bool
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
		oldOptionIndex := s.lines[s.cursor].optionIndex
		for oldOptionIndex == s.lines[s.cursor].optionIndex  ||	 !s.lines[s.cursor].isFirst {
			s.cursor--
			if s.cursor < 0 {
				s.cursor = len(s.lines) - 1
			}
		}
	} else if input == ControlDown {
		oldOptionIndex := s.lines[s.cursor].optionIndex
		for oldOptionIndex == s.lines[s.cursor].optionIndex {
			s.cursor++
			if s.cursor == len(s.lines) {
				s.cursor = 0
			}

			if s.offset > 0 {
				s.offset--
			}
		}
	} else if input == ControlEnter {
		s.output.showCursor()
		s.render(true)
		s.finish()
		return
	}

	s.render(false)
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
		s.output.writeColor("(Use arrow keys)", colorGreen)
	}
	s.output.nextLine()

	cursorLine := s.lines[s.cursor]

	startOffset := (-s.NumLinesToShow() / 2) + s.offset
	endOffset := (s.NumLinesToShow() / 2) + s.offset

	for offset := startOffset; offset <= endOffset; offset++ {
		lineIndex := s.cursor + offset
		if lineIndex < 0 {
			lineIndex += len(s.lines)
		} else if lineIndex >= len(s.lines) {
			lineIndex -= len(s.lines)
		}

		line := s.lines[lineIndex]
		if lineIndex == s.cursor {
			s.output.writeColor("> ", colorCyan)
		} else {
			s.output.write("  ")
		}

		if line.optionIndex == cursorLine.optionIndex {
			s.output.writeColorLn(line.text, colorCyan)
		} else {
			s.output.writeLn(line.text)
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
		wrappedDescription := wrapString(option.Description, s.output.outputWidth - longestName - 4)

		for i, wrapped := range wrappedDescription {
			var currentLineText string
			if i == 0 {
				padding := strings.Repeat(" ", longestName - uniseg.GraphemeClusterCount(option.Name))
				currentLineText = fmt.Sprintf("%s: %s%s", option.Name, padding, wrapped)
			} else {
				currentLineText = fmt.Sprintf("%s  %s", strings.Repeat(" ", longestName), wrapped)
			}

			s.lines = append(s.lines, line{
				optionIndex: optionIndex,
				text: currentLineText,
				isFirst: i == 0,
			})
		}
	}
}

func (s *Select) Response() SelectionOption {
	return s.Options[s.lines[s.cursor].optionIndex]
}

func (s *Select) longestName() int {
	longestName := 0
	for _, option := range s.Options {
		longestName = max(uniseg.GraphemeClusterCount(option.Name), longestName)
	}

	return longestName
}