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

	// Called when a key is pressed but before it is processed. Return `false` to cancel the event.
	OnKeyFunc func(Prompt, Key) bool

	editor *editor.TextEditor
}

// Show displays the prompt to the user and blocks the current Go routine until the user submits
func (b *Boolean) Show() error {
	err := b.show()
	if err != nil {
		return err
	}

	b.editor = editor.NewEditor()
	b.editor.SetWidth(b.output.outputWidth)
	b.render(false)

	for b.promptState == Showing {
		nextKey, err := b.nextKey()
		if err != nil {
			b.finish()
			return err
		}

		b.handleInput(nextKey)
	}

	return nil
}

func (b *Boolean) handleInput(input Key) {
	if b.State() != Showing {
		return
	}

	if b.OnKeyFunc != nil && !b.OnKeyFunc(b, input) {
		return
	}

	if input == ControlEnter {
		b.render(true)
		b.finish()
		return
	}

	applyKeyToEditor(input, b.editor)

	if b.State() != Waiting {
		b.render(false)
	}
}

func (b *Boolean) render(isFinished bool) {
	b.output.clear()

	b.output.writeColor("? ", colorGreen)
	b.output.write(b.Question)

	if b.defaultResponse() {
		b.output.writeColor("(Y/n) ", colorGreen)
	} else {
		b.output.writeColor(" (y/N) ", colorGreen)
	}

	b.editor.SetFirstLineIndent(b.output.cursorColumn)

	if isFinished {
		if b.Response() {
			b.output.writeColor("Yes", colorCyan)
		} else {
			b.output.writeColor("No", colorCyan)
		}
	} else {
		b.output.write(b.editor.String())
	}

	b.output.setCursor(b.editor.CursorRow(), b.editor.CursorColumn())
	b.output.flush()
}

func (b *Boolean) defaultResponse() bool {
	if b.IsTrueFunc != nil {
		return b.IsTrueFunc("")
	}

	return defaultIsYes("")
}

// Response returns the input from the user.
func (b *Boolean) Response() bool {
	input := ""
	if b.editor != nil {
		input = b.editor.String()
	}

	if b.IsTrueFunc != nil {
		return b.IsTrueFunc(input)
	}

	return defaultIsYes(input)
}

func defaultIsYes(input string) bool {
	lowercase := strings.ToLower(input)
	return lowercase == "y" || lowercase == "Y"
}
