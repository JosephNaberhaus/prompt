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
	state  State
}

func (b *base) Show() error {
	if b.state == Showing {
		return errors.New("cannot show a prompt multiple times")
	}

	if b.state == Finished {
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
	b.state = Showing

	return nil
}

func (b *base) Finish() {
	err := keyboard.Close()
	if err != nil {
		panic(err)
	}

	b.output.commit()
	b.state = Finished
}

func (b *base) State() State {
	return b.state
}

func (b *base) nextKey() (key, error) {
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
