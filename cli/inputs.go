package cli

import "github.com/charmbracelet/bubbles/textinput"

func NewDefaultTextInput() textinput.Model {
	input := textinput.New()
	input.PromptStyle = promptStyle
	input.TextStyle = textInputStyle
	input.Cursor.Style = cursorStyle

	return input
}
