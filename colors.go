package prompt

import escapes "github.com/snugfox/ansi-escapes"

type color uint8

const (
	colorBlack color = iota
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func (c color) toTextEscapes() string {
	switch c {
	case colorBlack:
		return escapes.TextColorBlack
	case colorRed:
		return escapes.TextColorRed
	case colorGreen:
		return escapes.TextColorGreen
	case colorYellow:
		return escapes.TextColorYellow
	case colorBlue:
		return escapes.TextColorBlue
	case colorMagenta:
		return escapes.TextColorMagenta
	case colorCyan:
		return escapes.TextColorCyan
	case colorWhite:
		return escapes.TextColorWhite
	}

	panic("invalid color")
}

func (c color) toBackgroundEscapes() string {
	switch c {
	case colorBlack:
		return escapes.BackgroundColorBlack
	case colorRed:
		return escapes.BackgroundColorRed
	case colorGreen:
		return escapes.BackgroundColorGreen
	case colorYellow:
		return escapes.BackgroundColorYellow
	case colorBlue:
		return escapes.BackgroundColorBlue
	case colorMagenta:
		return escapes.BackgroundColorMagenta
	case colorCyan:
		return escapes.BackgroundColorCyan
	case colorWhite:
		return escapes.BackgroundColorWhite
	}

	panic("invalid color")
}
