package prompt

import (
	"errors"
	"fmt"
	"github.com/eiannone/keyboard"
)

type State int

const (
	Waiting State = iota
	Showing
	Finished
)

type base struct {
	output *output
	promptState  State
}

func (b *base) show() error {
	if b.promptState == Showing {
		return errors.New("cannot show a prompt multiple times")
	}

	if b.promptState == Finished {
		return errors.New("cannot show a finished prompt")
	}

	err := keyboard.Open()
	if err != nil {
		return fmt.Errorf("can't listen to keyboard: %w", err)
	}

	output, err := newOutput()
	if err != nil {
		return err
	}

	b.output = output
	b.promptState = Showing

	return nil
}

func (b *base) Pause() error {
	if b.promptState == Waiting {
		return errors.New("cannot pause a prompt when it is already waiting")
	}

	if b.promptState == Finished {
		return errors.New("cannot pause a finished prompt")
	}

	b.output.clear()
	b.output.flush()
	b.promptState = Waiting

	err := keyboard.Close()
	if err != nil {
		panic(err)
	}

	return nil
}

func (b *base) ResetToWaiting() error {
	if b.promptState == Waiting {
		return errors.New("cannot reshow a response that is waiting")
	}

	if b.promptState == Showing {
		return errors.New("cannot reshow a response that is showing")
	}

	b.output.uncommit()
	b.promptState = Waiting
	return nil
}

func (b *base) finish() {
	err := keyboard.Close()
	if err != nil {
		panic(err)
	}

	b.output.commit()
	b.promptState = Finished
}

func (b *base) State() State {
	return b.promptState
}

func (b *base) nextKey() (Key, error) {
	r, key, err := keyboard.GetKey()
	if err != nil {
		if err.Error() == "Unrecognized escape sequence" {
			return b.nextKey()
		}
		return nil, fmt.Errorf("error getting key input: %w", err)
	}

	if key == keyboard.KeyCtrlC {
		return nil, errors.New("prompt loop aborted")
	}

	return ToKey(r, key), nil
}
