package cli

import (
	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
)

type PresentationItem struct {
	ID       int
	Name     string
	Duration string
}

func PresentationItemizer(items []api.Presentation) []list.Item {
	result := make([]list.Item, 0, len(items))

	for _, item := range items {
		result = append(result, PresentationItem{
			ID:       item.ID,
			Name:     item.Name,
			Duration: item.Duration.String(),
		})
	}

	return result
}

func (i PresentationItem) Title() string {
	return i.Name
}

func (i PresentationItem) Description() string {
	return i.Duration
}

func (i PresentationItem) FilterValue() string {
	return i.Name
}
