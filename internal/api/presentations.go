package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/helpers"
)

type Presentation struct {
	Name     string
	Duration time.Duration
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

type CreatePresentationMsg struct {
	Name string `json:"name"`
}

func CreatePresentation(client APIClient, get KeyValueRetriever, msg CreatePresentationMsg) error {
	path, err := url.JoinPath(client.Url(get), "/presentations")
	if err != nil {
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res, err := client.HTTPClient.Post(path, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating resource")
	}

	return nil
}
