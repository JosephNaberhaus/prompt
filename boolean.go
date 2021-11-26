package prompt

import (
	editor "github.com/JosephNaberhaus/go-text-editor"
	"strings"
)

type Boolean struct {
	base

	// The question to display to the user
	Question string

	// A function that will be called to determine if the user's input is equivalent to true.
	// Defaults to defaultIsYes which returns true for "y" or "yes" (ignore capitalization).
	IsTrueFunc func(string) bool

	editor *editor.TextEditor
}

func (b *Boolean) Show() error {
	err := b.base.Show()
	if err != nil {
		return err
	}

	b.editor = editor.NewEditor()
	b.editor.SetWidth(b.output.outputWidth)
	b.render(false)

	for b.state == Showing {
		nextKey, err := b.nextKey()
		if err != nil {
			b.Finish()
			return err
		}

		b.handleInput(nextKey)
	}

	return nil
}

func (b *Boolean) handleInput(input key) {
	if b.state != Showing {
		return
	}

	if input == controlEnter {
		b.render(true)
		b.Finish()
		return
	}

	applyKeyToEditor(input, b.editor)
	b.render(false)
}

func (b *Boolean) render(isFinished bool) {
	b.output.clear()

	b.output.writeColor("? ", colorGreen)
	b.output.write(b.Question)
	b.output.writeColor(" (y/N) ", colorGreen)

	b.editor.SetFirstLineIndent(b.output.cursorColumn)

	if isFinished {
		if b.Response() {
			b.output.writeColorLn("Yes", colorCyan)
		} else {
			b.output.writeColorLn("No", colorCyan)
		}
	} else {
		b.output.write(b.editor.String())
	}

	b.output.setCursor(b.editor.CursorRow(), b.editor.CursorColumn())
	b.output.flush()
}

func (b *Boolean) Response() bool {
	if b.IsTrueFunc != nil {
		return b.IsTrueFunc(b.editor.String())
	}

	return defaultIsYes(b.editor.String())
}

func defaultIsYes(input string) bool {
	lowercase := strings.ToLower(input)
	return lowercase == "y" || lowercase == "Y"
}
