package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/helpers"
)

type Presentation struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

func GetPresentations(client APIClient, get KeyValueRetriever) ([]Presentation, error) {
	path, err := url.JoinPath(client.Url(get), "/presentations")
	if err != nil {
		return nil, err
	}
	path += "?page_size=100"

	res, err := client.HTTPClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		response, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", string(response))
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
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return ExtractErrorMsg(res.Body)
	}

	return nil
}

func DeletePresentation(client APIClient, get KeyValueRetriever, ID int) error {
	path, err := url.JoinPath(client.Url(get), fmt.Sprintf("/presentations/%d", ID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ExtractErrorMsg(res.Body)
	}

	return nil
}
