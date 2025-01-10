package cli

import "github.com/charmbracelet/bubbles/textinput"

func NewDefaultTextInput() textinput.Model {
	input := textinput.New()
	input.PromptStyle = PromptStyle
	input.TextStyle = TextInputStyle
	input.Cursor.Style = CursorStyle

	return input
}
