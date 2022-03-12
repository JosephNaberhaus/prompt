package main

import "github.com/JosephNaberhaus/prompt"

func main() {
	input1 := prompt.Text{
		Question:                 "Input 1",
		IsSingleLine:             true,
	}

	input2 := prompt.Text{
		Question:                 "Input 2",
		IsSingleLine:             true,
		OnKeyFunc:                goBack,
	}

	for input2.State() != prompt.Finished {
		if input1.State() == prompt.Finished {
			input1.ResetToWaiting()
		}

		input1.Show()

		input2.Show()
	}
}

func goBack(p prompt.Prompt, key prompt.Key) bool {
	if key == prompt.ControlCtrlB {
		_ = p.Pause()
		return false
	}

	return true
}