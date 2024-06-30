package main

import (
	"fmt"
	"github.com/JosephNaberhaus/prompt"
)

func main() {
	var options []prompt.SelectionOption
	for i := 0; i < 30; i++ {
		var desc string
		if i%3 != 0 {
			desc = "a description too"
		}

		options = append(options, prompt.SelectionOption{
			Name:        fmt.Sprintf("Option %d", i),
			Description: desc,
		})
	}

	p := prompt.Select{
		Question:      "Select one",
		Options:       options,
		NumLinesShown: 10,
	}

	err := p.Show()
	if err != nil {
		panic(err)
	}

	println("You selected:", p.Response().Name)
}
