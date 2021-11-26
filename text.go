package prompt

import (
	"fmt"
	editor "github.com/JosephNaberhaus/go-text-editor"
	"unicode"
)

type Text struct {
	base

	// The question to display to the user
	Question string

	// Whether only a single line of input should be accepted
	IsSingleLine bool

	// Validates the current input given as an array of paragraphs. Return an empty string when the input is valid and
	// return a message to display to the user when it is not valid.
	ValidatorFunc func([]string) string

	// Whether to show the character count to the user
	ShouldShowCharacterCount bool

	// Whether all input should be converted to lowercase
	ShouldForceLowercase bool

	// The line length to wrap text to after the user submits
	OnSubmitMaxLineLength int

	didAttemptSubmit bool

	editor *editor.TextEditor
}

func (t *Text) Show() error {
	err := t.base.Show()
	if err != nil {
		return err
	}

	t.editor = editor.NewEditor()
	t.editor.SetWidth(t.output.outputWidth)
	t.render(false)

	for t.state == Showing {
		nextKey, err := t.nextKey()
		if err != nil {
			t.Finish()
			return err
		}

		t.handleInput(nextKey)
	}

	return nil
}

func (t *Text) handleInput(input key) {
	if t.state != Showing {
		return
	}

	t.didAttemptSubmit = false
	isFinished := false

	if input == controlEnter {
		if t.IsSingleLine || t.editor.Empty() {
			t.didAttemptSubmit = true
			isFinished = t.Validate() == ""
		} else {
			paragraphs := t.editor.Paragraphs()

			lastParagraphsAreEmpty := len(paragraphs) > 0 && paragraphs[len(paragraphs)-1] == "" && paragraphs[len(paragraphs)-2] == ""
			if lastParagraphsAreEmpty && t.editor.CursorIsOnLastParagraph() {
				t.didAttemptSubmit = true
				isFinished = t.Validate() == ""

				// The enter is going to happen below, so always remove at least on backspace to keep the cursor in the
				// same row.
				t.editor.Backspace()
				if isFinished {
					t.editor.Backspace()
				}
			}
		}
	}

	if isFinished {
		if t.OnSubmitMaxLineLength > 0 {
			t.editor.SetWidth(t.OnSubmitMaxLineLength)
		}

		t.render(true)
		t.Finish()
		return
	}

	if input != controlEnter || !t.IsSingleLine {
		if t.ShouldForceLowercase && input.isText() {
			input = runeKey(unicode.ToLower(input.rune()))
		}

		applyKeyToEditor(input, t.editor)
	}

	t.render(false)
}

func (t *Text) render(isFinished bool) {
	t.output.clear()

	t.output.writeColor("? ", colorGreen)
	t.output.write(t.Question)

	validatorMessage := t.Validate()
	isValid := validatorMessage == ""

	if !t.IsSingleLine {
		if t.editor.Empty() && isValid {
			t.output.write(": (press enter to skip)")
		} else {
			t.output.write(": (enter two empty lines to submit)")
		}
	}

	t.output.writeLn(":")

	if t.ShouldShowCharacterCount {
		prefix := fmt.Sprintf("(%d) ", t.editor.NumGraphemes())
		t.editor.SetFirstLineIndent(len(prefix))

		if isValid {
			t.output.writeColor(prefix, colorCyan)
		} else {
			t.output.writeColor(prefix, colorRed)
		}
	}

	if isFinished {
		t.output.writeColor(t.editor.String(), colorCyan)
	} else {

		if isValid {
			t.output.writeColor(t.editor.String(), colorWhite)
		} else {
			t.output.writeColor(t.editor.String(), colorRed)

			if t.didAttemptSubmit {
				t.output.nextLine()
				t.output.writeColor(">> ", colorRed)
				t.output.write(validatorMessage)
			}
		}

	}

	t.output.setCursor(t.editor.CursorRow()+1, t.editor.CursorColumn())
	t.output.flush()
}

func (t *Text) Validate() string {
	if t.ValidatorFunc != nil {
		return t.ValidatorFunc(t.editor.Paragraphs())
	}

	return ""
}

func (t *Text) Response() string {
	return t.editor.String()
}