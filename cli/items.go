package cli

import (
	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
)

type PresentationItem struct {
	id       int
	name     string
	duration string
}

func PresentationItemizer(items []api.Presentation) []list.Item {
	result := make([]list.Item, 0, len(items))

	for _, item := range items {
		result = append(result, PresentationItem{
			id:       item.ID,
			name:     item.Name,
			duration: item.Duration.String(),
		})
	}

	return result
}

func (i PresentationItem) Title() string {
	return i.name
}

func (i PresentationItem) Description() string {
	return i.duration
}

func (i PresentationItem) FilterValue() string {
	return i.name
}
