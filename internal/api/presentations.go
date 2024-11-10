package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/helpers"
)

type Presentation struct {
	Name string
}

func GetPresentations(client APIClient, get KeyValueRetriever) ([]Presentation, error) {
	path, err := url.JoinPath(client.Url(get), "/presentations")
	if err != nil {
		return nil, err
	}

	res, err := client.HTTPClient.Get(path)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request not successful")
	}

	var result PaginatedResponse[[]Presentation]
	if err := helpers.ReadJSON(res.Body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
