package prompt

import (
	editor "github.com/JosephNaberhaus/go-text-editor"
	"github.com/eiannone/keyboard"
)

type key interface {
	isText() bool
	rune() rune
}

type runeKey rune

func (r runeKey) isText() bool {
	return true
}

func (r runeKey) rune() rune {
	return rune(r)
}

type controlKey uint8

func (c controlKey) isText() bool {
	if c == controlSpace {
		return true
	}

	return false
}

func (c controlKey) rune() rune {
	if c == controlSpace {
		return ' '
	}

	return 0
}

const (
	noop controlKey = iota
	controlLeft
	controlRight
	controlUp
	controlDown
	controlEnter
	controlBackspace
	controlSpace
	controlHome
	controlEnd
)

func ToKey(rune rune, key keyboard.Key) key {
	if rune != 0 {
		return runeKey(rune)
	}

	switch key {
	case keyboard.KeyArrowLeft:
		return controlLeft
	case keyboard.KeyArrowRight:
		return controlRight
	case keyboard.KeyArrowUp:
		return controlUp
	case keyboard.KeyArrowDown:
		return controlDown
	case keyboard.KeyEnter:
		return controlEnter
	case keyboard.KeyBackspace:
		fallthrough
	case keyboard.KeyBackspace2:
		return controlBackspace
	case keyboard.KeySpace:
		return controlSpace
	case keyboard.KeyHome:
		return controlHome
	case keyboard.KeyEnd:
		return controlEnd
	default:
		return noop
	}
}

func applyKeyToEditor(k key, editor *editor.TextEditor) {
	if k.isText() {
		editor.Write(string(k.rune()))
	}

	switch k {
	case controlLeft:
		editor.Left()
	case controlRight:
		editor.Right()
	case controlUp:
		editor.Up()
	case controlDown:
		editor.Down()
	case controlEnter:
		editor.Newline()
	case controlBackspace:
		editor.Backspace()
	case controlHome:
		editor.Home()
	case controlEnd:
		editor.End()
	}
}
