package prompt

type Prompt interface {
	Show() error
	Pause() error
	ResetToWaiting() error
	State() State
}
