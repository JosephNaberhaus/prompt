# Prompt
Another GoLang prompt library with a couple of uncommon features:

- Supports more advanced text editing by using my [Go Text Editor](https://github.com/JosephNaberhaus/go-text-editor) library
    - Multiline text editing
    - Better cursor movement (`up`, `down`, `home`, and `end` all work)
- Doesn't clear out and of your screen

## Supported Prompts
- Yes/No questions with `Boolean{<options>}`
- Select from list with `Select{<options>}`
- Text (multiline and single line) with `Text{<options>}`

## Usage
All prompts are created by initializing their respective struct.

```go
input := prompt.Text{}
```

The zero value of the prompt is ready to use, but you will probably want to explore the available public members to customize the prompt.

```go
input := prompt.Text{
    Question: "What is your name?",
    IsSingleLine: true,
    ValidatorFunc: func(input []string) string {
        if input[0] == "" {
            return "name is required"
        }   

        return ""
    },
}
```

When the prompt is ready, call `Show` to display the prompt and block the active Go-Routine until a response is submitted. Then call `Response` to get the output of the prompt.

```go
err := input.Show()
if err != nil {
    // Handle error
}

name := input.Response()
fmt.Printf("Hello %s!", name)
```