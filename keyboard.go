package prompt

import (
	editor "github.com/JosephNaberhaus/go-text-editor"
	"github.com/eiannone/keyboard"
)

type Key interface {
	IsText() bool
	Rune() rune
}

type RuneKey rune

func (r RuneKey) IsText() bool {
	return true
}

func (r RuneKey) Rune() rune {
	return rune(r)
}

type ControlKey uint8

func (c ControlKey) IsText() bool {
	if c == ControlSpace {
		return true
	}

	return false
}

func (c ControlKey) Rune() rune {
	if c == ControlSpace {
		return ' '
	}

	return 0
}

const (
	Noop ControlKey = iota
	ControlLeft
	ControlRight
	ControlUp
	ControlDown
	ControlEnter
	ControlBackspace
	ControlSpace
	ControlHome
	ControlEnd
	ControlCtrlA
	ControlCtrlB
	ControlCtrlC
	ControlCtrlD
	ControlCtrlE
	ControlCtrlF
	ControlCtrlG
	ControlCtrlI
	ControlCtrlJ
	ControlCtrlK
	ControlCtrlL
	ControlCtrlN
	ControlCtrlO
	ControlCtrlP
	ControlCtrlQ
	ControlCtrlR
	ControlCtrlS
	ControlCtrlT
	ControlCtrlU
	ControlCtrlV
	ControlCtrlW
	ControlCtrlX
	ControlCtrlY
	ControlCtrlZ
)

func ToKey(rune rune, key keyboard.Key) Key {
	if rune != 0 {
		return RuneKey(rune)
	}

	switch key {
	case keyboard.KeyArrowLeft:
		return ControlLeft
	case keyboard.KeyArrowRight:
		return ControlRight
	case keyboard.KeyArrowUp:
		return ControlUp
	case keyboard.KeyArrowDown:
		return ControlDown
	case keyboard.KeyEnter:
		return ControlEnter
	case keyboard.KeyBackspace:
		fallthrough
	case keyboard.KeyBackspace2:
		return ControlBackspace
	case keyboard.KeySpace:
		return ControlSpace
	case keyboard.KeyHome:
		return ControlHome
	case keyboard.KeyEnd:
		return ControlEnd
	case keyboard.KeyCtrlA:
		return ControlCtrlA
	case keyboard.KeyCtrlB:
		return ControlCtrlB
	case keyboard.KeyCtrlC:
		return ControlCtrlC
	case keyboard.KeyCtrlD:
		return ControlCtrlD
	case keyboard.KeyCtrlE:
		return ControlCtrlE
	case keyboard.KeyCtrlF:
		return ControlCtrlF
	case keyboard.KeyCtrlG:
		return ControlCtrlG
	case keyboard.KeyCtrlI:
		return ControlCtrlI
	case keyboard.KeyCtrlJ:
		return ControlCtrlJ
	case keyboard.KeyCtrlK:
		return ControlCtrlK
	case keyboard.KeyCtrlL:
		return ControlCtrlL
	case keyboard.KeyCtrlN:
		return ControlCtrlN
	case keyboard.KeyCtrlO:
		return ControlCtrlO
	case keyboard.KeyCtrlP:
		return ControlCtrlP
	case keyboard.KeyCtrlQ:
		return ControlCtrlQ
	case keyboard.KeyCtrlR:
		return ControlCtrlR
	case keyboard.KeyCtrlS:
		return ControlCtrlS
	case keyboard.KeyCtrlT:
		return ControlCtrlT
	case keyboard.KeyCtrlU:
		return ControlCtrlU
	case keyboard.KeyCtrlV:
		return ControlCtrlV
	case keyboard.KeyCtrlW:
		return ControlCtrlW
	case keyboard.KeyCtrlX:
		return ControlCtrlX
	case keyboard.KeyCtrlY:
		return ControlCtrlY
	case keyboard.KeyCtrlZ:
		return ControlCtrlZ
	default:
		return Noop
	}
}

func applyKeyToEditor(k Key, editor *editor.TextEditor) {
	if k.IsText() {
		editor.Write(string(k.Rune()))
	}

	switch k {
	case ControlLeft:
		editor.Left()
	case ControlRight:
		editor.Right()
	case ControlUp:
		editor.Up()
	case ControlDown:
		editor.Down()
	case ControlEnter:
		editor.Newline()
	case ControlBackspace:
		editor.Backspace()
	case ControlHome:
		editor.Home()
	case ControlEnd:
		editor.End()
	}
}
