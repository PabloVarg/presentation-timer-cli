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

type SectionItem struct {
	ID       int
	Name     string
	Duration string
	Position int
}

func SectionItemizer(items []api.Section) []list.Item {
	result := make([]list.Item, 0, len(items))

	for _, item := range items {
		result = append(result, SectionItem{
			ID:       item.ID,
			Name:     item.Name,
			Duration: item.Duration.String(),
			Position: item.Position,
		})
	}

	return result
}

func (i SectionItem) Title() string {
	return i.Name
}

func (i SectionItem) Description() string {
	return i.Duration
}

func (i SectionItem) FilterValue() string {
	return i.Name
}
